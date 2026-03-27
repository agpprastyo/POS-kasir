import { Button } from '@/components/ui/button'
import { X, ActivitySquare } from 'lucide-react'

interface ActivityLogsHeaderProps {
    t: any
    onClearFilters: () => void
}

export function ActivityLogsHeader({ t, onClearFilters }: ActivityLogsHeaderProps) {
    return (
        <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center">
                    <ActivitySquare className="h-5 w-5 text-primary" />
                </div>
                <div>
                    <h2 className="text-2xl font-bold tracking-tight font-heading">{t('activity_logs.title')}</h2>
                    <p className="text-sm text-muted-foreground">
                        {t('activity_logs.description')}
                    </p>
                </div>
            </div>
            <Button variant="outline" onClick={onClearFilters} className="rounded-xl">
                <X className="mr-2 h-4 w-4" />
                {t('common.clear_filters')}
            </Button>
        </div>
    )
}
