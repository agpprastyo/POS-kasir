import { ProductCard } from './ProductCard'

interface ProductGridViewProps {
    products: any[]
    emptyMessage: string
    onEdit?: (product: any) => void
    onRestore?: (product: any) => void
    t: any
}

export function ProductGridView({
    products, emptyMessage, onEdit, onRestore
}: ProductGridViewProps) {
    if (products.length === 0) {
        return (
            <div className="col-span-full h-24 flex items-center justify-center text-muted-foreground border rounded-md border-dashed">
                {emptyMessage}
            </div>
        )
    }

    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
            {products.map((product: any) => (
                <ProductCard
                    key={product.id}
                    product={product}
                    onEdit={onEdit}
                    onRestore={onRestore}
                />
            ))}
        </div>
    )
}
