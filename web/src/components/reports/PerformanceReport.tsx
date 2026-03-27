import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { ResponsiveContainer, PieChart, Pie, Cell, Tooltip, Legend } from 'recharts'

const COLORS = ['#4F46E5', '#F59E0B', '#7C3AED', '#10B981', '#EC4899', '#06B6D4'];

interface PerformanceReportProps {
    paymentsData: any[]
    cashierData: any[]
    isLoadingPayments: boolean
    isLoadingCashier: boolean
    formatCurrency: (value: number) => string
    t: any
}

export function PerformanceReport({
    paymentsData, cashierData, isLoadingPayments, isLoadingCashier, formatCurrency, t
}: PerformanceReportProps) {
    return (
        <div className="grid gap-4 grid-cols-1 lg:grid-cols-7">
            <Card className="col-span-1 lg:col-span-4 border-0 shadow-sm">
                <CardHeader>
                    <CardTitle>{t('reports.performance.payment_methods')}</CardTitle>
                </CardHeader>
                <CardContent className="pl-2">
                    {isLoadingPayments ? (
                        <Skeleton className="h-[300px] w-full rounded-xl" />
                    ) : (
                        <div className="h-[300px] w-full">
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={paymentsData}
                                        cx="50%"
                                        cy="50%"
                                        labelLine={false}
                                        label={({ name, percent }: any) => `${name} ${((percent || 0) * 100).toFixed(0)}%`}
                                        outerRadius={100}
                                        fill="#4F46E5"
                                        dataKey="order_count"
                                        nameKey="payment_method_name"
                                    >
                                        {(paymentsData || []).map((_entry, index) => (
                                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip />
                                    <Legend />
                                </PieChart>
                            </ResponsiveContainer>
                        </div>
                    )}
                </CardContent>
            </Card>
            <Card className="col-span-3 border-0 shadow-sm">
                <CardHeader>
                    <CardTitle>{t('reports.performance.cashier')}</CardTitle>
                </CardHeader>
                <CardContent>
                    {isLoadingCashier ? (
                        <div className="space-y-2">
                            <Skeleton className="h-10 w-full rounded-lg" />
                            <Skeleton className="h-10 w-full rounded-lg" />
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>{t('reports.performance.staff')}</TableHead>
                                    <TableHead className="text-right">{t('reports.performance.orders')}</TableHead>
                                    <TableHead className="text-right">{t('reports.performance.sales')}</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {(cashierData || []).map((cashier: any) => (
                                    <TableRow key={cashier.user_id} className="hover:bg-muted/50">
                                        <TableCell className="font-medium">
                                            <div className="flex items-center gap-2">
                                                <Avatar className="h-7 w-7">
                                                    <AvatarFallback className="text-sm bg-primary/10 text-primary">{cashier.username?.substring(0, 2).toUpperCase()}</AvatarFallback>
                                                </Avatar>
                                                {cashier.username}
                                            </div>
                                        </TableCell>
                                        <TableCell className="text-right">{cashier.order_count ?? 0}</TableCell>
                                        <TableCell className="text-right">{formatCurrency(cashier.total_sales ?? 0)}</TableCell>
                                    </TableRow>
                                ))}
                                {(!cashierData || cashierData.length === 0) && (
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
