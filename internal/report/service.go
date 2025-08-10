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
	r.Log.Infof("GetSalesReports: %v", reports)
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

func NewRptService(store repository.Store, activityLogService activitylog.IActivityService, log logger.ILogger) IRptService {
	return &RptService{
		Store:              store,
		ActivityLogService: activityLogService,
		Log:                log,
	}
}
