import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { ChevronLeft, ChevronRight, ChevronsUpDown } from 'lucide-react'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { formatDateTime } from '@/lib/utils'
import { InternalActivitylogActivityLogResponse } from '@/lib/api/generated'

interface ActivityLogsTableProps {
    t: any
    logs: InternalActivitylogActivityLogResponse[] | undefined
    page: number
    totalPages: number
    onPageChange: (newPage: number) => void
}

export function ActivityLogsTable({
    t,
    logs,
    page,
    totalPages,
    onPageChange
}: ActivityLogsTableProps) {
    return (
        <div className="space-y-4">
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>{t('activity_logs.columns.date')}</TableHead>
                            <TableHead>{t('activity_logs.columns.user')}</TableHead>
                            <TableHead>{t('activity_logs.columns.action')}</TableHead>
                            <TableHead>{t('activity_logs.columns.entity')}</TableHead>
                            <TableHead>{t('activity_logs.columns.details')}</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {logs?.map((log: InternalActivitylogActivityLogResponse) => (
                            <TableRow key={log.id}>
                                <TableCell>{formatDateTime(log.created_at!)}</TableCell>
                                <TableCell>
                                    <div className="font-medium">{log.user_name}</div>
                                    <div className="text-sm text-muted-foreground">{t('activity_logs.table.id_prefix')} {log.user_id?.substring(0, 8)}...</div>
                                </TableCell>
                                <TableCell>
                                    <span className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-sm font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-primary text-primary-foreground hover:bg-primary/80">
                                        {log.action_type}
                                    </span>
                                </TableCell>
                                <TableCell>
                                    <div className="font-medium">{log.entity_type}</div>
                                    <div className="text-sm text-muted-foreground">{log.entity_id}</div>
                                </TableCell>
                                <TableCell>
                                    <Collapsible>
                                        <div className="flex items-center gap-2">
                                            <CollapsibleTrigger asChild>
                                                <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
                                                    <ChevronsUpDown className="h-4 w-4" />
                                                    <span className="sr-only">{t('activity_logs.table.toggle')}</span>
                                                </Button>
                                            </CollapsibleTrigger>
                                            <span className="text-sm text-muted-foreground truncate max-w-[200px]">
                                                {JSON.stringify(log.details)}
                                            </span>
                                        </div>
                                        <CollapsibleContent>
                                            <pre className="mt-2 w-[300px] overflow-auto rounded-md bg-muted p-2 text-sm">
                                                {JSON.stringify(log.details, null, 2)}
                                            </pre>
                                        </CollapsibleContent>
                                    </Collapsible>
                                </TableCell>
                            </TableRow>
                        ))}
                        {(!logs || logs.length === 0) && (
                            <TableRow>
                                <TableCell colSpan={5} className="h-24 text-center">
                                    {t('common.no_results')}
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>

            <div className="flex items-center justify-end space-x-2 py-4">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onPageChange(page - 1)}
                    disabled={page <= 1}
                >
                    <ChevronLeft className="h-4 w-4" />
                    {t('common.previous')}
                </Button>
                <div className="text-sm font-medium">
                    {t('common.page_info', {
                        current: page,
                        total: totalPages || 1
                    })}
                </div>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onPageChange(page + 1)}
                    disabled={page >= (totalPages || 1)}
                >
                    {t('common.next')}
                    <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
        </div>
    )
}
