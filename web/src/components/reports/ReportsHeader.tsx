import { BarChart3 } from 'lucide-react'
import { DateRangePicker } from '@/components/ui/date-range-picker'

interface ReportsHeaderProps {
    dateRange: { start: string, end: string }
    onDateChange: (type: 'start' | 'end', value: string) => void
    t: any
}

export function ReportsHeader({ dateRange, onDateChange, t }: ReportsHeaderProps) {
    return (
        <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center">
                    <BarChart3 className="h-5 w-5 text-primary" />
                </div>
                <div>
                    <h1 className="text-2xl font-bold tracking-tight font-heading">{t('reports.title')}</h1>
                    <p className="text-sm text-muted-foreground">{t('reports.subtitle')}</p>
                </div>
            </div>
            <div className="flex items-center gap-2">
                <div className="grid gap-1">
                    <label className="text-sm font-medium text-muted-foreground">{t('reports.date_range') || 'Date Range'}</label>
                    <DateRangePicker
                        date={{ from: dateRange.start, to: dateRange.end }}
                        onDateChange={({ from, to }) => {
                            onDateChange('start', from)
                            onDateChange('end', to)
                        }}
                    />
                </div>
            </div>
        </div>
    )
}
