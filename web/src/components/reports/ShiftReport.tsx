import { Download } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

interface ShiftReportProps {
    data: any[]
    isLoading: boolean
    onExport: () => void
    formatCurrency: (value: number) => string
    t: any
}

export function ShiftReport({
    data, isLoading, onExport, formatCurrency, t
}: ShiftReportProps) {
    return (
        <Card className="border-0 shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between">
                <div>
                    <CardTitle>{t('reports.shifts.title') || 'Work Shifts'}</CardTitle>
                    <CardDescription>{t('reports.shifts.description') || 'Monitor cash reconcile accurately.'}</CardDescription>
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
                                <TableHead>{t('reports.performance.staff')}</TableHead>
                                <TableHead>{t('reports.shifts.status') || 'Status'}</TableHead>
                                <TableHead className="text-right">{t('reports.shifts.expected') || 'Expected'}</TableHead>
                                <TableHead className="text-right">{t('reports.shifts.actual') || 'Actual'}</TableHead>
                                <TableHead className="text-right">{t('reports.shifts.diff') || 'Difference'}</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {(data || []).map((shift: any, idx: number) => (
                                <TableRow key={idx} className="hover:bg-muted/50">
                                    <TableCell className="font-medium">{shift.cashier_name}</TableCell>
                                    <TableCell>{shift.status}</TableCell>
                                    <TableCell className="text-right">{formatCurrency(shift.expected_cash_end ?? 0)}</TableCell>
                                    <TableCell className="text-right">{formatCurrency(shift.actual_cash_end ?? 0)}</TableCell>
                                    <TableCell className={`text-right font-bold ${(shift.cash_difference ?? 0) < 0 ? 'text-destructive' : (shift.cash_difference ?? 0) > 0 ? 'text-emerald-500' : ''}`}>
                                        {formatCurrency(shift.cash_difference ?? 0)}
                                    </TableCell>
                                </TableRow>
                            ))}
                            {(!data || data.length === 0) && (
                                <TableRow>
                                    <TableCell colSpan={5} className="text-center text-muted-foreground">{t('common.no_data')}</TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                )}
            </CardContent>
        </Card>
    )
}
