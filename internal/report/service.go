package report

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type IRptService interface {
	GetDashboardSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error)
	GetSalesReports(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.SalesReport, error)
	GetProductPerformance(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.ProductPerformanceResponse, error)
	GetPaymentMethodPerformance(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.PaymentMethodPerformanceResponse, error)
	GetCashierPerformance(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.CashierPerformanceResponse, error)
	GetCancellationReports(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.CancellationReportResponse, error)
}

type RptService struct {
	Store              repository.Store
	ActivityLogService activitylog.IActivityService
	Log                logger.ILogger
}

func (r *RptService) GetSalesReports(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.SalesReport, error) {

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

	reports, err := r.Store.GetSalesSummary(ctx, params)
	// r.Log.Infof("GetSalesReports: %v", reports)
	if err != nil {
		r.Log.Error("Failed to get sales reports", "error", err)
		return nil, err
	}

	salesReports := make([]dto.SalesReport, len(reports))
	for i, report := range reports {
		salesReports[i] = dto.SalesReport{
			Date:       report.Date.Time,
			OrderCount: report.OrderCount,
			TotalSales: 0, // Default value, will be set if valid
		}
		if n, ok := report.TotalSales.(pgtype.Numeric); ok && n.Valid {
			f8, _ := n.Float64Value()
			salesReports[i].TotalSales = f8.Float64
		}
	}
	return &salesReports, nil
}

func (r *RptService) GetDashboardSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error) {
	summary, err := r.Store.GetDashboardSummary(ctx)
	if err != nil {
		r.Log.Error("Failed to get dashboard summary", "error", err)
		return nil, err
	}

	var totalSales float64
	if n, ok := summary.TotalSales.(pgtype.Numeric); ok && n.Valid {
		f8, _ := n.Float64Value()
		totalSales = f8.Float64
	}

	response := &dto.DashboardSummaryResponse{
		TotalSales:    totalSales,
		TotalOrders:   summary.TotalOrders,
		UniqueCashier: summary.UniqueCashiers,
		TotalProducts: summary.TotalProducts,
	}

	return response, nil
}

func (r *RptService) GetProductPerformance(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.ProductPerformanceResponse, error) {
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

	results, err := r.Store.GetProductSalesPerformance(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get product performance", "error", err)
		return nil, err
	}

	responses := make([]dto.ProductPerformanceResponse, len(results))
	for i, row := range results {
		// TotalRevenue is int64 in generated code
		totalRevenue := float64(row.TotalRevenue)

		responses[i] = dto.ProductPerformanceResponse{
			ProductID:     row.ProductID.String(),
			ProductName:   row.ProductName,
			TotalQuantity: row.TotalQuantity,
			TotalRevenue:  totalRevenue,
		}
	}

	return &responses, nil
}

func (r *RptService) GetPaymentMethodPerformance(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.PaymentMethodPerformanceResponse, error) {
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

	results, err := r.Store.GetPaymentMethodSales(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get payment method performance", "error", err)
		return nil, err
	}

	responses := make([]dto.PaymentMethodPerformanceResponse, len(results))
	for i, row := range results {
		var totalSales float64
		// TotalSales is interface{} (numeric or int)
		if n, ok := row.TotalSales.(pgtype.Numeric); ok && n.Valid {
			f8, _ := n.Float64Value()
			totalSales = f8.Float64
		} else if v, ok := row.TotalSales.(int64); ok {
			totalSales = float64(v)
		}

		responses[i] = dto.PaymentMethodPerformanceResponse{
			PaymentMethodID:   row.PaymentMethodID,
			PaymentMethodName: row.PaymentMethodName,
			OrderCount:        row.OrderCount,
			TotalSales:        totalSales,
		}
	}

	return &responses, nil
}

func (r *RptService) GetCashierPerformance(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.CashierPerformanceResponse, error) {
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

	results, err := r.Store.GetCashierPerformance(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get cashier performance", "error", err)
		return nil, err
	}

	responses := make([]dto.CashierPerformanceResponse, len(results))
	for i, row := range results {
		var totalSales float64
		if n, ok := row.TotalSales.(pgtype.Numeric); ok && n.Valid {
			f8, _ := n.Float64Value()
			totalSales = f8.Float64
		} else if v, ok := row.TotalSales.(int64); ok {
			totalSales = float64(v)
		}

		responses[i] = dto.CashierPerformanceResponse{
			UserID:     row.UserID.String(),
			Username:   row.Username,
			OrderCount: row.OrderCount,
			TotalSales: totalSales,
		}
	}

	return &responses, nil
}

func (r *RptService) GetCancellationReports(ctx context.Context, req *dto.SalesReportRequest) (*[]dto.CancellationReportResponse, error) {
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

	results, err := r.Store.GetCancellationReasons(ctx, params)
	if err != nil {
		r.Log.Error("Failed to get cancellation reports", "error", err)
		return nil, err
	}

	responses := make([]dto.CancellationReportResponse, len(results))
	for i, row := range results {
		responses[i] = dto.CancellationReportResponse{
			ReasonID:        row.ReasonID,
			Reason:          row.Reason,
			CancelledOrders: row.CancelledOrders,
		}
	}

	return &responses, nil
}

func NewRptService(store repository.Store, activityLogService activitylog.IActivityService, log logger.ILogger) IRptService {
	return &RptService{
		Store:              store,
		ActivityLogService: activityLogService,
		Log:                log,
	}
}
