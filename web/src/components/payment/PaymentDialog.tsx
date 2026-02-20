import { useEffect, useState } from 'react'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Loader2, Banknote, TicketPercent } from 'lucide-react'
import { formatRupiah } from '@/lib/utils'
import {
    useOrderDetailQuery,
    useConfirmManualPaymentMutation,
    useCancelOrderMutation,
    useInitiateMidtransPaymentMutation,
    useApplyPromotionMutation
} from '@/lib/api/query/orders'
import { usePrinterSettingsQuery } from '@/lib/api/query/settings'
import { usePaymentMethodsListQuery } from '@/lib/api/query/payment-methods'
import { usePromotionsListQuery } from '@/lib/api/query/promotions'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { printerService } from "@/lib/printer"


interface PaymentDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    orderId: string | null
    onPaymentSuccess?: () => void
    mode?: 'checkout' | 'payment'
    onCancelOrder?: () => void
}

export function PaymentDialog({ open, onOpenChange, orderId, onPaymentSuccess, mode = 'checkout', onCancelOrder }: PaymentDialogProps) {
    const { t } = useTranslation()
    const queryClient = useQueryClient()
    const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<number | undefined>(undefined)
    const [cashReceived, setCashReceived] = useState<string>('')
    const [qrisUrl, setQrisUrl] = useState<string | null>(null)
    const [cancelDialogOpen, setCancelDialogOpen] = useState(false)

    const { data: order, isLoading: isLoadingOrder } = useOrderDetailQuery(orderId || '', {
        refetchInterval: open ? 3000 : false,
        enabled: !!orderId && open
    })

    const { data: paymentMethods } = usePaymentMethodsListQuery()
    const confirmManualPaymentMutation = useConfirmManualPaymentMutation()
    const initiateMidtransPaymentMutation = useInitiateMidtransPaymentMutation()
    const cancelOrderMutation = useCancelOrderMutation()
    const { data: printerSettings } = usePrinterSettingsQuery()

    // Promo states mostly for display if we keep the promo selection here
    // For now, let's assuming promo is already applied or we include the promo selector here too.
    // Based on existing code, promo is selected inside the dialog.
    const { data: promotionsData } = usePromotionsListQuery({ limit: 100, trash: false })
    const activePromotions = promotionsData?.promotions?.filter(p => p.is_active) || []
    const applyPromotionMutation = useApplyPromotionMutation()


    // Reset state when dialog opens
    useEffect(() => {
        if (open) {
            setQrisUrl(null)
            setCashReceived('')
            setSelectedPaymentMethod(undefined) // Optional: create default
            // If order has method, maybe pre-select?
        }
    }, [open])

    // Auto close
    useEffect(() => {
        if (open && order?.status === 'paid') {
            if (printerSettings?.auto_print) {
                handlePrint()
            }

            if (onPaymentSuccess) onPaymentSuccess()
            onOpenChange(false)
            toast.success(t('order.success.payment_success'), {
                description: t('order.payment_dialog.midtrans_auto_confirm'),
                style: { background: '#10B981', color: 'white', border: 'none' },
                action: {
                    label: 'Print',
                    onClick: () => handlePrint()
                }
            })
        }
    }, [order?.status, open, onPaymentSuccess, onOpenChange, t, printerSettings, orderId])

    const handleAttemptClose = (isOpen: boolean) => {
        if (!isOpen) {
            if (mode === 'checkout' && orderId && order?.status === 'open') {
                setCancelDialogOpen(true)
            } else {
                onOpenChange(false)
            }
        } else {
            onOpenChange(true)
        }
    }

    const handleConfirmCancel = async () => {
        if (!orderId) return

        try {
            await cancelOrderMutation.mutateAsync({
                id: orderId,
                body: {
                    cancellation_reason_id: 1,
                    cancellation_notes: t('order.payment_dialog.cancelled_by_user')
                }
            })
            onOpenChange(false)
            setCancelDialogOpen(false)
        } catch (error) {
            console.error(error)
        }
    }

    const handlePrint = async () => {
        if (!orderId) return
        try {
            await printerService.printInvoice(orderId)
            toast.success(t('payment.print_success', { defaultValue: 'Print command sent' }))
        } catch (error) {
            console.error(error)
            toast.error(t('payment.print_failed', { defaultValue: 'Failed to print receipt' }))
        }
    }

    const handlePayment = async () => {
        if (!orderId || !selectedPaymentMethod) {
            if (!selectedPaymentMethod) toast.error(t('order.errors.select_payment'))
            return
        }

        const totalAmount = order?.net_total || 0
        const method = paymentMethods?.find(m => m.id === selectedPaymentMethod)
        const isCash = method?.name?.toLowerCase().includes('cash')

        // Static QRIS (ID 3)
        if (selectedPaymentMethod === 3) {
            const payload = {
                payment_method_id: selectedPaymentMethod,
                cash_received: totalAmount
            }

            try {
                await confirmManualPaymentMutation.mutateAsync({
                    id: orderId,
                    body: payload
                })

                await queryClient.invalidateQueries({ queryKey: ['orders'] })

                if (printerSettings?.auto_print) {
                    handlePrint()
                }

                if (onPaymentSuccess) onPaymentSuccess()
                onOpenChange(false)

                toast.success(t('order.success.payment_complete'), {
                    action: {
                        label: 'Print',
                        onClick: () => handlePrint()
                    }
                })
            } catch (error) {
                console.error(error)
            }
            return
        }


        if (isCash) {
            const inputCash = Number(cashReceived)
            if (inputCash < totalAmount) {
                toast.error(t('order.errors.cash_insufficient'))
                return
            }

            const payload = {
                payment_method_id: selectedPaymentMethod,
                cash_received: inputCash
            }

            try {
                await confirmManualPaymentMutation.mutateAsync({
                    id: orderId,
                    body: payload
                })
                await queryClient.invalidateQueries({ queryKey: ['orders'] })

                if (printerSettings?.auto_print) {
                    handlePrint()
                }

                if (onPaymentSuccess) onPaymentSuccess()
                onOpenChange(false)

                const change = inputCash - totalAmount
                toast.success(`${t('order.success.payment_success')} ${formatRupiah(change)}`, {
                    duration: 5000,
                    description: `${t('order.success.received')}: ${formatRupiah(inputCash)} | ${t('order.total')}: ${formatRupiah(totalAmount)}`,
                    closeButton: true,
                    position: 'top-center',
                    style: { background: '#10B981', color: 'white', border: 'none' },
                    action: {
                        label: 'Print',
                        onClick: () => handlePrint()
                    }
                })
            } catch (error) {
                console.error(error)
            }
            return
        }

        // Dynamic QRIS
        try {
            const response = await initiateMidtransPaymentMutation.mutateAsync({
                id: orderId
            })
            const actions = response.actions
            const generateQrAction = actions?.find((action: any) => action.name === 'generate-qr-code')

            if (generateQrAction && generateQrAction.url) {
                setQrisUrl(generateQrAction.url)
                toast.success(t('order.payment_dialog.qris_generated'))
            } else {
                toast.error(t('order.payment_dialog.generate_qris_failed'))
            }
        } catch (error) {
            console.error(error)
        }
    }

    return (
        <>
            <Dialog open={open} onOpenChange={handleAttemptClose}>
                <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto overflow-x-hidden">
                    <DialogHeader>
                        <DialogTitle>{t('order.payment_dialog.title')}</DialogTitle>
                        <DialogDescription>{t('order.payment_dialog.desc')}</DialogDescription>
                    </DialogHeader>

                    {isLoadingOrder ? (
                        <div className="h-48 flex items-center justify-center">
                            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                        </div>
                    ) : (
                        <>
                            {/* Summary & Promo Section - Copied from OrderPage */}
                            <div className="bg-muted/30 p-4 rounded-lg space-y-3 mb-2">
                                <div className="flex justify-between items-center text-sm">
                                    <span className="text-muted-foreground">{t('order.subtotal')}</span>
                                    <span>{formatRupiah(order?.gross_total || 0)}</span>
                                </div>
                                <div className="flex justify-between items-center text-sm">
                                    <span className="text-muted-foreground flex items-center gap-1"><TicketPercent className="w-3 h-3" /> {t('order.payment_dialog.discount')}</span>
                                    <span className="text-green-600">-{formatRupiah(order?.discount_amount || 0)}</span>
                                </div>
                                {/* Promo Select */}
                                <div className="flex items-center gap-2 pt-2 border-t border-dashed">
                                    <Select
                                        value={order?.applied_promotion_id || "none"}
                                        onValueChange={(val) => {
                                            if (!orderId) return;
                                            if (val === "none") {

                                            } else {
                                                applyPromotionMutation.mutate({
                                                    id: orderId,
                                                    body: { promotion_id: val }
                                                })
                                            }
                                        }}
                                    >
                                        <SelectTrigger className="h-8 text-xs w-full">
                                            <SelectValue placeholder={t('order.payment_dialog.select_promo')} />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="none">{t('order.payment_dialog.no_promo')}</SelectItem>
                                            {activePromotions.map(p => (
                                                <SelectItem key={p.id} value={p.id || ''}>
                                                    {p.name} - {p.discount_type === 'percentage' ? `${p.discount_value}%` : formatRupiah(Number(p.discount_value))}
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>

                                <div className="flex justify-between items-center font-bold text-lg border-t pt-2">
                                    <span>{t('order.total')}</span>
                                    <span className="text-primary">{formatRupiah(order?.net_total || 0)}</span>
                                </div>
                            </div>

                            {/* Payment Method Tabs */}
                            <div className="space-y-4">
                                <Tabs
                                    defaultValue={paymentMethods?.[0]?.name}
                                    onValueChange={(val) => {
                                        const method = paymentMethods?.find(m => m.name === val)
                                        setSelectedPaymentMethod(method?.id)
                                        setQrisUrl(null)
                                    }}
                                    className="w-full"
                                >
                                    <TabsList className="grid w-full grid-cols-3">
                                        {paymentMethods?.map((method) => (
                                            <TabsTrigger key={method.id} value={method.name || ''} className="text-xs">
                                                {method.name}
                                            </TabsTrigger>
                                        ))}
                                    </TabsList>
                                    <div className="mt-4">
                                        {paymentMethods?.map((method) => (
                                            <TabsContent key={method.id} value={method.name || ''}>
                                                <div className="flex flex-col items-center p-4 border rounded-lg bg-muted/10 gap-4">
                                                    {/* Same Logic as Order Page */}
                                                    {method.name?.toLowerCase().includes('cash') && (
                                                        // Cash UI
                                                        <div className="w-full max-w-xs space-y-4 pt-2">
                                                            <div className="space-y-2">
                                                                <label className="text-sm font-medium">{t('order.payment_dialog.cash_received')}</label>
                                                                <Input
                                                                    autoFocus
                                                                    type="text"
                                                                    inputMode="numeric"
                                                                    value={cashReceived ? Number(cashReceived).toLocaleString('id-ID') : ''}
                                                                    onChange={(e) => {
                                                                        const val = e.target.value.replace(/\D/g, '')
                                                                        setCashReceived(val)
                                                                    }}
                                                                    className="text-center text-xl font-bold h-12"
                                                                    placeholder={t('order.payment_dialog.enter_amount')}
                                                                />
                                                            </div>
                                                            <div className="flex justify-between items-center text-sm py-3  rounded-lg ">
                                                                <span className="text-muted-foreground">{t('order.payment_dialog.change')}</span>
                                                                <span className="font-bold text-lg text-primary">
                                                                    {formatRupiah(Math.max(0, Number(cashReceived) - (order?.net_total || 0)))}
                                                                </span>
                                                            </div>
                                                        </div>
                                                    )}

                                                    {/* Static QRIS (ID 3) */}
                                                    {method.id === 3 && (
                                                        <div className="bg-white p-2 rounded-lg mt-2 mx-auto flex flex-col items-center justify-center">
                                                            <div className="h-48 w-48 mb-2">
                                                                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=PAY_ORDER_${orderId}`} alt="Static QRIS" className="w-full h-full object-contain" />
                                                            </div>
                                                            <span className="text-xs text-muted-foreground font-medium text-center">{t('order.payment_dialog.scan_qr_static')}</span>
                                                            <Button size="sm" className="mt-2 w-full" onClick={handlePayment}>
                                                                {confirmManualPaymentMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                                                                {t('order.payment_dialog.confirm_payment')}
                                                            </Button>
                                                        </div>
                                                    )}

                                                    {/* Dynamic QRIS */}
                                                    {method.name?.toLowerCase().includes('qris') && method.id !== 3 && (
                                                        <div className="bg-white p-2 rounded-lg mt-2 mx-auto flex items-center justify-center">
                                                            {qrisUrl ? (
                                                                <div className="flex flex-col items-center gap-2">
                                                                    <div className="h-48 w-48">
                                                                        <img src={qrisUrl} alt={t('order.qr_code_alt')} className="w-full h-full object-contain" />
                                                                    </div>

                                                                    <span className="text-xs text-muted-foreground font-medium">{t('order.payment_dialog.scan_qr_dynamic')}</span>
                                                                    <p className="text-[10px] text-muted-foreground break-all max-w-[200px] text-center mt-1 select-all">{qrisUrl}</p>
                                                                </div>
                                                            ) : (
                                                                <div className="flex flex-col items-center gap-2 py-4">
                                                                    <Button size="sm" onClick={handlePayment} disabled={initiateMidtransPaymentMutation.isPending}>
                                                                        {initiateMidtransPaymentMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                                                                        {t('order.payment_dialog.generate_qris')}
                                                                    </Button>
                                                                    <span className="text-xs text-muted-foreground text-center">{t('order.payment_dialog.click_to_generate')}</span>
                                                                </div>
                                                            )}
                                                        </div>
                                                    )}
                                                </div>
                                            </TabsContent>
                                        ))}
                                    </div>
                                </Tabs>
                            </div>
                        </>
                    )}

                    <DialogFooter className="gap-2 sm:justify-between">
                        <div className="flex gap-2 w-full">
                            {mode === 'payment' && (
                                <Button type="button" variant="destructive" onClick={onCancelOrder ? onCancelOrder : () => setCancelDialogOpen(true)}>
                                    {t('order.payment_dialog.cancel_order_btn')}
                                </Button>
                            )}
                            <Button type="button" variant="outline" className="flex-1" onClick={() => handleAttemptClose(false)}>
                                {mode === 'checkout' ? t('order.payment_dialog.cancel') : t('order.payment_dialog.close')}
                            </Button>

                            {paymentMethods?.find(m => m.id === selectedPaymentMethod)?.name?.toLowerCase().includes('cash') && (
                                <Button type="button" className="flex-1" onClick={handlePayment} disabled={!cashReceived}>
                                    {confirmManualPaymentMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Banknote className="h-4 w-4 mr-2" />}
                                    {t('order.payment_dialog.pay')}
                                </Button>
                            )}
                        </div>
                    </DialogFooter>

                </DialogContent>
            </Dialog>

            <AlertDialog open={cancelDialogOpen} onOpenChange={setCancelDialogOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>{t('order.cancel_confirm_title')}</AlertDialogTitle>
                        <AlertDialogDescription>
                            {t('order.cancel_confirm_desc')}
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel onClick={() => setCancelDialogOpen(false)}>{t('common.cancel')}</AlertDialogCancel>
                        <AlertDialogAction onClick={handleConfirmCancel} className="bg-destructive hover:bg-destructive/90">
                            {t('order.cancel_confirm_btn')}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}
