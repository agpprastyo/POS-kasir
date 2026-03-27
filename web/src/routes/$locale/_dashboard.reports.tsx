import { createFileRoute, redirect } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { meQueryOptions } from '@/lib/api/query/auth'
import { POSKasirInternalUserRepositoryUserRole } from '@/lib/api/generated'
import { useState } from 'react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useSalesReportQuery, useProductPerformanceQuery, usePaymentMethodPerformanceQuery, useCashierPerformanceQuery, useCancellationReportsQuery, useProfitSummaryQuery, useProductProfitReportsQuery, useLowStockReportQuery, usePromotionsReportQuery, useShiftSummaryReportQuery } from '@/lib/api/query/reports'
import { ReportsHeader } from '@/components/reports/ReportsHeader'
import { SalesReport } from '@/components/reports/SalesReport'
import { ProductsReport } from '@/components/reports/ProductsReport'
import { PerformanceReport } from '@/components/reports/PerformanceReport'
import { ProfitReport } from '@/components/reports/ProfitReport'
import { CancellationReport } from '@/components/reports/CancellationReport'
import { StockReport } from '@/components/reports/StockReport'
import { PromotionsReport } from '@/components/reports/PromotionsReport'
import { ShiftReport } from '@/components/reports/ShiftReport'
import { formatDate, formatRupiah } from '@/lib/utils'

export const Route = createFileRoute('/$locale/_dashboard/reports')({
    beforeLoad: async ({ context: { queryClient } }) => {
        const user = await queryClient.ensureQueryData(meQueryOptions())

        const allowedRoles: POSKasirInternalUserRepositoryUserRole[] = [
            POSKasirInternalUserRepositoryUserRole.UserRoleAdmin,
            POSKasirInternalUserRepositoryUserRole.UserRoleManager
        ]

        if (!user || !user.role || !allowedRoles.includes(user.role)) {
            throw redirect({
                to: '/$locale/login',
                params: { locale: 'en' },
                search: {
                    error: 'Unauthorized',
                },
            })
        }
    },
    component: ReportsPage,
})

function ReportsPage() {
    const { t } = useTranslation()

    const [dateRange, setDateRange] = useState({
        start: new Date(new Date().setDate(new Date().getDate() - 30)).toISOString().split('T')[0],
        end: new Date().toISOString().split('T')[0]
    })

    const { data: salesData, isLoading: isLoadingSales, isAllowed: canViewSales } = useSalesReportQuery(dateRange.start, dateRange.end)
    const { data: productsData, isLoading: isLoadingProducts, isAllowed: canViewProducts } = useProductPerformanceQuery(dateRange.start, dateRange.end)
    const { data: paymentsData, isLoading: isLoadingPayments, isAllowed: canViewPayments } = usePaymentMethodPerformanceQuery(dateRange.start, dateRange.end)
    const { data: cashierData, isLoading: isLoadingCashier, isAllowed: canViewCashier } = useCashierPerformanceQuery(dateRange.start, dateRange.end)
    const { data: cancellationData, isLoading: isLoadingCancellation, isAllowed: canViewCancellation } = useCancellationReportsQuery(dateRange.start, dateRange.end)
    const { data: profitSummaryData, isLoading: isLoadingProfitSummary, isAllowed: canViewProfit } = useProfitSummaryQuery(dateRange.start, dateRange.end)
    const { data: productProfitsData, isLoading: isLoadingProductProfits } = useProductProfitReportsQuery(dateRange.start, dateRange.end)

    const { data: lowStockData, isLoading: isLoadingLowStock, isAllowed: canViewLowStock } = useLowStockReportQuery(5)
    const { data: promotionsData, isLoading: isLoadingPromotions, isAllowed: canViewPromotions } = usePromotionsReportQuery(dateRange.start, dateRange.end)
    const { data: shiftData, isLoading: isLoadingShift, isAllowed: canViewShift } = useShiftSummaryReportQuery(dateRange.start, dateRange.end)

    const canViewPerformance = canViewPayments || canViewCashier

    const exportToCSV = (data: any[], filename: string, headers: string[]) => {
        if (!data || data.length === 0) return;
        const csvRows = [];
        csvRows.push(headers.join(','));
        for (const row of data) {
            const values = headers.map(header => {
                const value = row[header];
                return typeof value === 'string' ? `"${value.replace(/"/g, '""')}"` : value;
            });
            csvRows.push(values.join(','));
        }
        const csvString = csvRows.join('\n');
        const blob = new Blob([csvString], { type: 'text/csv;charset=utf-8;' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.setAttribute('href', url);
        link.setAttribute('download', `${filename}.csv`);
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    const handleDateChange = (type: 'start' | 'end', value: string) => {
        setDateRange(prev => ({
            ...prev,
            [type]: value
        }))
    }

    return (
        <div className="flex flex-col gap-4 ">
            <ReportsHeader 
                dateRange={dateRange}
                onDateChange={handleDateChange}
                t={t}
            />

            <Tabs defaultValue="sales" className="space-y-4">
                <TabsList>
                    {canViewSales && <TabsTrigger value="sales">{t('reports.tabs.sales')}</TabsTrigger>}
                    {canViewProfit && <TabsTrigger value="profit">{t('reports.tabs.profit')}</TabsTrigger>}
                    {canViewProducts && <TabsTrigger value="products">{t('reports.tabs.products')}</TabsTrigger>}
                    {canViewPerformance && <TabsTrigger value="performance">{t('reports.tabs.performance')}</TabsTrigger>}
                    {canViewCancellation && <TabsTrigger value="cancellations">{t('reports.tabs.cancellations')}</TabsTrigger>}
                    {canViewLowStock && <TabsTrigger value="low-stock">{t('reports.tabs.low_stock') || 'Low Stock'}</TabsTrigger>}
                    {canViewPromotions && <TabsTrigger value="promotions">{t('reports.tabs.promotions') || 'Promotions'}</TabsTrigger>}
                    {canViewShift && <TabsTrigger value="shifts">{t('reports.tabs.shifts') || 'Shifts'}</TabsTrigger>}
                </TabsList>

                {canViewSales && (
                    <TabsContent value="sales" className="space-y-4">
                        <SalesReport 
                            data={salesData || []}
                            isLoading={isLoadingSales}
                            onExport={() => exportToCSV(salesData || [], 'sales_report', ['date', 'total_sales', 'order_count'])}
                            formatCurrency={formatRupiah}
                            formatDate={formatDate}
                            t={t}
                        />
                    </TabsContent>
                )}

                {canViewProfit && (
                    <TabsContent value="profit" className="space-y-4">
                        <ProfitReport 
                            profitSummaryData={profitSummaryData || []}
                            productProfitsData={productProfitsData}
                            isLoadingSummary={isLoadingProfitSummary}
                            isLoadingProducts={isLoadingProductProfits}
                            formatCurrency={formatRupiah}
                            formatDate={formatDate}
                            t={t}
                        />
                    </TabsContent>
                )}

                {canViewProducts && (
                    <TabsContent value="products" className="space-y-4">
                        <ProductsReport 
                            data={productsData}
                            isLoading={isLoadingProducts}
                            formatCurrency={formatRupiah}
                            t={t}
                        />
                    </TabsContent>
                )}

                {canViewPerformance && (
                    <TabsContent value="performance" className="space-y-4">
                        <PerformanceReport 
                            paymentsData={paymentsData || []}
                            cashierData={cashierData || []}
                            isLoadingPayments={isLoadingPayments}
                            isLoadingCashier={isLoadingCashier}
                            formatCurrency={formatRupiah}
                            t={t}
                        />
                    </TabsContent>
                )}

                {canViewCancellation && (
                    <TabsContent value="cancellations" className="space-y-4">
                        <CancellationReport 
                            data={cancellationData || []}
                            isLoading={isLoadingCancellation}
                            t={t}
                        />
                    </TabsContent>
                )}

                {canViewLowStock && (
                    <TabsContent value="low-stock" className="space-y-4">
                        <StockReport 
                            data={lowStockData || []}
                            isLoading={isLoadingLowStock}
                            onExport={() => exportToCSV(lowStockData || [], 'low_stock_report', ['product_name', 'stock'])}
                            t={t}
                        />
                    </TabsContent>
                )}

                {canViewPromotions && (
                    <TabsContent value="promotions" className="space-y-4">
                        <PromotionsReport 
                            data={promotionsData || []}
                            isLoading={isLoadingPromotions}
                            onExport={() => exportToCSV(promotionsData || [], 'promotions_report', ['name', 'discount_type'])}
                            t={t}
                        />
                    </TabsContent>
                )}

                {canViewShift && (
                    <TabsContent value="shifts" className="space-y-4">
                        <ShiftReport 
                            data={shiftData || []}
                            isLoading={isLoadingShift}
                            onExport={() => exportToCSV(shiftData || [], 'shift_report', ['cashier_name', 'status'])}
                            formatCurrency={formatRupiah}
                            t={t}
                        />
                    </TabsContent>
                )}
            </Tabs>
        </div>
    )
}

