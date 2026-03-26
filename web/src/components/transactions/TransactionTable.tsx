import { Loader2, Utensils, ShoppingBag, CheckCircle, XCircle, Banknote } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
} from "@/components/ui/select"
import { formatRupiah } from '@/lib/utils'
import { POSKasirInternalOrdersRepositoryOrderStatus } from '@/lib/api/generated'

interface TransactionTableProps {
    orders: any[]
    isLoading: boolean
    userMap: Map<string, string>
    mutatingOrderId: string | null
    handleStatusUpdate: (id: string, newStatus: POSKasirInternalOrdersRepositoryOrderStatus) => void
    handleOpenPayment: (order: any) => void
    handleOpenCancel: (order: any) => void
    handleOpenRefund: (order: any) => void
    handleFinish: (order: any) => void
    t: any
}

export function TransactionTable({
    orders, isLoading, userMap, mutatingOrderId, handleStatusUpdate,
    handleOpenPayment, handleOpenCancel, handleOpenRefund, handleFinish, t
}: TransactionTableProps) {
    if (isLoading) {
        return (
            <div className="h-48 flex items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-primary/50" />
            </div>
        )
    }

    return (
        <Table>
            <TableHeader>
                <TableRow className="bg-muted/50">
                    <TableHead className="w-[100px]">{t('transactions.table.queue')}</TableHead>
                    <TableHead>{t('transactions.table.type')}</TableHead>
                    <TableHead>{t('transactions.table.cashier')}</TableHead>
                    <TableHead className="w-[300px]">{t('transactions.table.items')}</TableHead>
                    <TableHead>{t('transactions.table.total')}</TableHead>
                    <TableHead>{t('transactions.table.status')}</TableHead>
                    <TableHead className="text-right">{t('transactions.table.action')}</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {orders.length === 0 ? (
                    <TableRow>
                        <TableCell colSpan={7} className="h-32 text-center text-muted-foreground">
                            {t('transactions.table.no_data')}
                        </TableCell>
                    </TableRow>
                ) : (
                    orders.map((order) => (
                        <TableRow key={order.id} className="hover:bg-muted/30">
                            <TableCell className="font-medium">
                                <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-primary/10 text-primary font-mono ">
                                    {order.queue_number}
                                </div>
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-2">
                                    {order.type === 'dine_in' ? (
                                        <Utensils className="h-4 w-4 text-primary" />
                                    ) : (
                                        <ShoppingBag className="h-4 w-4 text-primary" />
                                    )}
                                    <div className="flex flex-col">
                                        <span className="font-medium capitalize">{order.type?.replace('_', ' ')}</span>
                                        <span className="text-xs text-muted-foreground">{new Date(order.created_at || '').toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
                                    </div>
                                </div>
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-2">
                                    <div className="h-6 w-6 rounded-full bg-muted flex items-center justify-center text-xs">
                                        {userMap.get(order.user_id || '')?.slice(0, 1).toUpperCase()}
                                    </div>
                                    <span className="text-sm">{userMap.get(order.user_id || '') || '-'}</span>
                                </div>
                            </TableCell>
                            <TableCell>
                                <div className="flex flex-col gap-1 text-sm text-muted-foreground">
                                    {order.items?.slice(0, 2).map((item: any, idx: number) => (
                                        <span key={idx}>
                                            {item.quantity}x <span className="text-foreground">{item.product_name}</span>
                                        </span>
                                    ))}
                                    {(order.items?.length || 0) > 2 && (
                                        <span className="text-xs italic">+{(order.items?.length || 0) - 2} more items...</span>
                                    )}
                                </div>
                            </TableCell>
                            <TableCell className="font-bold font-mono text-base">
                                <div className="flex flex-col items-start gap-1">
                                    <span>{formatRupiah(order.net_total || 0)}</span>
                                    {order.is_paid ? (
                                        <Badge variant="default" className="text-[10px] h-5">
                                            {t('transactions.status_badge.paid')}
                                        </Badge>
                                    ) : (
                                        <Badge variant="secondary" className="text-[10px] h-5">
                                            {t('transactions.status_badge.unpaid')}
                                        </Badge>
                                    )}
                                </div>
                            </TableCell>
                            <TableCell>
                                <Select
                                    value={order.status}
                                    disabled={mutatingOrderId === order.id || order.status === POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusPaid || order.status === POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusCancelled}
                                    onValueChange={(val) => {
                                        if (order.id && val !== POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusPaid && val !== POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusCancelled) {
                                            handleStatusUpdate(order.id, val as POSKasirInternalOrdersRepositoryOrderStatus)
                                        } else {
                                            // Optional: Alerting that specific status needs separate action
                                        }
                                    }}
                                >
                                    <SelectTrigger className="w-[140px] h-8 border-none bg-transparent hover:bg-muted focus:ring-0 p-0">
                                        <Badge
                                            variant={order.status === POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusPaid ? 'default' : order.status === POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusCancelled ? 'destructive' : 'secondary'}
                                            className={`capitalize cursor-pointer w-full justify-center gap-2`}
                                        >
                                            {mutatingOrderId === order.id && <Loader2 className="h-3 w-3 animate-spin" />}
                                            {order.status ? t(`transactions.status.${order.status}`) : order.status}
                                        </Badge>
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value={POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusOpen}>{t('transactions.status.open')}</SelectItem>
                                        <SelectItem value={POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusInProgress}>{t('transactions.status.in_progress')}</SelectItem>
                                        <SelectItem value={POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusServed}>{t('transactions.status.served')}</SelectItem>
                                    </SelectContent>
                                </Select>
                            </TableCell>
                            <TableCell className="text-right">
                                <div className="flex items-center justify-end gap-2">
                                    {!order.is_paid && order.status !== POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusCancelled && (
                                        <Button size="sm" variant="default" className="h-8 gap-1" onClick={() => handleOpenPayment(order)}>
                                            <Banknote className="h-3.5 w-3.5" />
                                            {t('transactions.actions_button.pay')}
                                        </Button>
                                    )}
                                    {order.is_paid && order.status === POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusInProgress && (
                                        <Button size="sm" variant="outline" className="h-8 gap-1 border-primary text-primary hover:bg-primary/10" onClick={() => handleStatusUpdate(order.id!, POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusServed)} disabled={mutatingOrderId === order.id}>
                                            <Utensils className="h-3.5 w-3.5" />
                                            {t('transactions.status.served')}
                                        </Button>
                                    )}
                                    {order.is_paid && order.status === POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusServed && (
                                        <Button size="sm" variant="default" className="h-8 gap-1" onClick={() => handleFinish(order)} disabled={mutatingOrderId === order.id}>
                                            <CheckCircle className="h-3.5 w-3.5" />
                                            {t('transactions.actions_button.complete')}
                                        </Button>
                                    )}
                                    {(!order.is_paid && order.status !== POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusCancelled && order.status !== POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusPaid) && (
                                        <Button size="sm" variant="destructive" className="h-8 gap-1" onClick={() => handleOpenCancel(order)}>
                                            <XCircle className="h-3.5 w-3.5" />
                                            {t('transactions.actions_button.cancel')}
                                        </Button>
                                    )}
                                    {order.is_paid && order.status !== POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusCancelled && (
                                        <Button size="sm" variant="destructive" className="h-8 gap-1" onClick={() => handleOpenRefund(order)}>
                                            <XCircle className="h-3.5 w-3.5" />
                                            {t('transactions.actions_button.refund', 'Refund')}
                                        </Button>
                                    )}
                                </div>
                            </TableCell>
                        </TableRow>
                    ))
                )}
            </TableBody>
        </Table>
    )
}
