import { Download } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

interface StockReportProps {
    data: any[]
    isLoading: boolean
    onExport: () => void
    t: any
}

export function StockReport({
    data, isLoading, onExport, t
}: StockReportProps) {
    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between">
                <div>
                    <CardTitle>{t('reports.low_stock.title') || 'Low Stock Products'}</CardTitle>
                    <CardDescription>{t('reports.low_stock.description') || 'Products with stock below minimum threshold.'}</CardDescription>
                </div>
                <Button variant="outline" size="sm" onClick={onExport}>
                    <Download className="mr-2 h-4 w-4" /> {t('reports.export_csv')}
                </Button>
            </CardHeader>
            <CardContent>
                {isLoading ? (
                    <div className="space-y-2">
                        <Skeleton className="h-10 w-full" />
                        <Skeleton className="h-10 w-full" />
                    </div>
                ) : (
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>{t('reports.products.product')}</TableHead>
                                <TableHead className="text-right">{t('reports.products.stock') || 'Stock'}</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {(data || []).map((product: any, idx: number) => (
                                <TableRow key={idx}>
                                    <TableCell className="font-medium">{product.product_name}</TableCell>
                                    <TableCell className="text-right text-red-500 font-bold">{product.stock ?? 0}</TableCell>
                                </TableRow>
                            ))}
                            {(!data || data.length === 0) && (
                                <TableRow>
                                    <TableCell colSpan={2} className="text-center text-muted-foreground">{t('common.no_data')}</TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                )}
            </CardContent>
        </Card>
    )
}
