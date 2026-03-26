import { Package, FileText } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'

interface DashboardQuickActionsProps {
    t: any
}

export function DashboardQuickActions({ t }: DashboardQuickActionsProps) {
    return (
        <Card className="col-span-1 lg:col-span-4">
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
    )
}
