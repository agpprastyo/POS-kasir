import { Download } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

interface PromotionsReportProps {
    data: any[]
    isLoading: boolean
    onExport: () => void
    t: any
}

export function PromotionsReport({
    data, isLoading, onExport, t
}: PromotionsReportProps) {
    return (
        <Card className="border-0 shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between">
                <div>
                    <CardTitle>{t('reports.promotions.title') || 'Active Promotions'}</CardTitle>
                    <CardDescription>{t('reports.promotions.description') || 'Check promotions efficacy.'}</CardDescription>
                </div>
                <Button variant="outline" size="sm" onClick={onExport} className="rounded-xl">
                    <Download className="mr-2 h-4 w-4" /> {t('reports.export_csv')}
                </Button>
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
                                <TableHead>{t('reports.promotions.name') || 'Promo Name'}</TableHead>
                                <TableHead>{t('reports.promotions.discount') || 'Discount'}</TableHead>
                                <TableHead>{t('reports.promotions.type') || 'Type'}</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {(data || []).map((promo: any, idx: number) => (
                                <TableRow key={idx} className="hover:bg-muted/50">
                                    <TableCell className="font-medium">{promo.name}</TableCell>
                                    <TableCell>{promo.discount_value}</TableCell>
                                    <TableCell>{promo.discount_type}</TableCell>
                                </TableRow>
                            ))}
                            {(!data || data.length === 0) && (
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
