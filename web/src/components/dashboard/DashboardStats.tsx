import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { DollarSign, Users, CreditCard, Activity } from 'lucide-react'
import { Skeleton } from '@/components/ui/skeleton'

interface DashboardStatsProps {
    t: any
    summary: any
    isLoading: boolean
    formatCurrency: (value: number) => string
}

export function DashboardStats({ t, summary, isLoading, formatCurrency }: DashboardStatsProps) {
    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">{t('dashboard.total_revenue')}</CardTitle>
                    <DollarSign className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    {isLoading ? (
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
                    {isLoading ? (
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
                    {isLoading ? (
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
                    {isLoading ? (
                        <Skeleton className="h-8 w-[50px]" />
                    ) : (
                        <div className="text-2xl font-bold">{summary?.total_products ?? 0}</div>
                    )}
                    <p className="text-xs text-muted-foreground">{t('dashboard.stats.all_time')}</p>
                </CardContent>
            </Card>
        </div>
    )
}
