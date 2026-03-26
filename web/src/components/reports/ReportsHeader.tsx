import { DateRangePicker } from '@/components/ui/date-range-picker'

interface ReportsHeaderProps {
    dateRange: { start: string, end: string }
    onDateChange: (type: 'start' | 'end', value: string) => void
    t: any
}

export function ReportsHeader({ dateRange, onDateChange, t }: ReportsHeaderProps) {
    return (
        <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div>
                <h1 className="text-2xl font-bold tracking-tight">{t('reports.title')}</h1>
                <p className="text-muted-foreground">{t('reports.subtitle')}</p>
            </div>
            <div className="flex items-center gap-2">
                <div className="grid gap-1">
                    <label className="text-xs font-medium text-muted-foreground">{t('reports.date_range') || 'Date Range'}</label>
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
