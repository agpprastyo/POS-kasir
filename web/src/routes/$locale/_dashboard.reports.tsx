import { createFileRoute, redirect } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { meQueryOptions } from '@/lib/api/query/auth'
import { queryClient } from '@/lib/queryClient'
import { POSKasirInternalRepositoryUserRole } from '@/lib/api/generated/models/poskasir-internal-repository-user-role'
import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useSalesReportQuery, useProductPerformanceQuery, usePaymentMethodPerformanceQuery, useCashierPerformanceQuery, useCancellationReportsQuery } from '@/lib/api/query/reports'
import { Skeleton } from '@/components/ui/skeleton'
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid, PieChart, Pie, Cell, Legend } from 'recharts'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'

export const Route = createFileRoute('/$locale/_dashboard/reports')({
    beforeLoad: async () => {
        const user = await queryClient.ensureQueryData(meQueryOptions())

        const allowedRoles: POSKasirInternalRepositoryUserRole[] = [
            POSKasirInternalRepositoryUserRole.UserRoleAdmin,
            POSKasirInternalRepositoryUserRole.UserRoleManager
        ]

        if (!user.role || !allowedRoles.includes(user.role)) {
            throw redirect({
                to: '/',
                search: {
                    error: 'Unauthorized',
                },
            })
        }
    },
    component: ReportsPage,
})

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d'];

function ReportsPage() {
    const { t } = useTranslation()

    const [dateRange, setDateRange] = useState({
        start: new Date(new Date().setDate(new Date().getDate() - 30)).toISOString().split('T')[0],
        end: new Date().toISOString().split('T')[0]
    })

    const { data: salesData, isLoading: isLoadingSales } = useSalesReportQuery(dateRange.start, dateRange.end)
    const { data: productsData, isLoading: isLoadingProducts } = useProductPerformanceQuery(dateRange.start, dateRange.end)
    const { data: paymentsData, isLoading: isLoadingPayments } = usePaymentMethodPerformanceQuery(dateRange.start, dateRange.end)
    const { data: cashierData, isLoading: isLoadingCashier } = useCashierPerformanceQuery(dateRange.start, dateRange.end)
    const { data: cancellationData, isLoading: isLoadingCancellation } = useCancellationReportsQuery(dateRange.start, dateRange.end)

    const formatCurrency = (value: number) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0
        }).format(value)
    }

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('id-ID', {
            day: 'numeric',
            month: 'short',
            year: 'numeric'
        })
    }

    const handleDateChange = (type: 'start' | 'end', value: string) => {
        setDateRange(prev => ({
            ...prev,
            [type]: value
        }))
    }

    return (
        <div className="flex flex-col gap-4 p-4 md:p-8">
            <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">{t('reports.title')}</h1>
                    <p className="text-muted-foreground">{t('reports.subtitle')}</p>
                </div>
                <div className="flex items-center gap-2">
                    <div className="grid gap-1.5">
                        <label htmlFor="start-date" className="text-xs font-medium text-muted-foreground">Start Date</label>
                        <div className="relative">
                            <input
                                type="date"
                                id="start-date"
                                value={dateRange.start}
                                onChange={(e) => handleDateChange('start', e.target.value)}
                                className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                            />
                        </div>
                    </div>
                    <div className="grid gap-1.5">
                        <label htmlFor="end-date" className="text-xs font-medium text-muted-foreground">End Date</label>
                        <div className="relative">
                            <input
                                type="date"
                                id="end-date"
                                value={dateRange.end}
                                onChange={(e) => handleDateChange('end', e.target.value)}
                                className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                            />
                        </div>
                    </div>
                </div>
            </div>

            <Tabs defaultValue="sales" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="sales">{t('reports.tabs.sales')}</TabsTrigger>
                    <TabsTrigger value="products">{t('reports.tabs.products')}</TabsTrigger>
                    <TabsTrigger value="performance">{t('reports.tabs.performance')}</TabsTrigger>
                    <TabsTrigger value="cancellations">{t('reports.tabs.cancellations')}</TabsTrigger>
                </TabsList>

                {/* Sales Report Tab */}
                <TabsContent value="sales" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>{t('reports.sales.title')}</CardTitle>
                            <CardDescription>{t('reports.sales.description')}</CardDescription>
                        </CardHeader>
                        <CardContent className="pl-2">
                            {isLoadingSales ? (
                                <Skeleton className="h-[400px] w-full" />
                            ) : (
                                <ResponsiveContainer width="100%" height={400}>
                                    <BarChart data={salesData}>
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
                                            formatter={(value: number) => formatCurrency(value)}
                                            labelFormatter={(label) => formatDate(label)}
                                        />
                                        <Legend />
                                        <Bar dataKey="total_sales" fill="#adfa1d" radius={[4, 4, 0, 0]} name={t('reports.sales.revenue')} />
                                        <Bar dataKey="order_count" fill="#8884d8" radius={[4, 4, 0, 0]} name={t('reports.sales.orders')} />
                                    </BarChart>
                                </ResponsiveContainer>
                            )}
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* Products Report Tab */}
                <TabsContent value="products" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>{t('reports.products.title')}</CardTitle>
                            <CardDescription>{t('reports.products.description')}</CardDescription>
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
                                            <TableHead className="text-right">{t('reports.products.quantity')}</TableHead>
                                            <TableHead className="text-right">{t('reports.products.revenue')}</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {(productsData || []).map((product) => (
                                            <TableRow key={product.product_id}>
                                                <TableCell className="font-medium">{product.product_name}</TableCell>
                                                <TableCell className="text-right">{product.total_quantity ?? 0}</TableCell>
                                                <TableCell className="text-right">{formatCurrency(product.total_revenue ?? 0)}</TableCell>
                                            </TableRow>
                                        ))}
                                        {(!productsData || productsData.length === 0) && (
                                            <TableRow>
                                                <TableCell colSpan={3} className="text-center text-muted-foreground">{t('common.no_data')}</TableCell>
                                            </TableRow>
                                        )}
                                    </TableBody>
                                </Table>
                            )}
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* Performance Tab (Payment & Cashier) */}
                <TabsContent value="performance" className="space-y-4">
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                        <Card className="col-span-4">
                            <CardHeader>
                                <CardTitle>{t('reports.performance.payment_methods')}</CardTitle>
                            </CardHeader>
                            <CardContent className="pl-2">
                                {isLoadingPayments ? (
                                    <Skeleton className="h-[300px] w-full" />
                                ) : (
                                    <div className="h-[300px] w-full">
                                        <ResponsiveContainer width="100%" height="100%">
                                            <PieChart>
                                                <Pie
                                                    data={paymentsData}
                                                    cx="50%"
                                                    cy="50%"
                                                    labelLine={false}
                                                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                                                    outerRadius={100}
                                                    fill="#8884d8"
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
                        <Card className="col-span-3">
                            <CardHeader>
                                <CardTitle>{t('reports.performance.cashier')}</CardTitle>
                            </CardHeader>
                            <CardContent>
                                {isLoadingCashier ? (
                                    <div className="space-y-2">
                                        <Skeleton className="h-10 w-full" />
                                        <Skeleton className="h-10 w-full" />
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
                                            {(cashierData || []).map((cashier) => (
                                                <TableRow key={cashier.user_id}>
                                                    <TableCell className="font-medium">
                                                        <div className="flex items-center gap-2">
                                                            <Avatar className="h-6 w-6">
                                                                <AvatarFallback>{cashier.username?.substring(0, 2).toUpperCase()}</AvatarFallback>
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
                </TabsContent>

                {/* Cancellation Reports Tab */}
                <TabsContent value="cancellations" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>{t('reports.cancellations.title')}</CardTitle>
                            <CardDescription>{t('reports.cancellations.description')}</CardDescription>
                        </CardHeader>
                        <CardContent>
                            {isLoadingCancellation ? (
                                <div className="space-y-2">
                                    <Skeleton className="h-10 w-full" />
                                    <Skeleton className="h-10 w-full" />
                                </div>
                            ) : (
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead>{t('reports.cancellations.reason')}</TableHead>
                                            <TableHead className="text-right">{t('reports.cancellations.count')}</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {(cancellationData || []).map((report, idx) => (
                                            <TableRow key={idx}>
                                                <TableCell className="font-medium">{report.reason}</TableCell>
                                                <TableCell className="text-right">{report.cancelled_orders ?? 0}</TableCell>
                                            </TableRow>
                                        ))}
                                        {(!cancellationData || cancellationData.length === 0) && (
                                            <TableRow>
                                                <TableCell colSpan={2} className="text-center text-muted-foreground">{t('common.no_data')}</TableCell>
                                            </TableRow>
                                        )}
                                    </TableBody>
                                </Table>
                            )}
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    )
}

