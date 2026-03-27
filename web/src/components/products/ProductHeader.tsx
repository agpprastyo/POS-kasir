import { Plus, Package } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface ProductHeaderProps {
    canCreate: boolean
    openCreateModal: () => void
    t: any
}

export function ProductHeader({ canCreate, openCreateModal, t }: ProductHeaderProps) {
    return (
        <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center">
                    <Package className="h-5 w-5 text-primary" />
                </div>
                <div>
                    <h1 className="text-2xl font-bold tracking-tight font-heading">{t('products.title')}</h1>
                    <p className="text-sm text-muted-foreground">{t('products.description')}</p>
                </div>
            </div>
            <div className="flex items-center gap-2">
                {canCreate && (
                    <Button onClick={openCreateModal} className="rounded-xl shadow-md shadow-primary/20">
                        <Plus className="mr-2 h-4 w-4" /> {t('products.add_button')}
                    </Button>
                )}
            </div>
        </div>
    )
}
