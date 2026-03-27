import { Loader2 } from 'lucide-react'

import { ScrollArea, ScrollBar } from '@/components/ui/scroll-area'
import { Product } from '@/lib/api/query/products'
import { ProductCard } from '../products/ProductCard'

interface ProductGridProps {
    inStockProducts: Product[]
    outOfStockProducts: Product[]
    isLoadingDetailsId: string | null
    onAddToCart: (product: Product) => void
    t: any
}

export function ProductGrid({
    inStockProducts, outOfStockProducts, isLoadingDetailsId, onAddToCart, t
}: ProductGridProps) {
    return (
        <ScrollArea className="h-full">
            <div className="p-4">
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 pb-4">
                    {inStockProducts.map(product => (
                        <div
                            key={product.id}
                            onClick={() => onAddToCart(product)}
                            className="relative cursor-pointer transition-transform active:scale-95"
                        >
                            <ProductCard
                                product={product}
                            />
                            {isLoadingDetailsId === product.id && (
                                <div className="absolute inset-0 bg-background/50 flex items-center justify-center rounded-xl z-10 backdrop-blur-[1px]">
                                    <Loader2 className="h-8 w-8 animate-spin text-primary" />
                                </div>
                            )}
                        </div>
                    ))}
                </div>

                {outOfStockProducts.length > 0 && (
                    <>
                        <div className="relative py-4">
                            <div className="absolute inset-0 flex items-center">
                                <span className="w-full border-t" />
                            </div>
                            <div className="relative flex justify-center text-sm uppercase">
                                <span className="bg-background px-2 text-muted-foreground">
                                    {t('order.out_of_stock')}
                                </span>
                            </div>
                        </div>
                        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 pb-4 opacity-60 grayscale">
                            {outOfStockProducts.map(product => (
                                <div key={product.id} className="relative cursor-not-allowed">
                                    <ProductCard product={product} />
                                </div>
                            ))}
                        </div>
                    </>
                )}

                {inStockProducts.length === 0 && outOfStockProducts.length === 0 && (
                    <div className="h-40 flex items-center justify-center text-muted-foreground">
                        {t('order.no_products')}
                    </div>
                )}
            </div>
            <ScrollBar orientation="vertical" />
        </ScrollArea>
    )
}
