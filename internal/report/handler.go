package report

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"errors"
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
// @Summary      Get sales reports
// @Description  Get aggregated sales data grouped by date within a specified range (Roles: admin, manager, cashier)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        start_date query string true "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string true "End Date (YYYY-MM-DD)"
// @Success      200 {object} common.SuccessResponse{data=[]SalesReport} "Sales reports retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /reports/sales [get]
func (r *RptHandler) GetSalesReportsHandler(c fiber.Ctx) error {
	var req SalesReportRequest
	if err := c.Bind().Query(&req); err != nil {
		r.log.Warnf("Get sales reports validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	serviceReq := &SalesReportServiceRequest{
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
// @Summary      Get product performance
// @Description  Get sales performance metrics for each product (Roles: admin, manager, cashier)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        start_date query string true "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string true "End Date (YYYY-MM-DD)"
// @Success      200 {object} common.SuccessResponse{data=[]ProductPerformanceResponse} "Product performance data retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /reports/products [get]
func (r *RptHandler) GetProductPerformanceHandler(c fiber.Ctx) error {
	var req SalesReportRequest
	if err := c.Bind().Query(&req); err != nil {
		r.log.Warnf("Get product performance validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	serviceReq := &SalesReportServiceRequest{
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
// @Summary      Get payment method performance
// @Description  Get usage counts and totals for each payment method (Roles: admin, manager, cashier)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        start_date query string true "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string true "End Date (YYYY-MM-DD)"
// @Success      200 {object} common.SuccessResponse{data=[]PaymentMethodPerformanceResponse} "Payment method performance data retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /reports/payment-methods [get]
func (r *RptHandler) GetPaymentMethodPerformanceHandler(c fiber.Ctx) error {
	var req SalesReportRequest
	if err := c.Bind().Query(&req); err != nil {
		r.log.Warnf("Get payment method performance validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	serviceReq := &SalesReportServiceRequest{
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
// @Summary      Get cashier performance
// @Description  Get order counts and sales totals handled by each cashier (Roles: admin, manager, cashier)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        start_date query string true "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string true "End Date (YYYY-MM-DD)"
// @Success      200 {object} common.SuccessResponse{data=[]CashierPerformanceResponse} "Cashier performance data retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /reports/cashier-performance [get]
func (r *RptHandler) GetCashierPerformanceHandler(c fiber.Ctx) error {
	var req SalesReportRequest
	if err := c.Bind().Query(&req); err != nil {
		r.log.Warnf("Get cashier performance validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	serviceReq := &SalesReportServiceRequest{
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
// @Summary      Get cancellation reports
// @Description  Get statistics on order cancellations grouped by reason (Roles: admin, manager, cashier)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        start_date query string true "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string true "End Date (YYYY-MM-DD)"
// @Success      200 {object} common.SuccessResponse{data=[]CancellationReportResponse} "Cancellation reports retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /reports/cancellations [get]
func (r *RptHandler) GetCancellationReportsHandler(c fiber.Ctx) error {
	var req SalesReportRequest
	if err := c.Bind().Query(&req); err != nil {
		r.log.Warnf("Get cancellation reports validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	serviceReq := &SalesReportServiceRequest{
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
// @Summary      Get dashboard summary
// @Description  Get high-level summary metrics (totals) for the dashboard (Roles: admin, manager, cashier)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Success      200 {object} common.SuccessResponse{data=DashboardSummaryResponse} "Dashboard summary retrieved successfully"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /reports/dashboard-summary [get]
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
// @Summary      Get profit summary
// @Description  Get gross profit analytics grouped by date (Roles: admin, manager, cashier)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        start_date query string true "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string true "End Date (YYYY-MM-DD)"
// @Success      200 {object} common.SuccessResponse{data=[]ProfitSummaryResponse} "Profit summary retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /reports/profit-summary [get]
func (r *RptHandler) GetProfitSummaryHandler(c fiber.Ctx) error {
	var req SalesReportRequest
	if err := c.Bind().Query(&req); err != nil {
		r.log.Warnf("Get profit summary validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	serviceReq := &SalesReportServiceRequest{
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
// @Summary      Get product profit reports
// @Description  Get profitability metrics for each product sold (Roles: admin, manager, cashier)
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        start_date query string true "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string true "End Date (YYYY-MM-DD)"
// @Success      200 {object} common.SuccessResponse{data=[]ProductProfitResponse} "Product profit reports retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /reports/profit-products [get]
func (r *RptHandler) GetProductProfitReportsHandler(c fiber.Ctx) error {
	var req SalesReportRequest
	if err := c.Bind().Query(&req); err != nil {
		r.log.Warnf("Get product profit reports validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	serviceReq := &SalesReportServiceRequest{
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

func NewRptHandler(service IRptService, log logger.ILogger) IRptHandler {
	return &RptHandler{
		Service: service,
		log:     log,
	}
}
