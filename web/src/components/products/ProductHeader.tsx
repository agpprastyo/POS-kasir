import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface ProductHeaderProps {
    canCreate: boolean
    openCreateModal: () => void
    t: any
}

export function ProductHeader({ canCreate, openCreateModal, t }: ProductHeaderProps) {
    return (
        <div className="flex items-center justify-between">
            <div>
                <h1 className="text-2xl font-bold tracking-tight">{t('products.title')}</h1>
                <p className="text-muted-foreground">{t('products.description')}</p>
            </div>
            <div className="flex items-center gap-2">
                {canCreate && (
                    <Button onClick={openCreateModal}>
                        <Plus className="mr-2 h-4 w-4" /> {t('products.add_button')}
                    </Button>
                )}
            </div>
        </div>
    )
}
