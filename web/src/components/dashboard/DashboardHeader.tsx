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
        <div className="rounded-2xl bg-linear-to-r from-primary via-primary/90 to-primary/70 p-6 md:p-8 text-primary-foreground shadow-lg shadow-primary/20 relative overflow-hidden">
            {/* Decorative circles */}
            <div className="absolute -top-10 -right-10 w-40 h-40 rounded-full bg-white/10" />
            <div className="absolute -bottom-12 -right-4 w-28 h-28 rounded-full bg-white/5" />
            <div className="absolute top-4 right-32 w-8 h-8 rounded-full bg-white/10" />

            <div className="relative flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-2xl md:text-3xl font-bold tracking-tight font-heading">
                        {t('auth.welcome_back', 'Welcome back')}, {username}!
                    </h1>
                    <p className="text-primary-foreground/70 mt-1">
                        {t('dashboard.welcome_subtitle', 'Here is what is happening with your store today.')}
                    </p>
                </div>

                <div className="flex items-center gap-2">
                    <div className="grid gap-1">
                        <Label className="text-sm text-primary-foreground/60">{t('reports.date_range', 'Date Range')}</Label>
                        <DateRangePicker
                            date={{ from: startDate, to: endDate }}
                            onDateChange={onDateChange}
                        />
                    </div>
                </div>
            </div>
        </div>
    )
}
