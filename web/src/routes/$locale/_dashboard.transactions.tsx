import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { ordersListQueryOptions, useCompleteManualPaymentMutation, useUpdateOrderStatusMutation, useCancelOrderMutation } from '@/lib/api/query/orders'
import { usersListQueryOptions } from '@/lib/api/query/user'
import { useCancellationReasonsListQuery } from '@/lib/api/query/cancel-reason'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { formatRupiah } from '@/lib/utils'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useState, useMemo } from 'react'
import { Loader2, Utensils, ShoppingBag, CreditCard, User as UserIcon, Banknote, CheckCircle, XCircle } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
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
    SelectValue,
} from "@/components/ui/select"
import { usePaymentMethodsListQuery } from '@/lib/api/query/payment-methods'
import { toast } from 'sonner'
import { POSKasirInternalDtoOrderListResponse, POSKasirInternalRepositoryOrderStatus } from '@/lib/api/generated'
import { useTranslation } from 'react-i18next'

const transactionsSearchSchema = z.object({
    page: z.number().catch(1),
    limit: z.number().min(1).catch(10),
    status: z.enum(['active', 'paid', 'cancelled', 'all']).catch('active'),
    user_id: z.string().optional(),
})

export const Route = createFileRoute('/$locale/_dashboard/transactions')({
    validateSearch: (search) => transactionsSearchSchema.parse(search),
    loaderDeps: ({ search }) => ({
        page: search.page,
        limit: search.limit,
        status: search.status,
        user_id: search.user_id,
    }),
    component: TransactionsPage,
})

type OrderWithItems = POSKasirInternalDtoOrderListResponse & {
    items?: Array<{
        product_name: string;
        quantity: number;
        product_option_name?: string;
        price: number;
    }>;
    payment_amount?: number;
    queue_number?: string;
    is_paid?: boolean;
}

function TransactionsPage() {
    const { t } = useTranslation()
    const navigate = useNavigate({ from: Route.fullPath })
    const searchParams = Route.useSearch()

    // Derived state from search params
    const selectedTab = searchParams.status || 'active'
    const selectedUserId = searchParams.user_id || 'all'
    const page = searchParams.page
    const limit = searchParams.limit

    // Fetch Users for Filter
    const { data: usersData } = useQuery(usersListQueryOptions({ page: 1, limit: 100 }))
    const users = usersData?.users || []

    // Create Map for User ID -> Username
    const userMap = useMemo(() => {
        const map = new Map<string, string>()
        users.forEach(u => {
            if (u.id) map.set(u.id, u.username || 'Unknown')
        })
        return map
    }, [users])

    // Map tab to API statuses
    const statusesFilter = useMemo(() => {
        if (selectedTab === 'active') return ['open', 'in_progress', 'served']
        if (selectedTab === 'paid') return ['paid']
        if (selectedTab === 'cancelled') return ['cancelled']
        return undefined // 'all'
    }, [selectedTab])

    // Fetch Orders
    const { data, isLoading } = useQuery(ordersListQueryOptions({
        page,
        limit,
        statuses: statusesFilter,
        userId: selectedUserId === 'all' ? undefined : selectedUserId
    }))

    const orders = (data?.orders || []) as OrderWithItems[]
    const pagination = data?.pagination || { total_pages: 1, current_page: 1 }

    // Payment Dialog State
    const [isPaymentOpen, setIsPaymentOpen] = useState(false)
    const [selectedOrder, setSelectedOrder] = useState<OrderWithItems | null>(null)
    const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<number | undefined>(undefined)
    const [cashReceived, setCashReceived] = useState<string>('')

    // Cancellation Dialog State
    const [isCancelDialogOpen, setIsCancelDialogOpen] = useState(false)
    const [orderToCancel, setOrderToCancel] = useState<OrderWithItems | null>(null)
    const [cancelReasonId, setCancelReasonId] = useState<string>('')
    const [cancelNotes, setCancelNotes] = useState('')

    const { data: paymentMethods } = usePaymentMethodsListQuery()
    const { data: cancellationReasons } = useCancellationReasonsListQuery()
    const completeManualPaymentMutation = useCompleteManualPaymentMutation()
    const updateOrderStatusMutation = useUpdateOrderStatusMutation()
    const cancelOrderMutation = useCancelOrderMutation()

    const handleStatusUpdate = async (id: string, newStatus: POSKasirInternalRepositoryOrderStatus) => {
        try {
            await updateOrderStatusMutation.mutateAsync({
                id,
                body: { status: newStatus }
            }) // Toast handled by mutation onSuccess
        } catch (error) {
            console.error(error)
        }
    }

    const handleOpenPayment = (order: OrderWithItems) => {
        setSelectedOrder(order)
        setIsPaymentOpen(true)
        setCashReceived('')
        setSelectedPaymentMethod(undefined)
    }

    const handleOpenCancel = (order: OrderWithItems) => {
        setOrderToCancel(order)
        setCancelReasonId('')
        setCancelNotes('')
        setIsCancelDialogOpen(true)
    }

    const handleConfirmCancel = async () => {
        if (!orderToCancel?.id || !cancelReasonId) return;

        try {
            await cancelOrderMutation.mutateAsync({
                id: orderToCancel.id,
                body: {
                    cancellation_reason_id: Number(cancelReasonId),
                    cancellation_notes: cancelNotes
                }
            })
            // Success toast handled by mutation
            setIsCancelDialogOpen(false)
            setOrderToCancel(null)
        } catch (error) {
            console.error(error)
        }
    }

    const handleFinish = async (order: OrderWithItems) => {
        if (!order.id) return;
        // Optional: Confirm dialog?
        try {
            await handleStatusUpdate(order.id, POSKasirInternalRepositoryOrderStatus.OrderStatusPaid)
            toast.success(t('transactions.messages.order_finished'))
        } catch (e) {
            console.error(e)
        }
    }


    const handleTabChange = (value: string) => {
        navigate({
            search: (prev) => ({ ...prev, status: value as any, page: 1 })
        })
    }

    const handleUserChange = (value: string) => {
        navigate({
            search: (prev) => ({ ...prev, user_id: value === 'all' ? undefined : value, page: 1 })
        })
    }

    const handlePageChange = (newPage: number) => {
        navigate({
            search: (prev) => ({ ...prev, page: newPage })
        })
    }

    const handlePayment = async () => {
        if (!selectedOrder || !selectedPaymentMethod) {
            toast.error(t('transactions.messages.select_payment_method'))
            return
        }
        if (!selectedOrder.id) return;

        const method = paymentMethods?.find((m: any) => m.id === selectedPaymentMethod)
        const isCash = method?.name?.toLowerCase().includes('cash')
        const totalAmount = selectedOrder.net_total || 0

        let payload: any = {
            payment_method_id: selectedPaymentMethod
        }

        let finalCashReceived = 0

        if (isCash) {
            const inputCash = Number(cashReceived)
            if (inputCash < totalAmount) {
                toast.error(t('transactions.payment_dialog.money_insufficient'))
                return
            }
            finalCashReceived = inputCash
            payload.cash_received = finalCashReceived
        }

        try {
            await completeManualPaymentMutation.mutateAsync({
                id: selectedOrder.id,
                body: payload
            })

            setIsPaymentOpen(false)
            setSelectedOrder(null)

            if (isCash) {
                const change = finalCashReceived - totalAmount
                toast.success(`${t('transactions.payment_dialog.success')} ${formatRupiah(change)}`)
            }
        } catch (error) {
            console.error(error)
        }
    }


    return (
        <div className="flex flex-col gap-6 p-4 md:p-8 max-w-7xl mx-auto w-full">
            <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">{t('transactions.title')}</h1>
                    <p className="text-muted-foreground">{t('transactions.subtitle')}</p>
                </div>

                {/* User Filter */}
                <div className="flex items-center gap-2">
                    <UserIcon className="h-4 w-4 text-muted-foreground" />
                    <Select value={selectedUserId} onValueChange={handleUserChange}>
                        <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder={t('transactions.filter_user')} />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">{t('transactions.filter_user')}</SelectItem>
                            {users.map(user => (
                                <SelectItem key={user.id} value={user.id || ''}>{user.username}</SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>
            </div>

            <div className="space-y-4">
                <Tabs value={selectedTab} onValueChange={handleTabChange} className="w-full">
                    <TabsList className="flex w-full overflow-x-auto justify-start h-auto p-1 gap-2 bg-muted/20">
                        <TabsTrigger value="active" className="flex-shrink-0">{t('transactions.tabs.active')}</TabsTrigger>
                        <TabsTrigger value="paid" className="flex-shrink-0">{t('transactions.tabs.history')}</TabsTrigger>
                        <TabsTrigger value="cancelled" className="flex-shrink-0">{t('transactions.tabs.cancelled')}</TabsTrigger>
                        <TabsTrigger value="all" className="flex-shrink-0">{t('transactions.status.all')}</TabsTrigger>
                    </TabsList>
                </Tabs>

                <div className="rounded-md border bg-card ">
                    {isLoading ? (
                        <div className="h-48 flex items-center justify-center">
                            <Loader2 className="h-8 w-8 animate-spin text-primary/50" />
                        </div>
                    ) : (
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
                                                        <Utensils className="h-4 w-4 text-orange-500" />
                                                    ) : (
                                                        <ShoppingBag className="h-4 w-4 text-blue-500" />
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
                                                    {order.items?.slice(0, 2).map((item, idx) => (
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
                                                        <Badge variant="outline" className="text-[10px] h-5 bg-green-50 text-green-700 border-green-200">
                                                            {t('transactions.status_badge.paid')}
                                                        </Badge>
                                                    ) : (
                                                        <Badge variant="outline" className="text-[10px] h-5 bg-yellow-50 text-yellow-700 border-yellow-200">
                                                            {t('transactions.status_badge.unpaid')}
                                                        </Badge>
                                                    )}
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                <Select
                                                    value={order.status}
                                                    disabled={order.status === POSKasirInternalRepositoryOrderStatus.OrderStatusPaid || order.status === POSKasirInternalRepositoryOrderStatus.OrderStatusCancelled}
                                                    onValueChange={(val) => {
                                                        if (order.id && val !== POSKasirInternalRepositoryOrderStatus.OrderStatusPaid && val !== POSKasirInternalRepositoryOrderStatus.OrderStatusCancelled) {
                                                            handleStatusUpdate(order.id, val as POSKasirInternalRepositoryOrderStatus)
                                                        } else {

                                                            toast.info(t('transactions.messages.use_action_button'))
                                                        }
                                                    }}
                                                >
                                                    <SelectTrigger className="w-[140px] h-8 border-none bg-transparent hover:bg-muted focus:ring-0 p-0">
                                                        <Badge
                                                            variant={order.status === POSKasirInternalRepositoryOrderStatus.OrderStatusPaid ? 'default' : order.status === POSKasirInternalRepositoryOrderStatus.OrderStatusCancelled ? 'destructive' : 'secondary'}
                                                            className={`capitalize cursor-pointer w-full justify-center ${order.status === POSKasirInternalRepositoryOrderStatus.OrderStatusPaid ? 'bg-green-600 hover:bg-green-700' : ''}`}
                                                        >
                                                            {order.status ? t(`transactions.status.${order.status}`) : order.status}
                                                        </Badge>
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        <SelectItem value={POSKasirInternalRepositoryOrderStatus.OrderStatusOpen}>{t('transactions.status.open')}</SelectItem>
                                                        <SelectItem value={POSKasirInternalRepositoryOrderStatus.OrderStatusInProgress}>{t('transactions.status.in_progress')}</SelectItem>
                                                        <SelectItem value={POSKasirInternalRepositoryOrderStatus.OrderStatusServed}>{t('transactions.status.served')}</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                            </TableCell>
                                            <TableCell className="text-right">
                                                <div className="flex items-center justify-end gap-2">
                                                    {!order.is_paid && order.status !== POSKasirInternalRepositoryOrderStatus.OrderStatusCancelled && (
                                                        <Button size="sm" variant="default" className="h-8 gap-1 bg-blue-600 hover:bg-blue-700" onClick={() => handleOpenPayment(order)}>
                                                            <Banknote className="h-3.5 w-3.5" />
                                                            {t('transactions.actions_button.pay')}
                                                        </Button>
                                                    )}
                                                    {order.is_paid && order.status !== POSKasirInternalRepositoryOrderStatus.OrderStatusPaid && order.status !== POSKasirInternalRepositoryOrderStatus.OrderStatusCancelled && (
                                                        <Button size="sm" variant="outline" className="h-8 gap-1 border-green-600 text-green-600 hover:bg-green-50" onClick={() => handleFinish(order)}>
                                                            <CheckCircle className="h-3.5 w-3.5" />
                                                            {t('transactions.actions_button.complete')}
                                                        </Button>
                                                    )}
                                                    {(!order.is_paid && order.status !== POSKasirInternalRepositoryOrderStatus.OrderStatusCancelled && order.status !== POSKasirInternalRepositoryOrderStatus.OrderStatusPaid) && (
                                                        <Button size="sm" variant="destructive" className="h-8 gap-1" onClick={() => handleOpenCancel(order)}>
                                                            <XCircle className="h-3.5 w-3.5" />
                                                            {t('transactions.actions_button.cancel')}
                                                        </Button>
                                                    )}
                                                </div>
                                            </TableCell>
                                        </TableRow>
                                    ))
                                )}
                            </TableBody>
                        </Table>
                    )}
                </div>

                {/* Pagination */}
                {pagination.total_pages > 1 && (
                    <div className="flex items-center justify-end gap-2 pt-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handlePageChange(Math.max(1, page - 1))}
                            disabled={page === 1}
                        >
                            {t('pagination.previous')}
                        </Button>
                        <div className="text-sm font-medium">{t('pagination.page_info', { current: page, total: pagination.total_pages })}</div>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handlePageChange(Math.min(pagination.total_pages, page + 1))}
                            disabled={page === pagination.total_pages}
                        >
                            {t('pagination.next')}
                        </Button>
                    </div>
                )}
            </div>

            {/* Payment Dialog */}
            <Dialog open={isPaymentOpen} onOpenChange={setIsPaymentOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>{t('transactions.payment_dialog.title')}</DialogTitle>
                        <DialogDescription>
                            {t('transactions.payment_dialog.description')} <b>#{selectedOrder?.queue_number}</b>
                        </DialogDescription>
                    </DialogHeader>

                    <div className="py-4">
                        <div className="mb-6 flex flex-col items-center justify-center rounded-lg border bg-muted/30 p-4">
                            <span className="text-sm text-muted-foreground">{t('transactions.payment_dialog.total_amount')}</span>
                            <span className="text-3xl font-bold text-foreground">{formatRupiah(selectedOrder?.net_total || 0)}</span>
                        </div>

                        <label className="text-sm font-medium mb-2 block">{t('transactions.payment_dialog.select_method')}</label>
                        <Tabs value={selectedPaymentMethod ? String(selectedPaymentMethod) : ""} onValueChange={(v) => setSelectedPaymentMethod(Number(v))} className="w-full">
                            <TabsList className="grid grid-cols-3 h-auto gap-2 bg-transparent p-0 mb-4">
                                {paymentMethods?.map((method: any) => (
                                    <TabsTrigger
                                        key={method.id}
                                        value={String(method.id)}
                                        className="h-20 flex flex-col items-center justify-center gap-2 border data-[state=active]:border-primary data-[state=active]:bg-primary/5 data-[state=active]:text-primary"
                                    >
                                        <CreditCard className="h-5 w-5 opacity-70" />
                                        <span className="text-xs text-wrap text-center">{method.name}</span>
                                    </TabsTrigger>
                                ))}
                            </TabsList>

                            {paymentMethods?.map((method: any) => (
                                <TabsContent key={method.id} value={String(method.id)} className="mt-0">
                                    {method.name?.toLowerCase().includes('cash') && (
                                        <div className="space-y-4 rounded-lg border border-dashed p-4 bg-muted/30 animate-in fade-in-0 zoom-in-95">
                                            <div className="space-y-2">
                                                <label className="text-sm font-medium">{t('transactions.payment_dialog.cash_received')}</label>
                                                <div className="relative">
                                                    <span className="absolute left-3 top-1/2 -translate-y-1/2 font-bold text-muted-foreground">Rp</span>
                                                    <Input
                                                        autoFocus
                                                        type="text"
                                                        inputMode="numeric"
                                                        value={cashReceived ? Number(cashReceived).toLocaleString('id-ID') : ''}
                                                        onChange={(e) => {
                                                            const val = e.target.value.replace(/\D/g, '')
                                                            setCashReceived(val)
                                                        }}
                                                        className="pl-10 text-lg font-bold h-12"
                                                        placeholder="0"
                                                    />
                                                </div>
                                            </div>
                                            <div className="flex justify-between items-center bg-background p-3 rounded border ">
                                                <span className="text-sm font-medium text-muted-foreground">{t('transactions.payment_dialog.change_due')}</span>
                                                <span className={`text-xl font-bold ${Number(cashReceived) >= (selectedOrder?.net_total || 0) ? 'text-green-600' : 'text-red-500'}`}>
                                                    {formatRupiah(Math.max(0, Number(cashReceived) - (selectedOrder?.net_total || 0)))}
                                                </span>
                                            </div>
                                        </div>
                                    )}
                                </TabsContent>
                            ))}
                        </Tabs>
                    </div>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsPaymentOpen(false)}>{t('common.cancel')}</Button>
                        <Button onClick={handlePayment} disabled={!selectedPaymentMethod || completeManualPaymentMutation.isPending}>
                            {completeManualPaymentMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t('transactions.payment_dialog.confirm')}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Cancellation Dialog */}
            <Dialog open={isCancelDialogOpen} onOpenChange={setIsCancelDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>{t('transactions.cancellation_dialog.title')}</DialogTitle>
                        <DialogDescription>
                            {t('transactions.cancellation_dialog.description')} <b>#{orderToCancel?.queue_number}</b>
                        </DialogDescription>
                    </DialogHeader>

                    <div className="space-y-4 py-4">
                        <div className="space-y-2">
                            <label className="text-sm font-medium">{t('transactions.cancellation_dialog.reason_label')}</label>
                            <Select value={cancelReasonId} onValueChange={setCancelReasonId}>
                                <SelectTrigger>
                                    <SelectValue placeholder={t('transactions.cancellation_dialog.reason_placeholder')} />
                                </SelectTrigger>
                                <SelectContent>
                                    {cancellationReasons?.map((reason) => (
                                        <SelectItem key={reason.id} value={String(reason.id)}>
                                            {reason.reason}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm font-medium">{t('transactions.cancellation_dialog.notes_label')}</label>
                            <Textarea
                                value={cancelNotes}
                                onChange={(e) => setCancelNotes(e.target.value)}
                                placeholder={t('transactions.cancellation_dialog.notes_placeholder')}
                            />
                        </div>
                    </div>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsCancelDialogOpen(false)}>{t('transactions.cancellation_dialog.cancel')}</Button>
                        <Button variant="destructive" onClick={handleConfirmCancel} disabled={!cancelReasonId || cancelOrderMutation.isPending}>
                            {cancelOrderMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t('transactions.cancellation_dialog.confirm')}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

        </div >
    )
}

