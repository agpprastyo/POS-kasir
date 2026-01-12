import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { ordersListQueryOptions, useCompleteManualPaymentMutation } from '@/lib/api/query/orders'
import { usersListQueryOptions } from '@/lib/api/query/user'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { formatRupiah } from '@/lib/utils'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useState, useMemo } from 'react'
import { Loader2, Utensils, ShoppingBag, CreditCard, User as UserIcon } from 'lucide-react'
import { Input } from '@/components/ui/input'
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
import { POSKasirInternalDtoOrderListResponse } from '@/lib/api/generated'

export const Route = createFileRoute('/_dashboard/transactions')({
    component: TransactionsPage,
})

// Extend the generated type because 'items' might be missing in the spec but present in response
type OrderWithItems = POSKasirInternalDtoOrderListResponse & {
    items?: Array<{
        product_name: string;
        quantity: number;
        product_option_name?: string;
        price: number;
    }>;
    payment_amount?: number;
    queue_number?: string;
}

function TransactionsPage() {
    const [selectedStatus, setSelectedStatus] = useState<string>('pending') // pending, completed, cancelled, all
    const [selectedUserId, setSelectedUserId] = useState<string>('all')
    const [page, setPage] = useState(1)
    const [limit] = useState(10)

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

    // Fetch Orders
    const { data, isLoading } = useQuery(ordersListQueryOptions({
        page,
        limit,
        status: selectedStatus === 'all' ? undefined : selectedStatus,
        userId: selectedUserId === 'all' ? undefined : selectedUserId
    }))

    const orders = (data?.orders || []) as OrderWithItems[]
    const pagination = data?.pagination || { total_pages: 1, current_page: 1 }

    // Payment Dialog State
    const [isPaymentOpen, setIsPaymentOpen] = useState(false)
    const [selectedOrder, setSelectedOrder] = useState<OrderWithItems | null>(null)
    const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<number | undefined>(undefined)
    const [cashReceived, setCashReceived] = useState<string>('')

    const { data: paymentMethods } = usePaymentMethodsListQuery()
    const completeManualPaymentMutation = useCompleteManualPaymentMutation()

    const handleOpenPayment = (order: OrderWithItems) => {
        setSelectedOrder(order)
        setIsPaymentOpen(true)
        setCashReceived('')
        setSelectedPaymentMethod(undefined)
    }

    const handlePayment = async () => {
        if (!selectedOrder || !selectedPaymentMethod) {
            toast.error("Please select a payment method")
            return
        }
        if (!selectedOrder.id) return;

        const method = paymentMethods?.find((m: any) => m.id === selectedPaymentMethod)
        const isCash = method?.name?.toLowerCase().includes('cash')
        const totalAmount = selectedOrder.payment_amount || 0

        let payload: any = {
            payment_method_id: selectedPaymentMethod
        }

        let finalCashReceived = 0

        if (isCash) {
            const inputCash = Number(cashReceived)
            if (inputCash < totalAmount) {
                toast.error("Uang tunai kurang!")
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
                toast.success(`Pembayaran Berhasil! Kembalian: ${formatRupiah(change)}`)
            }
        } catch (error) {
            console.error(error)
        }
    }


    return (
        <div className="flex flex-col gap-6 p-4 md:p-8 max-w-7xl mx-auto w-full">
            <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Transactions</h1>
                    <p className="text-muted-foreground">Monitor orders and manage payments.</p>
                </div>

                {/* User Filter */}
                <div className="flex items-center gap-2">
                    <UserIcon className="h-4 w-4 text-muted-foreground" />
                    <Select value={selectedUserId} onValueChange={setSelectedUserId}>
                        <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="All Cashiers" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">All Cashiers</SelectItem>
                            {users.map(user => (
                                <SelectItem key={user.id} value={user.id || ''}>{user.username}</SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>
            </div>

            <div className="space-y-4">
                <Tabs value={selectedStatus} onValueChange={setSelectedStatus} className="w-full">
                    <TabsList className="grid w-full max-w-md grid-cols-4">
                        <TabsTrigger value="pending">Pending</TabsTrigger>
                        <TabsTrigger value="completed">Done</TabsTrigger>
                        <TabsTrigger value="cancelled">Cancel</TabsTrigger>
                        <TabsTrigger value="all">All</TabsTrigger>
                    </TabsList>
                </Tabs>

                <div className="rounded-md border bg-card shadow-sm">
                    {isLoading ? (
                        <div className="h-48 flex items-center justify-center">
                            <Loader2 className="h-8 w-8 animate-spin text-primary/50" />
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow className="bg-muted/50">
                                    <TableHead className="w-[100px]">Queue</TableHead>
                                    <TableHead>Type</TableHead>
                                    <TableHead>Cashier</TableHead>
                                    <TableHead className="w-[300px]">Items</TableHead>
                                    <TableHead>Total</TableHead>
                                    <TableHead>Status</TableHead>
                                    <TableHead className="text-right">Action</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {orders.length === 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={7} className="h-32 text-center text-muted-foreground">
                                            No transactions found for this filter.
                                        </TableCell>
                                    </TableRow>
                                ) : (
                                    orders.map((order) => (
                                        <TableRow key={order.id} className="hover:bg-muted/30">
                                            <TableCell className="font-medium">
                                                <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-primary/10 text-primary font-mono text-lg font-bold">
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
                                                {formatRupiah(order.payment_amount || 0)}
                                            </TableCell>
                                            <TableCell>
                                                <Badge
                                                    variant={order.status === 'completed' ? 'default' : order.status === 'cancelled' ? 'destructive' : 'secondary'}
                                                    className={`capitalize ${order.status === 'completed' ? 'bg-green-600 hover:bg-green-700' : ''}`}
                                                >
                                                    {order.status}
                                                </Badge>
                                            </TableCell>
                                            <TableCell className="text-right">
                                                {order.status === 'pending' && (
                                                    <Button size="sm" onClick={() => handleOpenPayment(order)} className="gap-2">
                                                        <CreditCard className="h-4 w-4" /> Pay
                                                    </Button>
                                                )}
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
                            onClick={() => setPage(p => Math.max(1, p - 1))}
                            disabled={page === 1}
                        >
                            Previous
                        </Button>
                        <div className="text-sm font-medium">Page {page} of {pagination.total_pages}</div>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setPage(p => Math.min(pagination.total_pages, p + 1))}
                            disabled={page === pagination.total_pages}
                        >
                            Next
                        </Button>
                    </div>
                )}
            </div>

            {/* Payment Dialog */}
            <Dialog open={isPaymentOpen} onOpenChange={setIsPaymentOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>Complete Payment</DialogTitle>
                        <DialogDescription>
                            Process payment for Queue <b>#{selectedOrder?.queue_number}</b>
                        </DialogDescription>
                    </DialogHeader>

                    <div className="py-4">
                        <div className="mb-6 flex flex-col items-center justify-center rounded-lg border bg-muted/30 p-4">
                            <span className="text-sm text-muted-foreground">Total Amount</span>
                            <span className="text-3xl font-bold text-foreground">{formatRupiah(selectedOrder?.payment_amount || 0)}</span>
                        </div>

                        <label className="text-sm font-medium mb-2 block">Select Payment Method</label>
                        <Tabs value={selectedPaymentMethod ? String(selectedPaymentMethod) : undefined} onValueChange={(v) => setSelectedPaymentMethod(Number(v))} className="w-full">
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
                                                <label className="text-sm font-medium">Cash Received</label>
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
                                            <div className="flex justify-between items-center bg-background p-3 rounded border shadow-sm">
                                                <span className="text-sm font-medium text-muted-foreground">Change Due</span>
                                                <span className={`text-xl font-bold ${Number(cashReceived) >= (selectedOrder?.payment_amount || 0) ? 'text-green-600' : 'text-red-500'}`}>
                                                    {formatRupiah(Math.max(0, Number(cashReceived) - (selectedOrder?.payment_amount || 0)))}
                                                </span>
                                            </div>
                                        </div>
                                    )}
                                </TabsContent>
                            ))}
                        </Tabs>
                    </div>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsPaymentOpen(false)}>Cancel</Button>
                        <Button onClick={handlePayment} disabled={!selectedPaymentMethod || completeManualPaymentMutation.isPending}>
                            {completeManualPaymentMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            Confirm Payment
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    )
}
