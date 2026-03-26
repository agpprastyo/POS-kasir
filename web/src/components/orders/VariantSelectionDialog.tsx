import { Button } from '@/components/ui/button'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { formatRupiah } from '@/lib/utils'
import { Product } from '@/lib/api/query/products'
import { InternalProductsProductOptionResponse } from '@/lib/api/generated'

interface VariantSelectionDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    product: Product | null
    onSelect: (product: Product, variant?: InternalProductsProductOptionResponse) => void
    t: any
}

export function VariantSelectionDialog({
    open, onOpenChange, product, onSelect, t
}: VariantSelectionDialogProps) {
    if (!product) return null

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>{t('order.variant_dialog.title')}</DialogTitle>
                    <DialogDescription>
                        {t('order.variant_dialog.desc')} <span className="font-semibold">{product.name}</span>
                    </DialogDescription>
                </DialogHeader>

                <div className="grid grid-cols-2 gap-4 py-4">
                    {/* Base Product Option */}
                    <div
                        className="flex flex-col gap-2 p-3 border rounded-lg cursor-pointer hover:bg-accent transition-colors"
                        onClick={() => onSelect(product)}
                    >
                        <div className="aspect-square w-full rounded-md bg-muted overflow-hidden flex items-center justify-center">
                            {product.image_url ? (
                                <img src={product.image_url} alt={product.name} className="w-full h-full object-cover" />
                            ) : (
                                <span className="text-xs text-muted-foreground">{t('order.variant_dialog.original')}</span>
                            )}
                        </div>
                        <div className="flex flex-col">
                            <span className="font-medium text-sm">{t('order.variant_dialog.original')}</span>
                            <span className="text-xs text-muted-foreground">{formatRupiah(product.price || 0)}</span>
                        </div>
                    </div>

                    {product.options?.map(option => (
                        <div
                            key={option.id}
                            className="flex flex-col gap-2 p-3 border rounded-lg cursor-pointer hover:bg-accent transition-colors"
                            onClick={() => onSelect(product, option)}
                        >
                            <div className="aspect-square w-full rounded-md bg-muted overflow-hidden">
                                {option.image_url ? (
                                    <img src={option.image_url} alt={option.name} className="w-full h-full object-cover" />
                                ) : (
                                    <div className="w-full h-full flex items-center justify-center text-muted-foreground bg-muted/50">
                                        <span className="text-xs">{t('order.variant_dialog.no_image')}</span>
                                    </div>
                                )}
                            </div>
                            <div className="flex flex-col">
                                <span className="font-medium text-sm">{option.name}</span>
                                <span className="text-xs text-muted-foreground">+{formatRupiah(option.additional_price || 0)}</span>
                            </div>
                        </div>
                    ))}
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>{t('order.payment_dialog.cancel')}</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
