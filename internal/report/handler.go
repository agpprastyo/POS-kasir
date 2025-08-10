package report

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IRptHandler interface {
	GetDashboardSummaryHandler(c *fiber.Ctx) error
	GetSalesReportsHandler(c *fiber.Ctx) error
	GetProductPerformanceHandler(c *fiber.Ctx) error
	GetPaymentMethodPerformanceHandler(c *fiber.Ctx) error
	GetCashierPerformanceHandler(c *fiber.Ctx) error
	GetCancellationReportsHandler(c *fiber.Ctx) error
}

type RptHandler struct {
	Service IRptService
	log     logger.ILogger
}

func (r *RptHandler) GetSalesReportsHandler(c *fiber.Ctx) error {
	type queryParams struct {
		StartDate string `form:"start_date"`
		EndDate   string `form:"end_date"`
	}
	var q queryParams
	_ = c.QueryParser(&q) // optional, can remove if not working

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	r.log.Infof("Direct query: start_date=%q, end_date=%q", startDateStr, endDateStr)

	const layout = "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid start_date format, use YYYY-MM-DD",
		})
	}
	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid end_date format, use YYYY-MM-DD",
		})
	}

	// Use DTO with time.Time fields
	serviceReq := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	salesReports, err := r.Service.GetSalesReports(c.Context(), serviceReq)
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

func (r *RptHandler) GetProductPerformanceHandler(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r *RptHandler) GetPaymentMethodPerformanceHandler(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r *RptHandler) GetCashierPerformanceHandler(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r *RptHandler) GetCancellationReportsHandler(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (r *RptHandler) GetDashboardSummaryHandler(c *fiber.Ctx) error {

	summary, err := r.Service.GetDashboardSummary(c.Context())
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

func NewRptHandler(service IRptService, log logger.ILogger) IRptHandler {
	return &RptHandler{
		Service: service,
		log:     log,
	}
}
