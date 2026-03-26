import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid, Legend } from 'recharts'

interface ProfitReportProps {
    profitSummaryData: any[]
    productProfitsData: any
    isLoadingSummary: boolean
    isLoadingProducts: boolean
    formatCurrency: (value: number) => string
    formatDate: (date: string) => string
    t: any
}

export function ProfitReport({
    profitSummaryData, productProfitsData, isLoadingSummary, isLoadingProducts, formatCurrency, formatDate, t
}: ProfitReportProps) {
    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
            <Card className="col-span-4">
                <CardHeader>
                    <CardTitle>{t('reports.profit.title')}</CardTitle>
                    <CardDescription>{t('reports.profit.description')}</CardDescription>
                </CardHeader>
                <CardContent className="pl-2">
                    {isLoadingSummary ? (
                        <Skeleton className="h-[400px] w-full" />
                    ) : (
                        <ResponsiveContainer width="100%" height={400}>
                            <BarChart data={profitSummaryData}>
                                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                                <XAxis
                                    dataKey="date"
                                    stroke="#888888"
                                    fontSize={12}
                                    tickLine={false}
                                    axisLine={false}
                                    tickFormatter={(value) => formatDate(value)}
                                />
                                <YAxis
                                    stroke="#888888"
                                    fontSize={12}
                                    tickLine={false}
                                    axisLine={false}
                                    tickFormatter={(value) => `Rp${(value / 1000).toLocaleString()}k`}
                                />
                                <Tooltip
                                    formatter={(value: any) => formatCurrency(Number(value || 0))}
                                    labelFormatter={(label) => formatDate(label)}
                                />
                                <Legend />
                                <Bar dataKey="total_revenue" fill="#adfa1d" radius={[4, 4, 0, 0]} name={t('reports.profit.revenue')} />
                                <Bar dataKey="total_cogs" fill="#ef4444" radius={[4, 4, 0, 0]} name={t('reports.profit.cogs')} />
                                <Bar dataKey="gross_profit" fill="#3b82f6" radius={[4, 4, 0, 0]} name={t('reports.profit.gross_profit')} />
                            </BarChart>
                        </ResponsiveContainer>
                    )}
                </CardContent>
            </Card>
            <Card className="col-span-3">
                <CardHeader>
                    <CardTitle>{t('reports.profit.products')}</CardTitle>
                    <CardDescription>{t('reports.profit.products_desc')}</CardDescription>
                </CardHeader>
                <CardContent>
                    {isLoadingProducts ? (
                        <div className="space-y-2">
                            <Skeleton className="h-10 w-full" />
                            <Skeleton className="h-10 w-full" />
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>{t('reports.products.product')}</TableHead>
                                    <TableHead className="text-right">{t('reports.profit.profit')}</TableHead>
                                    <TableHead className="text-right">%</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {(productProfitsData?.products || []).map((product: any) => (
                                    <TableRow key={product.product_id}>
                                        <TableCell className="font-medium">
                                            <div className="flex flex-col">
                                                <span>{product.product_name}</span>
                                                <span className="text-xs text-muted-foreground">{product.total_sold ?? 0} {t('common.sold')}</span>
                                            </div>
                                        </TableCell>
                                        <TableCell className="text-right">
                                            <div className="flex flex-col items-end">
                                                <span className="font-medium">{formatCurrency(product.gross_profit ?? 0)}</span>
                                                <span className="text-xs text-muted-foreground">{formatCurrency(product.total_revenue ?? 0)} {t('common.rev')}</span>
                                            </div>
                                        </TableCell>
                                        <TableCell className="text-right">
                                            {product.total_revenue && product.total_revenue > 0
                                                ? `${(((product.gross_profit ?? 0) / product.total_revenue) * 100).toFixed(1)}%`
                                                : '0%'}
                                        </TableCell>
                                    </TableRow>
                                ))}
                                {(!productProfitsData?.products || productProfitsData.products.length === 0) && (
                                    <TableRow>
                                        <TableCell colSpan={3} className="text-center text-muted-foreground">{t('common.no_data')}</TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}
