import { useEffect, useState, useMemo } from 'react'
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
    useCheckoutOrderMutation,
    useCalculateOrderMutation,
    useInitiateMidtransPaymentMutation,
    orderDetailQueryOptions
} from '@/lib/api/query/orders'
import { usePrinterSettingsQuery } from '@/lib/api/query/settings'
import { usePaymentMethodsListQuery } from '@/lib/api/query/payment-methods'
import { usePromotionsListQuery } from '@/lib/api/query/promotions'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import { useQueryClient, useQuery } from '@tanstack/react-query'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { printerService } from "@/lib/printer"

interface CheckoutDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    cart: any[]
    selectedOrderType: 'dine_in' | 'takeaway'
    selectedCustomerId?: string | null
    onPaymentSuccess?: () => void
}

export function CheckoutDialog({ open, onOpenChange, cart = [], selectedOrderType, selectedCustomerId, onPaymentSuccess }: CheckoutDialogProps) {
    const { t } = useTranslation()

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto overflow-x-hidden">
                <DialogHeader>
                    <DialogTitle>{t('order.payment_dialog.title')}</DialogTitle>
                    <DialogDescription>{t('order.payment_dialog.desc')}</DialogDescription>
                </DialogHeader>

                {open && <CheckoutDialogForm 
                    open={open} 
                    onOpenChange={onOpenChange} 
                    cart={cart}
                    selectedOrderType={selectedOrderType}
                    selectedCustomerId={selectedCustomerId}
                    onPaymentSuccess={onPaymentSuccess}
                />}
            </DialogContent>
        </Dialog>
    )
}

function CheckoutDialogForm({ open, onOpenChange, cart, selectedOrderType, selectedCustomerId, onPaymentSuccess }: CheckoutDialogProps) {
    const { t } = useTranslation()
    const queryClient = useQueryClient()

    const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<number | undefined>(undefined)
    const [cashReceived, setCashReceived] = useState<string>('')
    const [qrisUrl, setQrisUrl] = useState<string | null>(null)
    const [selectedPromoId, setSelectedPromoId] = useState<string | null>(null)
    const [createdOrderId, setCreatedOrderId] = useState<string | null>(null)

    const checkoutOrderMutation = useCheckoutOrderMutation()
    const calculateOrderMutation = useCalculateOrderMutation()
    const initiateMidtransPaymentMutation = useInitiateMidtransPaymentMutation()

    const { data: paymentMethods } = usePaymentMethodsListQuery()
    const { data: printerSettings } = usePrinterSettingsQuery()
    const { data: promotionsData } = usePromotionsListQuery({ limit: 100, trash: false })
    const activePromotions = promotionsData?.promotions?.filter(p => p.is_active) || []

    const itemsForApi = useMemo(() => {
        return cart.map(item => ({
            product_id: item.product.id!,
            quantity: item.quantity,
            options: item.variant && item.variant.id ? [{ product_option_id: item.variant.id }] : []
        }))
    }, [cart])

    useEffect(() => {
        if (open && itemsForApi.length > 0) {
            calculateOrderMutation.mutate({
                items: itemsForApi,
                promotion_id: selectedPromoId !== "none" && selectedPromoId ? selectedPromoId : undefined
            })
        }
    }, [open, itemsForApi, selectedPromoId])

    useEffect(() => {
        if (!open) {
            setCreatedOrderId(null)
            setQrisUrl(null)
        }
    }, [open])

    useEffect(() => {
        if (paymentMethods && paymentMethods.length > 0 && !selectedPaymentMethod) {
            setSelectedPaymentMethod(paymentMethods[0].id)
        }
    }, [paymentMethods, selectedPaymentMethod])

    const calculationData = calculateOrderMutation.data

    const { data: orderData } = useQuery({
        ...orderDetailQueryOptions(createdOrderId!),
        enabled: !!createdOrderId && !!qrisUrl && open,
        refetchInterval: open && qrisUrl ? 3000 : false
    })

    // Auto close
    useEffect(() => {
        if (open && orderData?.is_paid) {
            if (printerSettings?.auto_print) {
                handlePrint(orderData.id || '')
            }
            if (onPaymentSuccess) onPaymentSuccess()
            onOpenChange(false)
            toast.success(t('order.success.payment_success'), {
                description: t('order.payment_dialog.midtrans_auto_confirm'),
                style: { background: '#10B981', color: 'white', border: 'none' },
                action: { label: t('common.print'), onClick: () => handlePrint(orderData.id || '') }
            })
        }
    }, [orderData?.is_paid, open])

    const handlePrint = async (orderId: string) => {
        try {
            await printerService.printInvoice(orderId)
            toast.success(t('settings.printer.print_success', { defaultValue: 'Print command sent' }))
        } catch (error) {
            console.error(error)
            toast.error(t('settings.printer.print_failed', { defaultValue: 'Failed to print receipt' }))
        }
    }

    const handleCheckout = async () => {
        if (!selectedPaymentMethod) {
            toast.error(t('order.errors.select_payment'))
            return
        }

        const totalAmount = calculationData?.net_total || 0
        const method = paymentMethods?.find(m => m.id === selectedPaymentMethod)
        const isCash = method?.name?.toLowerCase().includes('cash')

        let cashValue = 0
        if (isCash) {
            cashValue = Number(cashReceived)
            if (cashValue < totalAmount) {
                toast.error(t('order.errors.cash_insufficient'))
                return
            }
        }

        const requestBody = {
            type: selectedOrderType,
            items: itemsForApi,
            customer_id: selectedCustomerId || undefined,
            promotion_id: selectedPromoId !== "none" && selectedPromoId ? selectedPromoId : undefined,
            payment_method_id: selectedPaymentMethod,
            cash_received: isCash ? cashValue : totalAmount,
        }

        try {
            const result = await checkoutOrderMutation.mutateAsync(requestBody as any)
            setCreatedOrderId(result.id || null)
            await queryClient.invalidateQueries({ queryKey: ['orders'] })
            
            if (method?.id === 3 || isCash) {
                // Static QR or Cash pays immediately
                if (printerSettings?.auto_print) handlePrint(result.id || '')
                if (onPaymentSuccess) onPaymentSuccess()
                onOpenChange(false)
                
                if (isCash) {
                    const change = cashValue - totalAmount
                    toast.success(`${t('order.success.payment_success')} ${formatRupiah(change)}`, {
                        duration: 5000,
                        description: `${t('order.success.received')}: ${formatRupiah(cashValue)} | ${t('order.total')}: ${formatRupiah(totalAmount)}`,
                        closeButton: true, position: 'top-center',
                        style: { background: '#10B981', color: 'white', border: 'none' },
                        action: { label: t('common.print'), onClick: () => handlePrint(result.id || '') }
                    })
                } else {
                    toast.success(t('order.success.payment_complete'), { action: { label: t('common.print'), onClick: () => handlePrint(result.id || '') } })
                }
            } else {
                // Dynamic QRIS
                const response = await initiateMidtransPaymentMutation.mutateAsync({ id: result.id || '' })
                const generateQrAction = response.actions?.find((action: any) => action.name === 'generate-qr-code')
                if (generateQrAction && generateQrAction.url) {
                    setQrisUrl(generateQrAction.url)
                    toast.success(t('order.payment_dialog.qris_generated'))
                } else {
                    toast.error(t('order.payment_dialog.generate_qris_failed'))
                }
            }
        } catch (error: any) {
            console.error(error)
            toast.error(error?.response?.data?.message || 'Gagal memproses pesanan')
        }
    }

    if (!calculationData && calculateOrderMutation.isPending) {
        return <div className="h-48 flex items-center justify-center"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>
    }

    return (
        <>
            <div className="bg-muted/30 p-4 rounded-lg space-y-3 mb-2">
                <div className="flex justify-between items-center text-sm">
                    <span className="text-muted-foreground">{t('order.subtotal')}</span>
                    <span>{formatRupiah(calculationData?.gross_total || 0)}</span>
                </div>
                <div className="flex justify-between items-center text-sm">
                    <span className="text-muted-foreground flex items-center gap-1"><TicketPercent className="w-3 h-3" /> {t('order.payment_dialog.discount')}</span>
                    <span className="text-primary">-{formatRupiah(calculationData?.discount_amount || 0)}</span>
                </div>
                {calculationData?.tax_amount !== undefined && calculationData.tax_amount > 0 && (
                    <div className="flex justify-between items-center text-sm">
                        <span className="text-muted-foreground">{t('order.tax_amount')}</span>
                        <span>{formatRupiah(calculationData.tax_amount)}</span>
                    </div>
                )}
                <div className="flex items-center gap-2 pt-2 border-t border-dashed">
                    <Select
                        value={selectedPromoId || "none"}
                        onValueChange={(val) => setSelectedPromoId(val)}
                    >
                        <SelectTrigger className="h-8 text-sm w-full"><SelectValue placeholder={t('order.payment_dialog.select_promo')} /></SelectTrigger>
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
                <div className="flex justify-between items-center font-bold text-sm border-t pt-2">
                    <span>{t('order.total')}</span>
                    <span className="text-primary">{formatRupiah(calculationData?.net_total || 0)}</span>
                </div>
            </div>

            <div className="space-y-4">
                <Tabs defaultValue={paymentMethods?.[0]?.name} onValueChange={(val) => {
                    setSelectedPaymentMethod(paymentMethods?.find(m => m.name === val)?.id)
                    setQrisUrl(null)
                }} className="w-full">
                    <TabsList className="grid w-full grid-cols-3">
                        {paymentMethods?.map((method) => (<TabsTrigger key={method.id} value={method.name || ''} className="text-sm">{method.name}</TabsTrigger>))}
                    </TabsList>
                    <div className="mt-4">
                        {paymentMethods?.map((method) => (
                            <TabsContent key={method.id} value={method.name || ''}>
                                <div className="flex flex-col items-center p-4 border rounded-lg bg-muted/10 gap-4">
                                    {method.name?.toLowerCase().includes('cash') && (
                                        <div className="w-full max-w-xs space-y-4 pt-2">
                                            <div className="space-y-2">
                                                <label className="text-sm font-medium">{t('order.payment_dialog.cash_received')}</label>
                                                <Input
                                                    type="text" inputMode="numeric"
                                                    value={cashReceived ? Number(cashReceived).toLocaleString('id-ID') : ''}
                                                    onChange={(e) => setCashReceived(e.target.value.replace(/\D/g, ''))}
                                                    className="text-center text-xl font-bold h-12"
                                                    placeholder={t('order.payment_dialog.enter_amount')}
                                                />
                                            </div>
                                            <div className="flex justify-between items-center text-sm py-3 rounded-lg">
                                                <span className="text-muted-foreground">{t('order.payment_dialog.change')}</span>
                                                <span className="font-bold text-sm text-primary">
                                                    {formatRupiah(Math.max(0, Number(cashReceived) - (calculationData?.net_total || 0)))}
                                                </span>
                                            </div>
                                        </div>
                                    )}

                                    {method.id === 3 && (
                                        <div className="bg-white p-2 rounded-lg mt-2 mx-auto flex flex-col items-center justify-center">
                                            <div className="h-48 w-48 mb-2"><img src={`https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=STATIC_QRIS`} alt={t('order.static_qris')} className="w-full h-full object-contain" /></div>
                                            <span className="text-sm text-muted-foreground font-medium text-center">{t('order.payment_dialog.scan_qr_static')}</span>
                                        </div>
                                    )}

                                    {method.name?.toLowerCase().includes('qris') && method.id !== 3 && (
                                        <div className="bg-white p-2 rounded-lg mt-2 mx-auto flex items-center justify-center">
                                            {qrisUrl ? (
                                                <div className="flex flex-col items-center gap-2">
                                                    <div className="h-48 w-48"><img src={qrisUrl} alt={t('order.qr_code_alt')} className="w-full h-full object-contain" /></div>
                                                    <span className="text-sm text-muted-foreground font-medium">{t('order.payment_dialog.scan_qr_dynamic')}</span>
                                                    <p className="text-xs text-muted-foreground break-all max-w-[200px] text-center mt-1 select-all">{qrisUrl}</p>
                                                </div>
                                            ) : (
                                                <div className="flex flex-col items-center gap-2 py-4">
                                                    <span className="text-sm text-muted-foreground text-center">Menunggu pesanan dibuat untuk memunculkan QRIS</span>
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

            <DialogFooter className="gap-2 sm:justify-between border-t pt-4">
                <div className="flex gap-2 w-full">
                    <Button type="button" variant="outline" className="flex-1" onClick={() => onOpenChange(false)}>
                        {createdOrderId ? t('common.close', { defaultValue: 'Tutup' }) : t('order.payment_dialog.cancel')}
                    </Button>

                    {!createdOrderId && (
                        <Button type="button" className="flex-1" onClick={handleCheckout} disabled={checkoutOrderMutation.isPending || initiateMidtransPaymentMutation.isPending || (paymentMethods?.find(m => m.id === selectedPaymentMethod)?.name?.toLowerCase().includes('cash') && !cashReceived)}>
                            {checkoutOrderMutation.isPending || initiateMidtransPaymentMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Banknote className="h-4 w-4 mr-2" />}
                            {t('order.payment_dialog.pay')}
                        </Button>
                    )}
                </div>
            </DialogFooter>
        </>
    )
}
