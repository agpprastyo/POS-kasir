import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

interface CancellationReportProps {
    data: any[]
    isLoading: boolean
    t: any
}

export function CancellationReport({
    data, isLoading, t
}: CancellationReportProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>{t('reports.cancellations.title')}</CardTitle>
                <CardDescription>{t('reports.cancellations.description')}</CardDescription>
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
                                <TableHead>{t('reports.cancellations.reason')}</TableHead>
                                <TableHead className="text-right">{t('reports.cancellations.count')}</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {(data || []).map((report: any, idx: number) => (
                                <TableRow key={idx}>
                                    <TableCell className="font-medium">{report.reason}</TableCell>
                                    <TableCell className="text-right">{report.cancelled_orders ?? 0}</TableCell>
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
