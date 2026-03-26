import { Label } from '@/components/ui/label'
import { DateRangePicker } from '@/components/ui/date-range-picker'

interface DashboardHeaderProps {
    t: any
    username: string
    startDate: string
    endDate: string
    onDateChange: (range: { from: string; to: string }) => void
}

export function DashboardHeader({ t, username, startDate, endDate, onDateChange }: DashboardHeaderProps) {
    return (
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">
                    {t('auth.welcome_back', 'Welcome back')}, {username}!
                </h1>
                <p className="text-muted-foreground">
                    {t('dashboard.welcome_subtitle', 'Here is what is happening with your store today.')}
                </p>
            </div>

            <div className="flex items-center gap-2">
                <div className="grid gap-1">
                    <Label className="text-xs">{t('reports.date_range', 'Date Range')}</Label>
                    <DateRangePicker
                        date={{ from: startDate, to: endDate }}
                        onDateChange={onDateChange}
                    />
                </div>
            </div>
        </div>
    )
}
