import { ShoppingCart, Trash2, Minus, Plus, User } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { formatRupiah } from '@/lib/utils'
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Product } from '@/lib/api/query/products'
import { InternalProductsProductOptionResponse, POSKasirInternalOrdersRepositoryOrderType } from '@/lib/api/generated'

export interface CartItem {
    product: Product
    variant?: InternalProductsProductOptionResponse
    quantity: number
}

interface CartContentProps {
    cart: CartItem[]
    t: any
    customers: any[]
    selectedCustomerId: string | null
    setSelectedCustomerId: (id: string | null) => void
    selectedOrderType: POSKasirInternalOrdersRepositoryOrderType
    setSelectedOrderType: (type: POSKasirInternalOrdersRepositoryOrderType) => void
    updateQuantity: (id: string, delta: number, variantId?: string) => void
    removeFromCart: (id: string, variantId?: string) => void
    calculateTotal: () => number
    handleCheckout: () => void
    canCheckout: boolean
}

export function CartContent({
    cart, t, customers, selectedCustomerId, setSelectedCustomerId,
    selectedOrderType, setSelectedOrderType, updateQuantity, removeFromCart,
    calculateTotal, handleCheckout, canCheckout
}: CartContentProps) {
    return (
        <div className="flex flex-col h-full overflow-hidden">
            <div className="p-4 border-b bg-muted/40 flex items-center gap-2 shrink-0">
                <ShoppingCart className="h-5 w-5" />
                <h2 className="font-semibold">{t('order.current_order')}</h2>
                <span className="ml-auto text-sm font-medium bg-primary/10 text-primary px-2 py-0.5 rounded-full">
                    {cart.length} {t('order.items')}
                </span>
            </div>

            <div className="mx-4 mt-4 space-y-3 shrink-0">
                <Select value={selectedCustomerId || 'walk_in'} onValueChange={(v) => setSelectedCustomerId(v === 'walk_in' ? null : v)}>
                    <SelectTrigger className="w-full">
                        <div className="flex items-center gap-2 truncate">
                            <User className="h-4 w-4 shrink-0" />
                            <SelectValue placeholder={t('order.select_customer', 'Walk-in Customer')} />
                        </div>
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="walk_in">{t('order.walk_in', 'Walk-in Customer')}</SelectItem>
                        {customers.map(c => (
                            <SelectItem key={c.id} value={c.id!}>{c.name}</SelectItem>
                        ))}
                    </SelectContent>
                </Select>

                <Tabs value={selectedOrderType} onValueChange={(v) => setSelectedOrderType(v as POSKasirInternalOrdersRepositoryOrderType)} className="w-full">
                    <TabsList className="grid w-full grid-cols-2">
                        <TabsTrigger value={POSKasirInternalOrdersRepositoryOrderType.OrderTypeDineIn}>{t('order.order_type.dine_in')}</TabsTrigger>
                        <TabsTrigger value={POSKasirInternalOrdersRepositoryOrderType.OrderTypeTakeaway}>{t('order.order_type.take_away')}</TabsTrigger>
                    </TabsList>
                </Tabs>
            </div>


            <div className="flex-1 min-h-0">
                <ScrollArea className="h-full">
                    <div className="p-4">
                        <div className="flex flex-col gap-3">
                            {cart.length === 0 ? (
                                <div className="h-32 flex flex-col items-center justify-center text-muted-foreground text-sm">
                                    <ShoppingCart className="h-8 w-8 mb-2 opacity-50" />
                                    <span>{t('order.empty_cart')}</span>
                                </div>
                            ) : (
                                cart.map((item) => (
                                    <div key={`${item.product.id}-${item.variant?.id || 'base'}`} className="flex gap-3 bg-background p-3 rounded-lg border group">
                                        <div className="h-12 w-12 rounded-md bg-muted overflow-hidden shrink-0">
                                            {item.product.image_url && <img src={item.product.image_url} className="h-full w-full object-cover" />}
                                        </div>
                                        <div className="flex-1 min-w-0 flex flex-col justify-between">
                                            <div className="flex justify-between items-start gap-1">
                                                <div className="min-w-0">
                                                    <span className="font-medium text-sm truncate leading-tight block">{item.product.name}</span>
                                                    {item.variant && (
                                                        <span className="text-sm text-muted-foreground block truncate">{item.variant.name} (+{formatRupiah(item.variant.additional_price || 0)})</span>
                                                    )}
                                                </div>
                                                <span className="text-sm font-bold ml-1">{formatRupiah(((item.product.price || 0) + (item.variant?.additional_price || 0)) * item.quantity)}</span>
                                            </div>
                                            <div className="flex items-center justify-between mt-1">
                                                <div className="text-sm text-muted-foreground">
                                                    {formatRupiah((item.product.price || 0) + (item.variant?.additional_price || 0))} x {item.quantity}
                                                </div>
                                                <div className="flex items-center gap-2">
                                                    <Button variant="ghost" size="icon" className="h-6 w-6 rounded-full" onClick={(e) => { e.stopPropagation(); updateQuantity(item.product.id!, -1, item.variant?.id) }}>
                                                        <Minus className="h-3 w-3" />
                                                    </Button>
                                                    <span className="w-4 text-center text-sm font-medium">{item.quantity}</span>
                                                    <Button variant="ghost" size="icon" className="h-6 w-6 rounded-full" onClick={(e) => { e.stopPropagation(); updateQuantity(item.product.id!, 1, item.variant?.id) }}>
                                                        <Plus className="h-3 w-3" />
                                                    </Button>
                                                    <Button variant="ghost" size="icon" className="h-6 w-6 rounded-full text-destructive hover:text-destructive hover:bg-destructive/10" onClick={(e) => { e.stopPropagation(); removeFromCart(item.product.id!, item.variant?.id) }}>
                                                        <Trash2 className="h-3 w-3" />
                                                    </Button>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                ))
                            )}
                        </div>
                    </div>
                </ScrollArea>
            </div>


            <div className="p-4 border-t bg-muted/20 space-y-4 shrink-0 mt-auto">
                <div className="space-y-1.5">
                    <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">{t('order.subtotal')}</span>
                        <span>{formatRupiah(calculateTotal())}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">{t('order.tax', 'Tax (11%)')}</span>
                        <span>{formatRupiah(Math.floor(calculateTotal() * 0.11))}</span>
                    </div>

                    <div className="flex justify-between text-sm font-bold border-t pt-2 mt-2">
                        <span>{t('order.total')}</span>
                        <span className="text-primary">{formatRupiah(calculateTotal() + Math.floor(calculateTotal() * 0.11))}</span>
                    </div>
                </div>
                <Button className="w-full h-12 text-sm " size="lg" disabled={cart.length === 0 || !canCheckout} onClick={handleCheckout}>
                    {t('order.charge')} {formatRupiah(calculateTotal() + Math.floor(calculateTotal() * 0.11))}
                </Button>
            </div>
        </div>
    )
}
