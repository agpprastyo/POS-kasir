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
