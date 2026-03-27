import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

interface DashboardTopProductsProps {
    t: any
    isLoading: boolean
    topProducts: any[]
    formatCurrency: (value: number) => string
}

export function DashboardTopProducts({ t, isLoading, topProducts, formatCurrency }: DashboardTopProductsProps) {
    return (
        <Card className="col-span-1 lg:col-span-3 flex flex-col border-0 shadow-sm">
            <CardHeader>
                <CardTitle>{t('dashboard.widgets.top_products')}</CardTitle>
                <CardDescription>{t('reports.products.description')}</CardDescription>
            </CardHeader>
            <CardContent className="flex-1">
                {isLoading ? (
                    <div className="space-y-4">
                        <Skeleton className="h-12 w-full" />
                        <Skeleton className="h-12 w-full" />
                        <Skeleton className="h-12 w-full" />
                    </div>
                ) : (
                    <div className="space-y-8">
                        {topProducts.map((product: any, index: number) => (
                            <div key={index} className="flex items-center">
                                <div className="flex h-9 w-9 items-center justify-center rounded-full border border-muted bg-primary/10 text-sm font-medium text-primary">
                                    {index + 1}
                                </div>
                                <div className="ml-4 space-y-1">
                                    <p className="text-sm font-medium leading-none">{product.product_name}</p>
                                    <p className="text-sm text-muted-foreground">
                                        {product.total_quantity} sold
                                    </p>
                                </div>
                                <div className="ml-auto font-medium">
                                    {formatCurrency(product.total_revenue ?? 0)}
                                </div>
                            </div>
                        ))}
                        {topProducts.length === 0 && (
                            <div className="text-center text-sm text-muted-foreground py-8">
                                {t('common.no_data')}
                            </div>
                        )}
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
