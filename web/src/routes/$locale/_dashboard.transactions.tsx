import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { ordersListQueryOptions, useUpdateOrderStatusMutation, useCancelOrderMutation, useRefundOrderMutation, useConfirmManualPaymentMutation, useInitiateMidtransPaymentMutation } from '@/lib/api/query/orders'
import { usersListQueryOptions } from '@/lib/api/query/user'
import { useCancellationReasonsListQuery } from '@/lib/api/query/cancel-reason'
import { useState, useMemo } from 'react'
import { toast } from 'sonner'
import { InternalOrdersOrderListResponse, POSKasirInternalOrdersRepositoryOrderStatus } from '@/lib/api/generated'
import { useTranslation } from 'react-i18next'
import { NewPagination } from "@/components/pagination.tsx";
import { PaymentDialog } from "@/components/payment/PaymentDialog"
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { TransactionFilters } from '@/components/transactions/TransactionFilters'
import { TransactionTable } from '@/components/transactions/TransactionTable'
import { CancellationDialog } from '@/components/transactions/CancellationDialog'
import { RefundDialog } from '@/components/transactions/RefundDialog'

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

type OrderWithItems = InternalOrdersOrderListResponse & {
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

    // Cancellation Dialog State
    const [isCancelDialogOpen, setIsCancelDialogOpen] = useState(false)
    const [orderToCancel, setOrderToCancel] = useState<OrderWithItems | null>(null)
    const [cancelReasonId, setCancelReasonId] = useState<string>('')
    const [cancelNotes, setCancelNotes] = useState('')

    const { data: cancellationReasons } = useCancellationReasonsListQuery()
    const updateOrderStatusMutation = useUpdateOrderStatusMutation()
    const cancelOrderMutation = useCancelOrderMutation()
    const refundOrderMutation = useRefundOrderMutation()
    const manualPaymentMutation = useConfirmManualPaymentMutation()
    const midtransMutation = useInitiateMidtransPaymentMutation()

    const canUpdateStatus = updateOrderStatusMutation.isAllowed
    const canCancel = cancelOrderMutation.isAllowed
    const canRefund = refundOrderMutation.isAllowed
    const canPay = manualPaymentMutation.isAllowed || midtransMutation.isAllowed

    // Refund Dialog State
    const [isRefundDialogOpen, setIsRefundDialogOpen] = useState(false)
    const [orderToRefund, setOrderToRefund] = useState<OrderWithItems | null>(null)
    const [refundReason, setRefundReason] = useState('')

    const [mutatingOrderId, setMutatingOrderId] = useState<string | null>(null)

    const handleStatusUpdate = async (id: string, newStatus: POSKasirInternalOrdersRepositoryOrderStatus) => {
        setMutatingOrderId(id)
        const promise = updateOrderStatusMutation.mutateAsync({
            id,
            body: { status: newStatus }
        })

        toast.promise(promise, {
            loading: t('common.loading', 'Updating status...'),
            success: t('order.status_updated', 'Status updated successfully'),
            error: (err: any) => {
                const msg = err.response?.data?.message || err.message || 'Failed to update status'
                return `Error: ${msg}`
            },
            finally: () => setMutatingOrderId(null)
        })

        try {
            await promise
        } catch (error) {
            console.error(error)
        }
    }

    const handleOpenPayment = (order: OrderWithItems) => {
        setSelectedOrder(order)
        setIsPaymentOpen(true)
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

    const handleOpenRefund = (order: OrderWithItems) => {
        setOrderToRefund(order)
        setRefundReason('')
        setIsRefundDialogOpen(true)
    }

    const handleConfirmRefund = async () => {
        if (!orderToRefund?.id || !refundReason) return;

        try {
            await refundOrderMutation.mutateAsync({
                id: orderToRefund.id,
                body: { reason: refundReason }
            })
            setIsRefundDialogOpen(false)
            setOrderToRefund(null)
        } catch (error) {
            console.error(error)
        }
    }

    const handleFinish = async (order: OrderWithItems) => {
        if (!order.id) return;
        handleStatusUpdate(order.id, POSKasirInternalOrdersRepositoryOrderStatus.OrderStatusPaid)
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

    return (
        <div className="flex flex-col gap-6  mx-auto w-full">
            <TransactionFilters 
                users={users}
                selectedUserId={selectedUserId}
                onUserChange={handleUserChange}
                t={t}
            />

            <div className="space-y-4">
                <Tabs value={selectedTab} onValueChange={handleTabChange} className="w-full">
                    <TabsList className="flex w-full overflow-x-auto justify-start h-auto p-1 gap-2 bg-muted/20">
                        <TabsTrigger value="active" className="shrink-0">{t('transactions.tabs.active')}</TabsTrigger>
                        <TabsTrigger value="paid" className="shrink-0">{t('transactions.tabs.history')}</TabsTrigger>
                        <TabsTrigger value="cancelled" className="shrink-0">{t('transactions.tabs.cancelled')}</TabsTrigger>
                        <TabsTrigger value="all" className="shrink-0">{t('transactions.status.all')}</TabsTrigger>
                    </TabsList>
                </Tabs>

                <div className="rounded-2xl border-0 shadow-sm bg-card overflow-hidden">
                    <TransactionTable 
                        orders={orders}
                        isLoading={isLoading}
                        userMap={userMap}
                        mutatingOrderId={mutatingOrderId}
                        handleStatusUpdate={handleStatusUpdate}
                        handleOpenPayment={handleOpenPayment}
                        handleOpenCancel={handleOpenCancel}
                        handleOpenRefund={handleOpenRefund}
                        handleFinish={handleFinish}
                        canUpdateStatus={canUpdateStatus}
                        canCancel={canCancel}
                        canRefund={canRefund}
                        canPay={canPay}
                        t={t}
                    />
                </div>

                {pagination && (
                    <NewPagination
                        pagination={pagination}
                        onClickPrev={() => handlePageChange((pagination.current_page || 1) - 1)}
                        onClickNext={() => handlePageChange((pagination.current_page || 1) + 1)}
                    />
                )}
            </div>

            {/* Payment Dialog Component */}
            <PaymentDialog
                open={isPaymentOpen}
                onOpenChange={setIsPaymentOpen}
                orderId={selectedOrder?.id || null}
                onPaymentSuccess={() => {
                    setIsPaymentOpen(false);
                    setSelectedOrder(null);
                }}
                mode="payment"
                onCancelOrder={() => {
                    setIsPaymentOpen(false)
                    if (selectedOrder) handleOpenCancel(selectedOrder)
                }}
            />

            {/* Cancellation Dialog */}
            <CancellationDialog 
                open={isCancelDialogOpen}
                onOpenChange={setIsCancelDialogOpen}
                order={orderToCancel}
                cancellationReasons={cancellationReasons}
                cancelReasonId={cancelReasonId}
                setCancelReasonId={setCancelReasonId}
                cancelNotes={cancelNotes}
                setCancelNotes={setCancelNotes}
                onConfirm={handleConfirmCancel}
                isPending={cancelOrderMutation.isPending}
                t={t}
            />

            {/* Refund Dialog */}
            <RefundDialog 
                open={isRefundDialogOpen}
                onOpenChange={setIsRefundDialogOpen}
                order={orderToRefund}
                refundReason={refundReason}
                setRefundReason={setRefundReason}
                onConfirm={handleConfirmRefund}
                isPending={refundOrderMutation.isPending}
                t={t}
            />
        </div >
    )
}

