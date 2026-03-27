import { lazy, Suspense } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

const DashboardCharts = lazy(() => import('./DashboardCharts'))

interface DashboardPaymentChartProps {
    t: any
    isLoading: boolean
    salesData: any
    paymentsData: any
    formatCurrency: (value: number) => string
    formatDate: (dateString: string) => string
}

export function DashboardPaymentChart({
    t,
    isLoading,
    salesData,
    paymentsData,
    formatCurrency,
    formatDate
}: DashboardPaymentChartProps) {
    return (
        <Card className="col-span-1 lg:col-span-3 border-0 shadow-sm">
            <CardHeader>
                <CardTitle>{t('dashboard.widgets.payment_distribution')}</CardTitle>
            </CardHeader>
            <CardContent>
                {isLoading ? (
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
    )
}
