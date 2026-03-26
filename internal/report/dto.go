package report

import (
	"POS-kasir/internal/common/pagination"
	"time"
)

type DashboardSummaryResponse struct {
	TotalSales    float64 `json:"total_sales"`
	TotalOrders   int64   `json:"total_orders"`
	UniqueCashier int64   `json:"unique_cashier"`
	TotalProducts int64   `json:"total_products"`
}

type SalesReportRequest struct {
	pagination.PaginationRequest
	StartDate string `json:"start_date" query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" query:"end_date" validate:"required,datetime=2006-01-02"`
	Export    string `json:"export" query:"export"`
}

type SalesReportServiceRequest struct {
	pagination.PaginationRequest
	StartDate time.Time
	EndDate   time.Time
}

type SalesReport struct {
	Date       time.Time `json:"date"`
	OrderCount int64     `json:"order_count"`
	TotalSales float64   `json:"total_sales"`
}

type ProductPerformanceResponse struct {
	Products   []ProductPerformanceRow `json:"products"`
	Pagination pagination.Pagination   `json:"pagination"`
}

type ProductPerformanceRow struct {
	ProductID     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	TotalQuantity int64   `json:"total_quantity"`
	TotalRevenue  float64 `json:"total_revenue"`
}

type PaymentMethodPerformanceResponse struct {
	PaymentMethodID   int32   `json:"payment_method_id"`
	PaymentMethodName string  `json:"payment_method_name"`
	OrderCount        int64   `json:"order_count"`
	TotalSales        float64 `json:"total_sales"`
}

type CashierPerformanceResponse struct {
	UserID     string  `json:"user_id"`
	Username   string  `json:"username"`
	OrderCount int64   `json:"order_count"`
	TotalSales float64 `json:"total_sales"`
}

type CancellationReportResponse struct {
	ReasonID        int32  `json:"reason_id"`
	Reason          string `json:"reason"`
	CancelledOrders int64  `json:"cancelled_orders"`
}

type ProfitSummaryResponse struct {
	Date         time.Time `json:"date"`
	TotalRevenue float64   `json:"total_revenue"`
	TotalCOGS    float64   `json:"total_cogs"`
	GrossProfit  float64   `json:"gross_profit"`
}

type ProductProfitResponse struct {
	Products   []ProductProfitRow    `json:"products"`
	Pagination pagination.Pagination `json:"pagination"`
}

type ProductProfitRow struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	TotalSold    int64   `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalCOGS    float64 `json:"total_cogs"`
	GrossProfit  float64 `json:"gross_profit"`
}

type LowStockRequest struct {
	Threshold int32  `json:"threshold" query:"threshold"`
	Export    string `json:"export" query:"export"`
}

type LowStockProductResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Stock       int32  `json:"stock"`
}

type PromotionPerformanceResponse struct {
	PromotionID             string  `json:"promotion_id"`
	PromotionName           string  `json:"promotion_name"`
	UsageCount              int64   `json:"usage_count"`
	TotalDiscountGiven      float64 `json:"total_discount_given"`
	TotalSalesWithPromotion float64 `json:"total_sales_with_promotion"`
}

type ShiftSummaryResponse struct {
	ShiftID         string    `json:"shift_id"`
	CashierName     string    `json:"cashier_name"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Status          string    `json:"status"`
	StartCash       int64     `json:"start_cash"`
	ActualCashEnd   *int64    `json:"actual_cash_end"`
	ExpectedCashEnd *int64    `json:"expected_cash_end"`
	CashDifference  int64     `json:"cash_difference"`
}
