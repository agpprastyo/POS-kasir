import { createFileRoute } from '@tanstack/react-router'
import { useState, useMemo, useEffect } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { useProductsListQuery, Product, productDetailQueryOptions } from '@/lib/api/query/products'
import { ProductCard } from '@/components/ProductCard'
import { Input } from '@/components/ui/input'
import { Search, ShoppingCart, Trash2, Banknote, Minus, Plus, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ScrollArea, ScrollBar } from '@/components/ui/scroll-area'
import { formatRupiah } from '@/lib/utils'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { useCreateOrderMutation, useCompleteManualPaymentMutation } from '@/lib/api/query/orders'
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { toast } from 'sonner'
import { POSKasirInternalDtoProductOptionResponse, POSKasirInternalRepositoryOrderType } from '@/lib/api/generated'
import { usePaymentMethodsListQuery } from '@/lib/api/query/payment-methods'


export const Route = createFileRoute('/_dashboard/order')({
    component: OrderPage,
})

interface CartItem {
    product: Product
    variant?: POSKasirInternalDtoProductOptionResponse
    quantity: number
}


function OrderPage() {
    const queryClient = useQueryClient()

    const [searchTerm, setSearchTerm] = useState('')
    const { data: productsData } = useProductsListQuery({ limit: 100, search: searchTerm })
    const products = productsData?.products || []

    const [selectedCategory, setSelectedCategory] = useState("all")
    const [selectedOrderType, setSelectedOrderType] = useState<POSKasirInternalRepositoryOrderType>(POSKasirInternalRepositoryOrderType.OrderTypeDineIn)

    const categories = useMemo(() => {
        const unique = new Map<string, string>()
        products.forEach(p => {
            if (p.category_id && p.category_name) {
                unique.set(String(p.category_id), p.category_name)
            }
        })
        return Array.from(unique.entries())
            .map(([id, name]) => ({ id, name }))
            .sort((a, b) => a.name.localeCompare(b.name))
    }, [products])

    const filteredProducts = useMemo(() => {
        if (selectedCategory === "all") return products
        return products.filter(p => String(p.category_id) === selectedCategory)
    }, [selectedCategory, products])

    const inStockProducts = filteredProducts.filter(p => (p.stock || 0) > 0)
    const outOfStockProducts = filteredProducts.filter(p => (p.stock || 0) <= 0)

    const [cart, setCart] = useState<CartItem[]>([])
    const [isPaymentOpen, setIsPaymentOpen] = useState(false)
    const { data: paymentMethods } = usePaymentMethodsListQuery()
    const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<number | undefined>(undefined)
    const [cashReceived, setCashReceived] = useState<string>('')
    const [createdOrderId, setCreatedOrderId] = useState<string | null>(null)
    const [isLoadingDetailsId, setIsLoadingDetailsId] = useState<string | null>(null)


    const [variantSelectionOpen, setVariantSelectionOpen] = useState(false)
    const [productForVariantSelection, setProductForVariantSelection] = useState<Product | null>(null)

    const createOrderMutation = useCreateOrderMutation()
    const completeManualPaymentMutation = useCompleteManualPaymentMutation()


    const addToCart = async (product: Product) => {
        if ((product.stock || 0) <= 0) {
            toast.error("Stok produk habis")
            return
        }
        setIsLoadingDetailsId(product.id || null)
        try {

            const detail = await queryClient.fetchQuery(productDetailQueryOptions(product.id!))

            if (detail.options && detail.options.length > 0) {
                setProductForVariantSelection(detail)
                setVariantSelectionOpen(true)
                return
            }
            addCartItem(detail)
        } catch (error) {
            console.error(error)
            toast.error("Gagal memuat info produk")
        } finally {
            setIsLoadingDetailsId(null)
        }
    }

    const addCartItem = (product: Product, variant?: POSKasirInternalDtoProductOptionResponse) => {
        setCart(prev => {
            const existing = prev.find(item =>
                item.product.id === product.id && item.variant?.id === variant?.id
            )

            const currentQty = existing ? existing.quantity : 0
            if (currentQty + 1 > (product.stock || 0)) {
                toast.error("Stok tidak mencukupi")
                return prev
            }

            if (existing) {
                return prev.map(item =>
                    item.product.id === product.id && item.variant?.id === variant?.id
                        ? { ...item, quantity: item.quantity + 1 }
                        : item
                )
            }
            return [...prev, { product, variant, quantity: 1 }]
        })
        setVariantSelectionOpen(false)
        setProductForVariantSelection(null)
    }

    const removeFromCart = (productId: string, variantId?: string) => {
        setCart(prev => prev.filter(item => !(item.product.id === productId && item.variant?.id === variantId)))
    }

    const updateQuantity = (productId: string, delta: number, variantId?: string) => {
        setCart(prev => {
            return prev.map(item => {
                if (item.product.id === productId && item.variant?.id === variantId) {
                    if (delta > 0 && item.quantity + delta > (item.product.stock || 0)) {
                        toast.error("Stok tidak mencukupi")
                        return item
                    }
                    const newQty = item.quantity + delta
                    return { ...item, quantity: Math.max(1, newQty) }
                }
                return item
            })
        })
    }

    const calculateTotal = () => {
        return cart.reduce((total, item) => {
            const price = (item.product.price || 0) + (item.variant?.additional_price || 0)
            return total + (price * item.quantity)
        }, 0)
    }

    const handleCheckout = () => {
        if (cart.length === 0) return

        const orderData = {
            items: cart.map(item => ({
                product_id: item.product.id!,
                product_option_id: item.variant?.id,
                quantity: item.quantity
            })),
            type: selectedOrderType
        }

        toast.promise(createOrderMutation.mutateAsync(orderData), {
            loading: 'Creating order...',
            success: (data) => {
                setCreatedOrderId(data.id)
                setIsPaymentOpen(true)
                return 'Order created! Proceed to payment.'
            },
            error: 'Failed to create order'
        })
    }

    useEffect(() => {
        if (isPaymentOpen) {
            setCashReceived('')
        }
    }, [isPaymentOpen])

    const handlePayment = async () => {
        if (!createdOrderId || !selectedPaymentMethod) {
            if (!selectedPaymentMethod) toast.error("Please select a payment method")
            return
        }

        const totalAmount = calculateTotal() * 1.11

        const method = paymentMethods?.find(m => m.id === selectedPaymentMethod)
        const isCash = method?.name?.toLowerCase().includes('cash')

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
                id: createdOrderId,
                body: payload
            })

            await queryClient.invalidateQueries({ queryKey: ['products', 'list'] })

            setIsPaymentOpen(false)
            setCart([])
            setCreatedOrderId(null)

            if (isCash) {
                const change = finalCashReceived - totalAmount
                toast.success(`Pembayaran Berhasil! Kembalian: ${formatRupiah(change)}`, {
                    duration: 5000,
                    description: `Diterima: ${formatRupiah(finalCashReceived)} | Total: ${formatRupiah(totalAmount)}`,
                    closeButton: true,
                    position: 'top-center',
                    style: { background: '#10B981', color: 'white', border: 'none' }
                })
            } else {
                toast.success("Payment completed successfully")
            }
        } catch (error) {
            console.error(error)
            toast.error("Payment failed")
        }
    }


    return (
        <div className="flex  h-[calc(100vh-4rem)] gap-2 ">

            {/* Left: Product Grid */}
            <div className="flex-1 flex flex-col gap-4 overflow-hidden   bg-background  min-h-0">
                <div className="px-4 ">
                    <div className="relative">
                        <Search className="absolute left-2.5 top-3 h-6 w-6 text-muted-foreground" />
                        <Input
                            type="search"
                            placeholder="Search products..."
                            className="pl-12 py-6"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                    </div>
                </div>

                <div className="px-4">
                    <Tabs value={selectedCategory} onValueChange={setSelectedCategory} className="w-full">
                        <TabsList className="w-full justify-start overflow-x-auto h-auto p-1 no-scrollbar">
                            <TabsTrigger value="all" className="rounded-full px-4">All</TabsTrigger>
                            {categories.map(category => (
                                <TabsTrigger key={category.id} value={category.id} className="rounded-full px-4">
                                    {category.name}
                                </TabsTrigger>
                            ))}
                        </TabsList>
                    </Tabs>
                </div>


                <div className="flex-1 min-h-0">
                    <ScrollArea className="h-full">
                        <div className="p-4">
                            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 pb-4">
                                {inStockProducts.map(product => (
                                    <div key={product.id} onClick={() => addToCart(product)} className="relative cursor-pointer transition-transform active:scale-95">
                                        <ProductCard
                                            product={product}
                                        />
                                        {isLoadingDetailsId === product.id && (
                                            <div className="absolute inset-0 bg-background/50 flex items-center justify-center rounded-xl z-10 backdrop-blur-[1px]">
                                                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                                            </div>
                                        )}
                                    </div>
                                ))}
                            </div>

                            {outOfStockProducts.length > 0 && (
                                <>
                                    <div className="relative py-4">
                                        <div className="absolute inset-0 flex items-center">
                                            <span className="w-full border-t" />
                                        </div>
                                        <div className="relative flex justify-center text-xs uppercase">
                                            <span className="bg-background px-2 text-muted-foreground">Out of Stock</span>
                                        </div>
                                    </div>
                                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 pb-4 opacity-60 grayscale">
                                        {outOfStockProducts.map(product => (
                                            <div key={product.id} className="relative cursor-not-allowed">
                                                <ProductCard product={product} />
                                            </div>
                                        ))}
                                    </div>
                                </>
                            )}

                            {products.length === 0 && (
                                <div className="h-40 flex items-center justify-center text-muted-foreground">
                                    No products found.
                                </div>
                            )}
                        </div>
                        <ScrollBar orientation="vertical" />
                    </ScrollArea>
                </div>

            </div>

            {/* Right: Cart Sidebar */}
            <div className="w-[350px] flex flex-col rounded-xl border bg-card  overflow-hidden min-h-0">
                <div className="p-4 border-b bg-muted/40 flex items-center gap-2 shrink-0">
                    <ShoppingCart className="h-5 w-5" />
                    <h2 className="font-semibold">Current Order</h2>
                    <span className="ml-auto text-xs font-medium bg-primary/10 text-primary px-2 py-0.5 rounded-full">
                        {cart.length} items
                    </span>
                </div>

                <div className="mx-4 mt-4">
                    <Tabs value={selectedOrderType} onValueChange={(v) => setSelectedOrderType(v as POSKasirInternalRepositoryOrderType)} className="w-full">
                        <TabsList className="grid w-full grid-cols-2">
                            <TabsTrigger value={POSKasirInternalRepositoryOrderType.OrderTypeDineIn}>Dine In</TabsTrigger>
                            <TabsTrigger value={POSKasirInternalRepositoryOrderType.OrderTypeTakeaway}>Take Away</TabsTrigger>
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
                                        <span>Cart is empty</span>
                                    </div>
                                ) : (
                                    cart.map((item) => (
                                        <div key={item.product.id} className="flex gap-3 bg-background p-3 rounded-lg border group">
                                            <div className="h-12 w-12 rounded-md bg-muted overflow-hidden shrink-0">
                                                {item.product.image_url && <img src={item.product.image_url} className="h-full w-full object-cover" />}
                                            </div>
                                            <div className="flex-1 min-w-0 flex flex-col justify-between">
                                                <div className="flex justify-between items-start gap-1">
                                                    <div className="min-w-0">
                                                        <span className="font-medium text-sm truncate leading-tight block">{item.product.name}</span>
                                                        {item.variant && (
                                                            <span className="text-xs text-muted-foreground block truncate">{item.variant.name} (+{formatRupiah(item.variant.additional_price || 0)})</span>
                                                        )}
                                                    </div>
                                                    <span className="text-sm font-bold ml-1">{formatRupiah(((item.product.price || 0) + (item.variant?.additional_price || 0)) * item.quantity)}</span>
                                                </div>
                                                <div className="flex items-center justify-between mt-1">
                                                    <div className="text-xs text-muted-foreground">
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
                        <ScrollBar orientation="vertical" />
                    </ScrollArea>
                </div>


                <div className="p-4 border-t bg-muted/20 space-y-4 shrink-0">
                    <div className="space-y-1.5">
                        <div className="flex justify-between text-sm">
                            <span className="text-muted-foreground">Subtotal</span>
                            <span>{formatRupiah(calculateTotal())}</span>
                        </div>
                        <div className="flex justify-between text-sm">
                            <span className="text-muted-foreground">Tax (11%)</span>
                            <span>{formatRupiah(calculateTotal() * 0.11)}</span>
                        </div>
                        <div className="flex justify-between text-lg font-bold border-t pt-2 mt-2">
                            <span>Total</span>
                            <span className="text-primary">{formatRupiah(calculateTotal() * 1.11)}</span>
                        </div>
                    </div>
                    <Button className="w-full h-12 text-lg " size="lg" disabled={cart.length === 0} onClick={handleCheckout}>
                        Charge {formatRupiah(calculateTotal() * 1.11)}
                    </Button>
                </div>
            </div>

            {/* Payment Dialog */}
            <Dialog open={isPaymentOpen} onOpenChange={setIsPaymentOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>Payment Method</DialogTitle>
                        <DialogDescription>
                            Select how the customer is paying. Total: <span className="font-bold text-foreground">{formatRupiah(calculateTotal() * 1.11)}</span>
                        </DialogDescription>
                    </DialogHeader>

                    <Tabs value={selectedPaymentMethod ? String(selectedPaymentMethod) : undefined} onValueChange={(v) => setSelectedPaymentMethod(Number(v))} className="w-full">
                        <TabsList className="flex flex-wrap h-auto w-full gap-2 bg-transparent p-0">
                            {paymentMethods?.map(method => (
                                <TabsTrigger
                                    key={method.id}
                                    value={String(method.id)}
                                    className="flex-1 min-w-[100px] border data-[state=active]:border-primary data-[state=active]:bg-primary/5"
                                >
                                    {method.name}
                                </TabsTrigger>
                            ))}
                        </TabsList>
                        <div className="py-4">
                            {paymentMethods?.map(method => (
                                <TabsContent key={method.id} value={String(method.id)} className="mt-0">
                                    <div className="flex flex-col items-center justify-center gap-4 py-4 border-2 border-dashed rounded-lg bg-muted/50">
                                        <Banknote className="h-10 w-10 text-muted-foreground" />
                                        <p className="text-sm text-muted-foreground">Process payment via {method.name}</p>

                                        {method.name?.toLowerCase().includes('cash') && (
                                            <div className="w-full max-w-xs space-y-4 pt-2">
                                                <div className="space-y-2">
                                                    <label className="text-sm font-medium">Cash Received</label>
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
                                                        placeholder="Enter amount"
                                                    />
                                                </div>
                                                <div className="flex justify-between items-center text-sm py-3  rounded-lg ">
                                                    <span className="text-muted-foreground">Change</span>
                                                    <span className="font-bold text-lg text-primary">
                                                        {formatRupiah(Math.max(0, Number(cashReceived) - (calculateTotal() * 1.11)))}
                                                    </span>
                                                </div>
                                            </div>
                                        )}

                                        {/* Simple placeholder for QRIS if name matches */}
                                        {method.name?.toLowerCase().includes('qris') && (
                                            <div className="h-32 w-32 bg-white p-2 rounded-lg mt-2">
                                                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=PAY_ORDER_${createdOrderId}`} alt="QR Code" className="w-full h-full" />
                                            </div>
                                        )}
                                    </div>
                                </TabsContent>
                            ))}
                        </div>
                    </Tabs>

                    <DialogFooter className="gap-2 sm:justify-between">
                        <div className="flex gap-2 w-full">
                            <Button type="button" variant="outline" className="flex-1" onClick={() => setIsPaymentOpen(false)}>
                                Cancel
                            </Button>
                            {selectedOrderType === POSKasirInternalRepositoryOrderType.OrderTypeDineIn && (
                                <Button type="button" variant="secondary" className="flex-1" onClick={() => {
                                    setIsPaymentOpen(false)
                                    setCart([])
                                    setCreatedOrderId(null)
                                    toast.success("Order saved! You can pay later in Transactions.")
                                }}>
                                    Pay Later
                                </Button>
                            )}
                            <Button type="button" className="flex-1" onClick={handlePayment}>
                                Complete Payment
                            </Button>
                        </div>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
            {/* Variant Selection Dialog */}
            <Dialog open={variantSelectionOpen} onOpenChange={setVariantSelectionOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>Select Variant</DialogTitle>
                        <DialogDescription>
                            Choose an option for <span className="font-semibold">{productForVariantSelection?.name}</span>
                        </DialogDescription>
                    </DialogHeader>

                    <div className="grid grid-cols-2 gap-4 py-4">
                        {/* Base Product Option */}
                        <div
                            className="flex flex-col gap-2 p-3 border rounded-lg cursor-pointer hover:bg-accent transition-colors"
                            onClick={() => addCartItem(productForVariantSelection!)}
                        >
                            <div className="aspect-square w-full rounded-md bg-muted overflow-hidden flex items-center justify-center">
                                {productForVariantSelection?.image_url ? (
                                    <img src={productForVariantSelection.image_url} alt={productForVariantSelection.name} className="w-full h-full object-cover" />
                                ) : (
                                    <span className="text-xs text-muted-foreground">Original</span>
                                )}
                            </div>
                            <div className="flex flex-col">
                                <span className="font-medium text-sm">Original</span>
                                <span className="text-xs text-muted-foreground">{formatRupiah(productForVariantSelection?.price || 0)}</span>
                            </div>
                        </div>

                        {productForVariantSelection?.options?.map(option => (
                            <div
                                key={option.id}
                                className="flex flex-col gap-2 p-3 border rounded-lg cursor-pointer hover:bg-accent transition-colors"
                                onClick={() => addCartItem(productForVariantSelection!, option)}
                            >
                                <div className="aspect-square w-full rounded-md bg-muted overflow-hidden">
                                    {option.image_url ? (
                                        <img src={option.image_url} alt={option.name} className="w-full h-full object-cover" />
                                    ) : (
                                        <div className="w-full h-full flex items-center justify-center text-muted-foreground bg-muted/50">
                                            <span className="text-xs">No Image</span>
                                        </div>
                                    )}
                                </div>
                                <div className="flex flex-col">
                                    <span className="font-medium text-sm">{option.name}</span>
                                    <span className="text-xs text-muted-foreground">+{formatRupiah(option.additional_price || 0)}</span>
                                </div>
                            </div>
                        ))}
                    </div>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setVariantSelectionOpen(false)}>Cancel</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    )
}
