package report

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v3"
)

type IRptHandler interface {
	GetDashboardSummaryHandler(c fiber.Ctx) error
	GetSalesReportsHandler(c fiber.Ctx) error
	GetProductPerformanceHandler(c fiber.Ctx) error
	GetPaymentMethodPerformanceHandler(c fiber.Ctx) error
	GetCashierPerformanceHandler(c fiber.Ctx) error
	GetCancellationReportsHandler(c fiber.Ctx) error
	GetProfitSummaryHandler(c fiber.Ctx) error
	GetProductProfitReportsHandler(c fiber.Ctx) error
}

type RptHandler struct {
	Service IRptService
	log     logger.ILogger
}

// GetSalesReportsHandler retrieves sales reports within a date range
// @Summary Get sales reports
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} common.SuccessResponse{data=[]dto.SalesReport}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /reports/sales [get]
func (r *RptHandler) GetSalesReportsHandler(c fiber.Ctx) error {
	startDate, endDate, err := r.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	serviceReq := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	salesReports, err := r.Service.GetSalesReports(c.RequestCtx(), serviceReq)
	if err != nil {
		r.log.Error("Failed to get sales reports", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get sales reports",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Sales reports retrieved successfully",
		Data:    salesReports,
	})
}

// GetProductPerformanceHandler retrieves product performance analytics
// @Summary Get product performance
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} common.SuccessResponse{data=[]dto.ProductPerformanceResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /reports/products [get]
func (r *RptHandler) GetProductPerformanceHandler(c fiber.Ctx) error {
	startDate, endDate, err := r.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	serviceReq := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	results, err := r.Service.GetProductPerformance(c.RequestCtx(), serviceReq)
	if err != nil {
		r.log.Error("Failed to get product performance", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get product performance",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product performance retrieved successfully",
		Data:    results,
	})
}

// GetPaymentMethodPerformanceHandler retrieves payment method performance analytics
// @Summary Get payment method performance
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} common.SuccessResponse{data=[]dto.PaymentMethodPerformanceResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /reports/payment-methods [get]
func (r *RptHandler) GetPaymentMethodPerformanceHandler(c fiber.Ctx) error {
	startDate, endDate, err := r.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	serviceReq := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	results, err := r.Service.GetPaymentMethodPerformance(c.RequestCtx(), serviceReq)
	if err != nil {
		r.log.Error("Failed to get payment method performance", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get payment method performance",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Payment method performance retrieved successfully",
		Data:    results,
	})
}

// GetCashierPerformanceHandler retrieves cashier performance analytics
// @Summary Get cashier performance
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} common.SuccessResponse{data=[]dto.CashierPerformanceResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /reports/cashier-performance [get]
func (r *RptHandler) GetCashierPerformanceHandler(c fiber.Ctx) error {
	startDate, endDate, err := r.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	serviceReq := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	results, err := r.Service.GetCashierPerformance(c.RequestCtx(), serviceReq)
	if err != nil {
		r.log.Error("Failed to get cashier performance", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get cashier performance",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Cashier performance retrieved successfully",
		Data:    results,
	})
}

// GetCancellationReportsHandler retrieves cancellation reports
// @Summary Get cancellation reports
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} common.SuccessResponse{data=[]dto.CancellationReportResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /reports/cancellations [get]
func (r *RptHandler) GetCancellationReportsHandler(c fiber.Ctx) error {
	startDate, endDate, err := r.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	serviceReq := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	results, err := r.Service.GetCancellationReports(c.RequestCtx(), serviceReq)
	if err != nil {
		r.log.Error("Failed to get cancellation reports", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get cancellation reports",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Cancellation reports retrieved successfully",
		Data:    results,
	})
}

// GetDashboardSummaryHandler retrieves dashboard summary
// @Summary Get dashboard summary
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {object} common.SuccessResponse{data=dto.DashboardSummaryResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /reports/dashboard-summary [get]
func (r *RptHandler) GetDashboardSummaryHandler(c fiber.Ctx) error {

	summary, err := r.Service.GetDashboardSummary(c.RequestCtx())
	if err != nil {
		r.log.Error("Failed to get dashboard summary", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get dashboard summary",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Dashboard summary retrieved successfully",
		Data:    summary,
	})
}

// GetProfitSummaryHandler retrieves profit summary
// @Summary Get profit summary
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} common.SuccessResponse{data=[]dto.ProfitSummaryResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /reports/profit-summary [get]
func (r *RptHandler) GetProfitSummaryHandler(c fiber.Ctx) error {
	startDate, endDate, err := r.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	serviceReq := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	results, err := r.Service.GetProfitSummary(c.RequestCtx(), serviceReq)
	if err != nil {
		r.log.Error("Failed to get profit summary", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get profit summary",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Profit summary retrieved successfully",
		Data:    results,
	})
}

// GetProductProfitReportsHandler retrieves product profit reports
// @Summary Get product profit reports
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} common.SuccessResponse{data=[]dto.ProductProfitResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /reports/profit-products [get]
func (r *RptHandler) GetProductProfitReportsHandler(c fiber.Ctx) error {
	startDate, endDate, err := r.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	serviceReq := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	results, err := r.Service.GetProductProfitReports(c.RequestCtx(), serviceReq)
	if err != nil {
		r.log.Error("Failed to get product profit reports", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get product profit reports",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product profit reports retrieved successfully",
		Data:    results,
	})
}

func (r *RptHandler) parseDateRange(c fiber.Ctx) (time.Time, time.Time, error) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	const layout = "2006-01-02"

	if startDateStr == "" {
		// Default to today if not provided (optional behavior, but strict validation requested)
		// Or return error. Let's return error as per doc.
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "start_date is required")
	}

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date format, use YYYY-MM-DD")
	}

	if endDateStr == "" {
		// Default to today
		endDateStr = time.Now().Format(layout)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "Invalid end_date format, use YYYY-MM-DD")
	}

	// Adjust endDate to end of day? Or just date.
	// Postgres date type comparison is inclusive if using BETWEEN 'YYYY-MM-DD' AND 'YYYY-MM-DD' (casted to 00:00:00).
	// If the column is timestamptz, we might need to adjust.
	// The query uses `created_at::date BETWEEN $1 AND $2`, so pure date is fine.

	return startDate, endDate, nil
}

func NewRptHandler(service IRptService, log logger.ILogger) IRptHandler {
	return &RptHandler{
		Service: service,
		log:     log,
	}
}

// fiber:context-methods migrated
