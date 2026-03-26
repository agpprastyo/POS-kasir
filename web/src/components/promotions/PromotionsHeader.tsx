import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface PromotionsHeaderProps {
    t: any
    onCreateClick: () => void
    canCreate: boolean
}

export function PromotionsHeader({ t, onCreateClick, canCreate }: PromotionsHeaderProps) {
    return (
        <div className="flex items-center justify-between">
            <div>
                <h1 className="text-2xl font-bold tracking-tight">{t('promotions.title')}</h1>
                <p className="text-muted-foreground">{t('promotions.description')}</p>
            </div>
            <div className="flex items-center gap-2">
                {canCreate && (
                    <Button onClick={onCreateClick}>
                        <Plus className="mr-2 h-4 w-4" /> {t('promotions.add_button')}
                    </Button>
                )}
            </div>
        </div>
    )
}
