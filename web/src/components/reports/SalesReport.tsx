import { Download } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid, Legend } from 'recharts'

interface SalesReportProps {
    data: any[]
    isLoading: boolean
    onExport: () => void
    formatCurrency: (value: number) => string
    formatDate: (date: string) => string
    t: any
}

export function SalesReport({
    data, isLoading, onExport, formatCurrency, formatDate, t
}: SalesReportProps) {
    return (
        <Card className="border-0 shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between">
                <div>
                    <CardTitle>{t('reports.sales.title')}</CardTitle>
                    <CardDescription>{t('reports.sales.description')}</CardDescription>
                </div>
                <Button variant="outline" size="sm" onClick={onExport} className="rounded-xl">
                    <Download className="mr-2 h-4 w-4" /> {t('reports.export_csv')}
                </Button>
            </CardHeader>
            <CardContent className="pl-2">
                {isLoading ? (
                    <Skeleton className="h-[400px] w-full rounded-xl" />
                ) : (
                    <ResponsiveContainer width="100%" height={400}>
                        <BarChart data={data}>
                            <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#E5E7EB" />
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
                            <Bar dataKey="total_sales" fill="#4F46E5" radius={[8, 8, 0, 0]} name={t('reports.sales.revenue')} />
                            <Bar dataKey="order_count" fill="#F59E0B" radius={[8, 8, 0, 0]} name={t('reports.sales.orders')} />
                        </BarChart>
                    </ResponsiveContainer>
                )}
            </CardContent>
        </Card>
    )
}
