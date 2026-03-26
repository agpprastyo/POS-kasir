import { createFileRoute } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { useDashboardSummaryQuery, useSalesReportQuery, useProductPerformanceQuery, usePaymentMethodPerformanceQuery } from '@/lib/api/query/reports'
import { useQueryClient } from '@tanstack/react-query'
import { meQueryOptions } from '@/lib/api/query/auth'
import { useState } from 'react'
import { DashboardHeader } from "@/components/dashboard/DashboardHeader"
import { DashboardStats } from "@/components/dashboard/DashboardStats"
import { DashboardSalesChart } from "@/components/dashboard/DashboardSalesChart"
import { DashboardTopProducts } from "@/components/dashboard/DashboardTopProducts"
import { DashboardPaymentChart } from "@/components/dashboard/DashboardPaymentChart"
import { DashboardQuickActions } from "@/components/dashboard/DashboardQuickActions"
import { formatDate, formatRupiah } from '@/lib/utils'


export const Route = createFileRoute('/$locale/_dashboard/')(({
    component: DashboardIndex,
    loader: ({ context: { queryClient } }: any) => queryClient.ensureQueryData(meQueryOptions()),
} as any))

function DashboardIndex() {
    const { t } = useTranslation()
    const queryClient = useQueryClient()
    const user = queryClient.getQueryData(meQueryOptions().queryKey) as any

    const [endDate, setEndDate] = useState<string>(new Date().toISOString().split('T')[0])
    const [startDate, setStartDate] = useState<string>(
        new Date(new Date().setDate(new Date().getDate() - 30)).toISOString().split('T')[0]
    )

    const { data: summary, isLoading: isLoadingSummary } = useDashboardSummaryQuery(startDate, endDate)
    const { data: salesData, isLoading: isLoadingSales } = useSalesReportQuery(startDate, endDate)
    const { data: productsData, isLoading: isLoadingProducts } = useProductPerformanceQuery(startDate, endDate)
    const { data: paymentsData, isLoading: isLoadingPayments } = usePaymentMethodPerformanceQuery(startDate, endDate)

    const topProducts = productsData?.products?.slice(0, 5) || []

    return (
        <div className="flex flex-col gap-8">
            <DashboardHeader
                t={t}
                username={user?.username || 'User'}
                startDate={startDate}
                endDate={endDate}
                onDateChange={({ from, to }) => {
                    setStartDate(from)
                    setEndDate(to)
                }}
            />

            <DashboardStats
                t={t}
                summary={summary}
                isLoading={isLoadingSummary}
                formatCurrency={formatRupiah}
            />

            <div className="grid gap-4 grid-cols-1 lg:grid-cols-7">
                <DashboardSalesChart
                    t={t}
                    isLoading={isLoadingSales}
                    salesData={salesData}
                    paymentsData={paymentsData}
                    formatCurrency={formatRupiah}
                    formatDate={formatDate}
                />

                <DashboardTopProducts
                    t={t}
                    isLoading={isLoadingProducts}
                    topProducts={topProducts}
                    formatCurrency={formatRupiah}
                />
            </div>

            <div className="grid gap-4 grid-cols-1 lg:grid-cols-7">
                <DashboardPaymentChart
                    t={t}
                    isLoading={isLoadingPayments}
                    salesData={salesData}
                    paymentsData={paymentsData}
                    formatCurrency={formatRupiah}
                    formatDate={formatDate}
                />

                <DashboardQuickActions t={t} />
            </div>
        </div>
    )
}