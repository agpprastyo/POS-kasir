import { Card, CardContent } from '@/components/ui/card'
import { DollarSign, Users, CreditCard, Activity } from 'lucide-react'
import { Skeleton } from '@/components/ui/skeleton'

interface DashboardStatsProps {
    t: any
    summary: any
    isLoading: boolean
    formatCurrency: (value: number) => string
}

const statConfig = [
    { key: 'revenue', icon: DollarSign, colorClass: 'bg-primary/10 text-primary' },
    { key: 'orders', icon: CreditCard, colorClass: 'bg-amber/10 text-amber' },
    { key: 'cashiers', icon: Users, colorClass: 'bg-emerald-500/10 text-emerald-500' },
    { key: 'products', icon: Activity, colorClass: 'bg-violet-500/10 text-violet-500' },
]

export function DashboardStats({ t, summary, isLoading, formatCurrency }: DashboardStatsProps) {
    const stats = [
        { ...statConfig[0], title: t('dashboard.total_revenue'), value: formatCurrency(summary?.total_sales ?? 0), sub: t('dashboard.stats.today_sales') },
        { ...statConfig[1], title: t('dashboard.total_orders'), value: String(summary?.total_orders ?? 0), sub: t('dashboard.stats.today_orders') },
        { ...statConfig[2], title: t('dashboard.active_cashiers'), value: String(summary?.unique_cashier ?? 0), sub: t('dashboard.stats.active_today') },
        { ...statConfig[3], title: t('dashboard.total_products'), value: String(summary?.total_products ?? 0), sub: t('dashboard.stats.all_time') },
    ]

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {stats.map(({ key, icon: Icon, colorClass, title, value, sub }) => (
                <Card key={key} className="border-0 shadow-sm hover:shadow-md transition-shadow duration-300">
                    <CardContent className="p-5">
                        <div className="flex items-center justify-between mb-3">
                            <span className="text-sm font-medium text-muted-foreground">{title}</span>
                            <div className={`h-9 w-9 rounded-xl flex items-center justify-center ${colorClass}`}>
                                <Icon className="h-4 w-4" />
                            </div>
                        </div>
                        {isLoading ? (
                            <Skeleton className="h-8 w-24 rounded-lg" />
                        ) : (
                            <div className="text-2xl font-bold font-heading tracking-tight">{value}</div>
                        )}
                        <p className="text-sm text-muted-foreground mt-1">{sub}</p>
                    </CardContent>
                </Card>
            ))}
        </div>
    )
}
