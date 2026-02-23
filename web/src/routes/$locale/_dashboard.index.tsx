import { lazy, Suspense } from 'react'
import { createFileRoute, Link } from '@tanstack/react-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { DollarSign, Users, CreditCard, Activity, Package, FileText } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useDashboardSummaryQuery, useSalesReportQuery, useProductPerformanceQuery, usePaymentMethodPerformanceQuery } from '@/lib/api/query/reports'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { useQueryClient } from '@tanstack/react-query'
import { meQueryOptions } from '@/lib/api/query/auth'

// Lazy load Recharts — tidak masuk initial bundle
const DashboardCharts = lazy(() => import('@/components/dashboard/DashboardCharts'))

export const Route = createFileRoute('/$locale/_dashboard/')(({
    component: DashboardIndex,
    loader: ({ context: { queryClient } }: any) => queryClient.ensureQueryData(meQueryOptions()),
} as any))

function DashboardIndex() {
    const { t } = useTranslation()
    const { data: summary, isLoading: isLoadingSummary } = useDashboardSummaryQuery()
    const queryClient = useQueryClient()
    const user = queryClient.getQueryData(meQueryOptions().queryKey)

    const endDate = new Date().toISOString().split('T')[0]
    const startDate = new Date(new Date().setDate(new Date().getDate() - 30)).toISOString().split('T')[0]

    const { data: salesData, isLoading: isLoadingSales } = useSalesReportQuery(startDate, endDate)
    const { data: productsData, isLoading: isLoadingProducts } = useProductPerformanceQuery(startDate, endDate)
    const { data: paymentsData, isLoading: isLoadingPayments } = usePaymentMethodPerformanceQuery(startDate, endDate)

    const formatCurrency = (value: number) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0
        }).format(value)
    }

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('id-ID', {
            day: 'numeric',
            month: 'short'
        })
    }

    // Get top 5 products
    const topProducts = productsData?.slice(0, 5) || []

    return (
        <div className="flex flex-col gap-8">
            {/* Welcome Header */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">
                        {t('auth.welcome_back')}, {user?.username || 'User'}!
                    </h1>
                    <p className="text-muted-foreground">
                        {t('dashboard.welcome_subtitle')}
                    </p>
                </div>

            </div>

            {/* Stat Cards */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">{t('dashboard.total_revenue')}</CardTitle>
                        <DollarSign className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        {isLoadingSummary ? (
                            <Skeleton className="h-8 w-[100px]" />
                        ) : (
                            <div className="text-2xl font-bold">
                                {formatCurrency(summary?.total_sales ?? 0)}
                            </div>
                        )}
                        <p className="text-xs text-muted-foreground">{t('dashboard.stats.today_sales')}</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">{t('dashboard.total_orders')}</CardTitle>
                        <CreditCard className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        {isLoadingSummary ? (
                            <Skeleton className="h-8 w-[50px]" />
                        ) : (
                            <div className="text-2xl font-bold">{summary?.total_orders ?? 0}</div>
                        )}
                        <p className="text-xs text-muted-foreground">{t('dashboard.stats.today_orders')}</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">{t('dashboard.active_cashiers')}</CardTitle>
                        <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        {isLoadingSummary ? (
                            <Skeleton className="h-8 w-[50px]" />
                        ) : (
                            <div className="text-2xl font-bold">{summary?.unique_cashier ?? 0}</div>
                        )}
                        <p className="text-xs text-muted-foreground">{t('dashboard.stats.active_today')}</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">{t('dashboard.total_products')}</CardTitle>
                        <Activity className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        {isLoadingSummary ? (
                            <Skeleton className="h-8 w-[50px]" />
                        ) : (
                            <div className="text-2xl font-bold">{summary?.total_products ?? 0}</div>
                        )}
                        <p className="text-xs text-muted-foreground">{t('dashboard.stats.all_time')}</p>
                    </CardContent>
                </Card>
            </div>

            {/* Main Charts Area */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                {/* Sales Trend Chart */}
                <Card className="col-span-4">
                    <CardHeader>
                        <CardTitle>{t('dashboard.widgets.sales_trend')}</CardTitle>
                        <CardDescription>{t('reports.sales.description')}</CardDescription>
                    </CardHeader>
                    <CardContent className="pl-2">
                        {isLoadingSales ? (
                            <Skeleton className="h-[350px] w-full" />
                        ) : (
                            <Suspense fallback={<Skeleton className="h-[350px] w-full" />}>
                                <DashboardCharts
                                    salesData={salesData}
                                    paymentsData={paymentsData}
                                    formatCurrency={formatCurrency}
                                    formatDate={formatDate}
                                    chartType="bar"
                                />
                            </Suspense>
                        )}
                    </CardContent>
                </Card>

                {/* Top Products */}
                <Card className="col-span-3 flex flex-col">
                    <CardHeader>
                        <CardTitle>{t('dashboard.widgets.top_products')}</CardTitle>
                        <CardDescription>{t('reports.products.description')}</CardDescription>
                    </CardHeader>
                    <CardContent className="flex-1">
                        {isLoadingProducts ? (
                            <div className="space-y-4">
                                <Skeleton className="h-12 w-full" />
                                <Skeleton className="h-12 w-full" />
                                <Skeleton className="h-12 w-full" />
                            </div>
                        ) : (
                            <div className="space-y-8">
                                {topProducts.map((product, index) => (
                                    <div key={index} className="flex items-center">
                                        <div className="flex h-9 w-9 items-center justify-center rounded-full border border-muted bg-primary/10 text-sm font-medium text-primary">
                                            {index + 1}
                                        </div>
                                        <div className="ml-4 space-y-1">
                                            <p className="text-sm font-medium leading-none">{product.product_name}</p>
                                            <p className="text-xs text-muted-foreground">
                                                {product.total_quantity} sold
                                            </p>
                                        </div>
                                        <div className="ml-auto font-medium">
                                            {formatCurrency(product.total_revenue ?? 0)}
                                        </div>
                                    </div>
                                ))}
                                {topProducts.length === 0 && (
                                    <div className="text-center text-sm text-muted-foreground py-8">
                                        {t('common.no_data')}
                                    </div>
                                )}
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>

            {/* Lower Charts Area */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                <Card className="col-span-3">
                    <CardHeader>
                        <CardTitle>{t('dashboard.widgets.payment_distribution')}</CardTitle>
                    </CardHeader>
                    <CardContent>
                        {isLoadingPayments ? (
                            <Skeleton className="h-[300px] w-full" />
                        ) : (
                            <div className="h-[300px] w-full">
                                <Suspense fallback={<Skeleton className="h-[300px] w-full" />}>
                                    <DashboardCharts
                                        salesData={salesData}
                                        paymentsData={paymentsData}
                                        formatCurrency={formatCurrency}
                                        formatDate={formatDate}
                                        chartType="pie"
                                    />
                                </Suspense>
                            </div>
                        )}
                    </CardContent>
                </Card>

                <Card className="col-span-4">
                    <CardHeader>
                        <CardTitle>{t('dashboard.widgets.quick_actions')}</CardTitle>
                        <CardDescription>{t('dashboard.widgets.quick_actions_desc')}</CardDescription>
                    </CardHeader>
                    <CardContent className="grid gap-4 md:grid-cols-2">
                        <Button variant="outline" className="h-24 flex flex-col gap-2 items-center justify-center text-lg hover:border-primary hover:bg-primary/5" asChild>
                            <Link to="/$locale/product" params={{ locale: 'en' }} search={{ page: 1, limit: 10, tab: 'active' }}>
                                <Package className="h-8 w-8 text-primary" />
                                {t('dashboard.widgets.manage_products')}
                            </Link>
                        </Button>
                        <Button variant="outline" className="h-24 flex flex-col gap-2 items-center justify-center text-lg hover:border-primary hover:bg-primary/5" asChild>
                            <Link to="/$locale/reports" params={{ locale: 'en' }}>
                                <FileText className="h-8 w-8 text-primary" />
                                {t('dashboard.widgets.view_reports')}
                            </Link>
                        </Button>
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}