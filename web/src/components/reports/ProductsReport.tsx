import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

interface ProductsReportProps {
    data: any
    isLoading: boolean
    formatCurrency: (value: number) => string
    t: any
}

export function ProductsReport({
    data, isLoading, formatCurrency, t
}: ProductsReportProps) {
    return (
        <Card className="border-0 shadow-sm">
            <CardHeader>
                <CardTitle>{t('reports.products.title')}</CardTitle>
                <CardDescription>{t('reports.products.description')}</CardDescription>
            </CardHeader>
            <CardContent>
                {isLoading ? (
                    <div className="space-y-2">
                        <Skeleton className="h-10 w-full rounded-lg" />
                        <Skeleton className="h-10 w-full rounded-lg" />
                    </div>
                ) : (
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>{t('reports.products.product')}</TableHead>
                                <TableHead className="text-right">{t('reports.products.quantity')}</TableHead>
                                <TableHead className="text-right">{t('reports.products.revenue')}</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {(data?.products || []).map((product: any) => (
                                <TableRow key={product.product_id} className="hover:bg-muted/50">
                                    <TableCell className="font-medium">{product.product_name}</TableCell>
                                    <TableCell className="text-right">{product.total_quantity ?? 0}</TableCell>
                                    <TableCell className="text-right">{formatCurrency(product.total_revenue ?? 0)}</TableCell>
                                </TableRow>
                            ))}
                            {(!data?.products || data.products.length === 0) && (
                                <TableRow>
                                    <TableCell colSpan={3} className="text-center text-muted-foreground">{t('common.no_data')}</TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                )}
            </CardContent>
        </Card>
    )
}
