package dto

import "time"

type DashboardSummaryResponse struct {
	TotalSales    float64 `json:"total_sales"`
	TotalOrders   int64   `json:"total_orders"`
	UniqueCashier int64   `json:"unique_cashier"`
	TotalProducts int64   `json:"total_products"`
}

type SalesReportRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type SalesReport struct {
	Date       time.Time `json:"date"`
	OrderCount int64     `json:"order_count"`
	TotalSales float64   `json:"total_sales"`
}

type ProductPerformanceResponse struct {
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
