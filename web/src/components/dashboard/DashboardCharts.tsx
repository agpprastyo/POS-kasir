/**
 * DashboardCharts.tsx
 *
 * Komponen ini sengaja dipisah dari dashboard index agar Recharts (~440 KB)
 * di-lazy load dan tidak masuk initial bundle.
 *
 * Di-import via: const DashboardCharts = lazy(() => import('@/components/dashboard/DashboardCharts'))
 */
import {
    ResponsiveContainer,
    BarChart,
    Bar,
    XAxis,
    YAxis,
    Tooltip as RechartsTooltip,
    CartesianGrid,
    PieChart,
    Pie,
    Cell,
    Legend,
} from 'recharts'
import { useTranslation } from 'react-i18next'

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d']

interface SalesData {
    date?: string
    total_sales?: number
}

interface PaymentData {
    payment_method_name?: string
    order_count?: number
}

interface DashboardChartsProps {
    salesData: SalesData[] | undefined
    paymentsData: PaymentData[] | undefined
    formatCurrency: (value: number) => string
    formatDate: (dateString: string) => string
    chartType: 'bar' | 'pie'
}

export default function DashboardCharts({
    salesData,
    paymentsData,
    formatCurrency,
    formatDate,
    chartType,
}: DashboardChartsProps) {
    const { t } = useTranslation()

    if (chartType === 'bar') {
        return (
            <ResponsiveContainer width="100%" height={350}>
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
                    <RechartsTooltip
                        formatter={(value: number) => formatCurrency(value)}
                        labelFormatter={(label) => formatDate(label)}
                    />
                    <Bar dataKey="total_sales" fill="#adfa1d" radius={[4, 4, 0, 0]} name={t('reports.sales.revenue')} />
                </BarChart>
            </ResponsiveContainer>
        )
    }

    // chartType === 'pie'
    return (
        <ResponsiveContainer width="100%" height={300}>
            <PieChart>
                <Pie
                    data={paymentsData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={100}
                    fill="#8884d8"
                    paddingAngle={5}
                    dataKey="order_count"
                    nameKey="payment_method_name"
                    label={({ percent }) => `${(percent * 100).toFixed(0)}%`}
                >
                    {(paymentsData || []).map((_entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                </Pie>
                <RechartsTooltip />
                <Legend verticalAlign="bottom" height={36} />
            </PieChart>
        </ResponsiveContainer>
    )
}
