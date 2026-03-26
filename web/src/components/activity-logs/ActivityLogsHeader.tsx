import { Button } from '@/components/ui/button'
import { X } from 'lucide-react'

interface ActivityLogsHeaderProps {
    t: any
    onClearFilters: () => void
}

export function ActivityLogsHeader({ t, onClearFilters }: ActivityLogsHeaderProps) {
    return (
        <div className="flex items-center justify-between">
            <div>
                <h2 className="text-3xl font-bold tracking-tight">{t('activity_logs.title')}</h2>
                <p className="text-muted-foreground">
                    {t('activity_logs.description')}
                </p>
            </div>
            <Button variant="outline" onClick={onClearFilters}>
                <X className="mr-2 h-4 w-4" />
                {t('common.clear_filters')}
            </Button>
        </div>
    )
}
