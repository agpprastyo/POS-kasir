import { Package } from 'lucide-react'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ProductActions } from './ProductActions'
import { formatRupiah } from '@/lib/utils'

interface ProductTableProps {
    products: any[]
    emptyMessage: string
    isTrash?: boolean
    canRestore?: boolean
    onRestore?: (product: any) => void
    onEdit?: (product: any) => void
    t: any
    hasActions?: boolean
}

export function ProductTable({
    products, emptyMessage, isTrash, canRestore, onRestore, onEdit, t, hasActions = true
}: ProductTableProps) {
    return (
        <div className="rounded-md border bg-card">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-[80px] hidden sm:table-cell">{t('products.table.image')}</TableHead>
                        <TableHead>{t('products.table.name')}</TableHead>
                        <TableHead className="hidden md:table-cell">{t('products.table.category')}</TableHead>
                        <TableHead>{t('products.table.price')}</TableHead>
                        <TableHead>{t('products.table.stock')}</TableHead>
                        <TableHead className="text-right">{t('products.table.actions')}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {products.length === 0 ? (
                        <TableRow>
                            <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                                {emptyMessage}
                            </TableCell>
                        </TableRow>
                    ) : (
                        products.map((product: any) => (
                            <TableRow key={product.id} className={isTrash ? "opacity-70 bg-muted/30" : ""}>
                                <TableCell className="hidden sm:table-cell">
                                    <Avatar className={`h-10 w-10 rounded-md border ${isTrash ? "grayscale" : ""}`}>
                                        <AvatarImage src={product.image_url} alt={product.name}
                                            className="object-cover" />
                                        <AvatarFallback className="rounded-md">
                                            <Package className="h-5 w-5 text-muted-foreground" />
                                        </AvatarFallback>
                                    </Avatar>
                                </TableCell>
                                <TableCell className={`font-medium ${isTrash ? "text-muted-foreground" : ""}`}>
                                    {product.name}
                                </TableCell>
                                <TableCell className="hidden md:table-cell">
                                    <div className="flex flex-wrap gap-1">
                                        {product.categories?.map((cat: any) => (
                                            <Badge key={cat.id} variant="outline" className={`text-sm ${isTrash ? "opacity-50" : ""}`}>{cat.name}</Badge>
                                        )) || '-'}
                                    </div>
                                </TableCell>
                                <TableCell className={isTrash ? "text-muted-foreground" : ""}>
                                    {formatRupiah(product.price || 0)}
                                </TableCell>
                                <TableCell>
                                    {isTrash ? (
                                        <span className="text-muted-foreground">{product.stock}</span>
                                    ) : (
                                        <Badge
                                            variant={product.stock && product.stock > 0 ? 'secondary' : 'destructive'}>
                                            {product.stock}
                                        </Badge>
                                    )}
                                </TableCell>
                                <TableCell className="text-right">
                                    {isTrash ? (
                                        canRestore && (
                                            <Button
                                                size="sm"
                                                variant="outline"
                                                onClick={() => onRestore?.(product)}
                                            >
                                                {t('products.actions.restore')}
                                            </Button>
                                        )
                                    ) : (
                                        hasActions ? <ProductActions product={product} onEdit={() => onEdit?.(product)} /> : null
                                    )}
                                </TableCell>
                            </TableRow>
                        ))
                    )}
                </TableBody>
            </Table>
        </div>
    )
}
