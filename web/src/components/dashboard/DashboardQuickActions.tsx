import { Package, FileText } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Link } from '@tanstack/react-router'

interface DashboardQuickActionsProps {
    t: any
}

export function DashboardQuickActions({ t }: DashboardQuickActionsProps) {
    return (
        <Card className="col-span-1 lg:col-span-4 border-0 shadow-sm">
            <CardHeader>
                <CardTitle>{t('dashboard.widgets.quick_actions')}</CardTitle>
                <CardDescription>{t('dashboard.widgets.quick_actions_desc')}</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-2">
                <Link
                    to="/$locale/product"
                    params={{ locale: 'en' }}
                    search={{ page: 1, limit: 10, tab: 'active' }}
                    className="group h-28 flex flex-col gap-3 items-center justify-center rounded-2xl border-2 border-dashed border-border/60 bg-muted/20 transition-all duration-300 hover:border-primary/40 hover:bg-primary/5 hover:shadow-md"
                >
                    <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                        <Package className="h-6 w-6 text-primary" />
                    </div>
                    <span className="text-sm font-semibold">{t('dashboard.widgets.manage_products')}</span>
                </Link>
                <Link
                    to="/$locale/reports"
                    params={{ locale: 'en' }}
                    className="group h-28 flex flex-col gap-3 items-center justify-center rounded-2xl border-2 border-dashed border-border/60 bg-muted/20 transition-all duration-300 hover:border-amber/40 hover:bg-amber/5 hover:shadow-md"
                >
                    <div className="h-12 w-12 rounded-xl bg-amber/10 flex items-center justify-center group-hover:bg-amber/20 transition-colors">
                        <FileText className="h-6 w-6 text-amber" />
                    </div>
                    <span className="text-sm font-semibold">{t('dashboard.widgets.view_reports')}</span>
                </Link>
            </CardContent>
        </Card>
    )
}
