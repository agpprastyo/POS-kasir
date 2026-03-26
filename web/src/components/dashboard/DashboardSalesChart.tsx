import { lazy, Suspense } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

const DashboardCharts = lazy(() => import('./DashboardCharts'))

interface DashboardSalesChartProps {
    t: any
    isLoading: boolean
    salesData: any
    paymentsData: any
    formatCurrency: (value: number) => string
    formatDate: (dateString: string) => string
}

export function DashboardSalesChart({
    t,
    isLoading,
    salesData,
    paymentsData,
    formatCurrency,
    formatDate
}: DashboardSalesChartProps) {
    return (
        <Card className="col-span-1 lg:col-span-4">
            <CardHeader>
                <CardTitle>{t('dashboard.widgets.sales_trend')}</CardTitle>
                <CardDescription>{t('reports.sales.description')}</CardDescription>
            </CardHeader>
            <CardContent className="pl-2">
                {isLoading ? (
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
    )
}
