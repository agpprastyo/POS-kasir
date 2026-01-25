import { Badge } from "@/components/ui/badge"
import { formatRupiah } from "@/lib/utils"
import { Package, RotateCcw } from "lucide-react"
import { ProductActions } from "@/components/ProductActions"
import { Button } from "@/components/ui/button"
import { Product } from "@/lib/api/query/products"
import { cn } from "@/lib/utils"
import { useTranslation } from 'react-i18next';

interface ProductCardProps {
    product: Product
    onEdit?: (product: Product) => void
    onRestore?: (product: Product) => void
}

export function ProductCard({ product, onEdit, onRestore }: ProductCardProps) {
    const { t } = useTranslation();
    return (
        <div className={cn(
            "group relative rounded-lg bg-card border border-border/40 text-card-foreground transition-all duration-300 hover:border-border/80 hover:shadow-sm overflow-hidden",
            onRestore && "opacity-75 border-dashed bg-muted/30"
        )}>

            <div className="aspect-square w-full relative bg-secondary/20 overflow-hidden">
                {product.image_url ? (
                    <img
                        src={product.image_url}
                        alt={product.name}
                        className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
                    />
                ) : (
                    <div className="flex h-full w-full items-center justify-center text-muted-foreground/20">
                        <Package className="h-10 w-10" strokeWidth={1} />
                    </div>
                )}


                {onEdit && (
                    <div className="absolute top-2 right-2 z-10">
                        <div className="rounded-full bg-background/90 backdrop-blur-sm p-1 opacity-0 translate-y-1 group-hover:opacity-100 group-hover:translate-y-0 transition-all duration-200">
                            <ProductActions product={product} onEdit={() => onEdit(product)} />
                        </div>
                    </div>
                )}

                {onRestore && (
                    <div className="absolute top-2 right-2 z-10">
                        <Button
                            size="icon"
                            variant="secondary"
                            className="h-8 w-8 rounded-full bg-background/90 backdrop-blur-sm opacity-0 translate-y-1 group-hover:opacity-100 group-hover:translate-y-0 transition-all duration-200 hover:bg-green-100 hover:text-green-600 shadow-sm"
                            onClick={(e) => {
                                e.stopPropagation();
                                onRestore(product);
                            }}
                        >
                            <RotateCcw className="h-4 w-4" />
                        </Button>
                    </div>
                )}


                <div className="absolute top-2 left-2 z-10">
                    {product.stock !== undefined && product.stock <= 5 && (
                        <Badge
                            variant={product.stock === 0 ? "destructive" : "secondary"}
                            className={cn(
                                "px-1.5 py-0 text-[10px] font-medium border-0 shadow-none backdrop-blur-sm h-5",
                                product.stock !== 0 && "bg-background/90 text-foreground/70"
                            )}
                        >
                            {product.stock === 0 ? t('products.card.out_of_stock') : t('products.card.stock_left', { count: product.stock })}
                        </Badge>
                    )}
                </div>
            </div>


            <div className="p-3">
                <div className="flex justify-between items-start gap-3">


                    <div className="flex flex-col gap-0.5 min-w-0 flex-1">
                        <h3
                            className="text-sm font-medium text-foreground truncate"
                            title={product.name}
                        >
                            {product.name}
                        </h3>
                        <p className="text-[11px] text-muted-foreground font-normal truncate">
                            {product.category_name || t('products.card.uncategorized')}
                        </p>
                    </div>


                    <div className="shrink-0">
                        <p className="text-sm font-medium text-foreground/90">
                            {formatRupiah(product.price || 0)}
                        </p>
                    </div>

                </div>
            </div>
        </div>
    )
}