package report

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common/store"
	"POS-kasir/internal/report/repository"
	"POS-kasir/pkg/logger"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type IRptService interface {
	GetDashboardSummary(ctx context.Context) (*DashboardSummaryResponse, error)
	GetSalesReports(ctx context.Context, req *SalesReportServiceRequest) (*[]SalesReport, error)
	GetProductPerformance(ctx context.Context, req *SalesReportServiceRequest) (*[]ProductPerformanceResponse, error)
	GetPaymentMethodPerformance(ctx context.Context, req *SalesReportServiceRequest) (*[]PaymentMethodPerformanceResponse, error)
	GetCashierPerformance(ctx context.Context, req *SalesReportServiceRequest) (*[]CashierPerformanceResponse, error)
	GetCancellationReports(ctx context.Context, req *SalesReportServiceRequest) (*[]CancellationReportResponse, error)
	GetProfitSummary(ctx context.Context, req *SalesReportServiceRequest) (*[]ProfitSummaryResponse, error)
	GetProductProfitReports(ctx context.Context, req *SalesReportServiceRequest) (*[]ProductProfitResponse, error)
}

func NewRptService(store store.Store, repo repository.Querier, activityLogService activitylog.IActivityService, log logger.ILogger) IRptService {
	return &RptService{
		repo:               repo,
		Store:              store,
		ActivityLogService: activityLogService,
		Log:                log,
	}
}

type RptService struct {
	repo               repository.Querier
	Store              store.Store
	ActivityLogService activitylog.IActivityService
	Log                logger.ILogger
}

func (r *RptService) GetSalesReports(ctx context.Context, req *SalesReportServiceRequest) (*[]SalesReport, error) {

	params := repository.GetSalesSummaryParams{
		CreatedAt: pgtype.Timestamptz{
			Time:  req.StartDate,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  req.EndDate,
			Valid: true,
		},
	}

	reports, err := r.repo.GetSalesSummary(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get sales reports", "error", err)
		return nil, err
	}

	salesReports := make([]SalesReport, len(reports))
	for i, report := range reports {
		salesReports[i] = SalesReport{
			Date:       report.Date.Time,
			OrderCount: report.OrderCount,
			TotalSales: 0, 
		}
		if n, ok := report.TotalSales.(pgtype.Numeric); ok && n.Valid {
			f8, _ := n.Float64Value()
			salesReports[i].TotalSales = f8.Float64
		}
	}
	return &salesReports, nil
}

func (r *RptService) GetDashboardSummary(ctx context.Context) (*DashboardSummaryResponse, error) {
	summary, err := r.repo.GetDashboardSummary(ctx)
	if err != nil {
		r.Log.Error("Failed to get dashboard summary", "error", err)
		return nil, err
	}

	var totalSales float64
	if n, ok := summary.TotalSales.(pgtype.Numeric); ok && n.Valid {
		f8, _ := n.Float64Value()
		totalSales = f8.Float64
	}

	response := &DashboardSummaryResponse{
		TotalSales:    totalSales,
		TotalOrders:   summary.TotalOrders,
		UniqueCashier: summary.UniqueCashiers,
		TotalProducts: summary.TotalProducts,
	}

	return response, nil
}

func (r *RptService) GetProductPerformance(ctx context.Context, req *SalesReportServiceRequest) (*[]ProductPerformanceResponse, error) {
	params := repository.GetProductSalesPerformanceParams{
		CreatedAt: pgtype.Timestamptz{
			Time:  req.StartDate,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  req.EndDate,
			Valid: true,
		},
	}

	results, err := r.repo.GetProductSalesPerformance(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get product performance", "error", err)
		return nil, err
	}

	responses := make([]ProductPerformanceResponse, len(results))
	for i, row := range results {
		// TotalRevenue is int64 in generated code
		totalRevenue := float64(row.TotalRevenue)

		responses[i] = ProductPerformanceResponse{
			ProductID:     row.ProductID.String(),
			ProductName:   row.ProductName,
			TotalQuantity: row.TotalQuantity,
			TotalRevenue:  totalRevenue,
		}
	}

	return &responses, nil
}

func (r *RptService) GetPaymentMethodPerformance(ctx context.Context, req *SalesReportServiceRequest) (*[]PaymentMethodPerformanceResponse, error) {
	params := repository.GetPaymentMethodSalesParams{
		CreatedAt: pgtype.Timestamptz{
			Time:  req.StartDate,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  req.EndDate,
			Valid: true,
		},
	}

	results, err := r.repo.GetPaymentMethodSales(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get payment method performance", "error", err)
		return nil, err
	}

	responses := make([]PaymentMethodPerformanceResponse, len(results))
	for i, row := range results {
		var totalSales float64
		if n, ok := row.TotalSales.(pgtype.Numeric); ok && n.Valid {
			f8, _ := n.Float64Value()
			totalSales = f8.Float64
		} else if v, ok := row.TotalSales.(int64); ok {
			totalSales = float64(v)
		}

		responses[i] = PaymentMethodPerformanceResponse{
			PaymentMethodID:   row.PaymentMethodID,
			PaymentMethodName: row.PaymentMethodName,
			OrderCount:        row.OrderCount,
			TotalSales:        totalSales,
		}
	}

	return &responses, nil
}

func (r *RptService) GetCashierPerformance(ctx context.Context, req *SalesReportServiceRequest) (*[]CashierPerformanceResponse, error) {
	params := repository.GetCashierPerformanceParams{
		CreatedAt: pgtype.Timestamptz{
			Time:  req.StartDate,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  req.EndDate,
			Valid: true,
		},
	}

	results, err := r.repo.GetCashierPerformance(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get cashier performance", "error", err)
		return nil, err
	}

	responses := make([]CashierPerformanceResponse, len(results))
	for i, row := range results {
		var totalSales float64
		if n, ok := row.TotalSales.(pgtype.Numeric); ok && n.Valid {
			f8, _ := n.Float64Value()
			totalSales = f8.Float64
		} else if v, ok := row.TotalSales.(int64); ok {
			totalSales = float64(v)
		}

		responses[i] = CashierPerformanceResponse{
			UserID:     row.UserID.String(),
			Username:   row.Username,
			OrderCount: row.OrderCount,
			TotalSales: totalSales,
		}
	}

	return &responses, nil
}

func (r *RptService) GetCancellationReports(ctx context.Context, req *SalesReportServiceRequest) (*[]CancellationReportResponse, error) {
	params := repository.GetCancellationReasonsParams{
		CreatedAt: pgtype.Timestamptz{
			Time:  req.StartDate,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  req.EndDate,
			Valid: true,
		},
	}

	results, err := r.repo.GetCancellationReasons(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get cancellation reports", "error", err)
		return nil, err
	}

	responses := make([]CancellationReportResponse, len(results))
	for i, row := range results {
		responses[i] = CancellationReportResponse{
			ReasonID:        row.ReasonID,
			Reason:          row.Reason,
			CancelledOrders: row.CancelledOrders,
		}
	}

	return &responses, nil
}

func (r *RptService) GetProfitSummary(ctx context.Context, req *SalesReportServiceRequest) (*[]ProfitSummaryResponse, error) {
	params := repository.GetProfitSummaryParams{
		CreatedAt: pgtype.Timestamptz{
			Time:  req.StartDate,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  req.EndDate,
			Valid: true,
		},
	}

	results, err := r.repo.GetProfitSummary(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get profit summary", "error", err)
		return nil, err
	}

	responses := make([]ProfitSummaryResponse, len(results))
	for i, row := range results {
		var totalRevenue, totalCOGS, grossProfit float64

		// TotalRevenue (interface{})
		if n, ok := row.TotalRevenue.(pgtype.Numeric); ok && n.Valid {
			f, _ := n.Float64Value()
			totalRevenue = f.Float64
		} else if v, ok := row.TotalRevenue.(int64); ok {
			totalRevenue = float64(v)
		} else if v, ok := row.TotalRevenue.(float64); ok {
			totalRevenue = v
		}

		// TotalCogs (interface{})
		if n, ok := row.TotalCogs.(pgtype.Numeric); ok && n.Valid {
			f, _ := n.Float64Value()
			totalCOGS = f.Float64
		} else if v, ok := row.TotalCogs.(int64); ok {
			totalCOGS = float64(v)
		} else if v, ok := row.TotalCogs.(float64); ok {
			totalCOGS = v
		}

		// GrossProfit (int32)
		grossProfit = float64(row.GrossProfit)

		var date time.Time
		if row.Date.Valid {
			date = row.Date.Time
		}

		responses[i] = ProfitSummaryResponse{
			Date:         date,
			TotalRevenue: totalRevenue,
			TotalCOGS:    totalCOGS,
			GrossProfit:  grossProfit,
		}
	}

	return &responses, nil
}

func (r *RptService) GetProductProfitReports(ctx context.Context, req *SalesReportServiceRequest) (*[]ProductProfitResponse, error) {
	params := repository.GetProductProfitReportsParams{
		CreatedAt: pgtype.Timestamptz{
			Time:  req.StartDate,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamptz{
			Time:  req.EndDate,
			Valid: true,
		},
	}

	results, err := r.repo.GetProductProfitReports(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get product profit reports", "error", err)
		return nil, err
	}

	responses := make([]ProductProfitResponse, len(results))
	for i, row := range results {
		var totalRevenue, totalCOGS, grossProfit float64

		// Direct fields are int64/int32 in generated struct
		totalRevenue = float64(row.TotalRevenue)
		totalCOGS = float64(row.TotalCogs)
		grossProfit = float64(row.GrossProfit)

		responses[i] = ProductProfitResponse{
			ProductID:    row.ProductID.String(),
			ProductName:  row.ProductName,
			TotalSold:    row.TotalSold,
			TotalRevenue: totalRevenue,
			TotalCOGS:    totalCOGS,
			GrossProfit:  grossProfit,
		}
	}

	return &responses, nil
}
