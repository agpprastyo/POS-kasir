import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { useProductsListQuery, Product, productDetailQueryOptions } from '@/lib/api/query/products'
import { ProductCard } from '@/components/ProductCard'
import { Input } from '@/components/ui/input'
import { Search, ShoppingCart, Trash2, Minus, Plus, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ScrollArea, ScrollBar } from '@/components/ui/scroll-area'
import { formatRupiah } from '@/lib/utils'
import { PaymentDialog } from "@/components/payment/PaymentDialog"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { useCreateOrderMutation } from '@/lib/api/query/orders'
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { toast } from 'sonner'
import { InternalProductsProductOptionResponse, POSKasirInternalOrdersRepositoryOrderType } from '@/lib/api/generated'

const orderSearchSchema = z.object({
    category: z.string().optional().catch('all'),
})

export const Route = createFileRoute('/$locale/_dashboard/order')({
    validateSearch: (search) => orderSearchSchema.parse(search),
    loaderDeps: ({ search }) => ({ category: search.category }),
    component: OrderPage,
})

interface CartItem {
    product: Product
    variant?: InternalProductsProductOptionResponse
    quantity: number
}


function OrderPage() {
    const { t } = useTranslation()
    const queryClient = useQueryClient()
    const navigate = useNavigate({ from: Route.fullPath })
    const searchParams = Route.useSearch()

    const [searchTerm, setSearchTerm] = useState('')
    const { data: productsData } = useProductsListQuery({ limit: 100, search: searchTerm })
    const products = productsData?.products || []

    const selectedCategory = searchParams.category || "all"
    const [selectedOrderType, setSelectedOrderType] = useState<POSKasirInternalOrdersRepositoryOrderType>(POSKasirInternalOrdersRepositoryOrderType.OrderTypeDineIn)

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
    const [createdOrderId, setCreatedOrderId] = useState<string | null>(null)
    const [isLoadingDetailsId, setIsLoadingDetailsId] = useState<string | null>(null)

    const [variantSelectionOpen, setVariantSelectionOpen] = useState(false)
    const [productForVariantSelection, setProductForVariantSelection] = useState<Product | null>(null)

    const createOrderMutation = useCreateOrderMutation()

    const addToCart = async (product: Product) => {
        if ((product.stock || 0) <= 0) {
            toast.error(t('order.errors.stock_empty'))
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
            toast.error(t('order.errors.load_failed'))
        } finally {
            setIsLoadingDetailsId(null)
        }
    }

    const addCartItem = (product: Product, variant?: InternalProductsProductOptionResponse) => {
        setCart(prev => {
            const existing = prev.find(item =>
                item.product.id === product.id && item.variant?.id === variant?.id
            )

            const currentQty = existing ? existing.quantity : 0
            if (currentQty + 1 > (product.stock || 0)) {
                toast.error(t('order.errors.insufficient_stock'))
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
                        toast.error(t('order.errors.insufficient_stock'))
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
                quantity: item.quantity,
                options: item.variant && item.variant.id ? [{ product_option_id: item.variant.id }] : []
            })),
            type: selectedOrderType
        }

        toast.promise(createOrderMutation.mutateAsync(orderData), {
            loading: t('order.success.creating'),
            success: (data) => {
                setCreatedOrderId(data.id || null)
                setIsPaymentOpen(true)
                return t('order.success.created')
            },
            error: t('order.errors.create_failed')
        })
    }

    return (
        <div className="flex  h-[calc(100vh-4rem)] gap-2 ">
            {/* Left: Product Grid */}
            <div className="flex-1 flex flex-col gap-4 overflow-hidden   bg-background  min-h-0">
                <div className="">
                    <div className="relative">
                        <Search className="absolute left-2.5 top-3 h-6 w-6 text-muted-foreground" />
                        <Input
                            type="search"
                            placeholder={t('order.search_placeholder')}
                            className="pl-12 py-6"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                    </div>
                </div>

                <div className="px-4">
                    <Tabs value={selectedCategory} onValueChange={(v) => navigate({ search: (prev) => ({ ...prev, category: v }) })} className="w-full">
                        <TabsList className="w-full justify-start overflow-x-auto h-auto p-1 no-scrollbar">
                            <TabsTrigger value="all" className="rounded-full px-4">{t('order.all_categories')}</TabsTrigger>
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
                                            <span className="bg-background px-2 text-muted-foreground">{t('order.out_of_stock')}</span>
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
                                    {t('order.no_products')}
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
                    <h2 className="font-semibold">{t('order.current_order')}</h2>
                    <span className="ml-auto text-xs font-medium bg-primary/10 text-primary px-2 py-0.5 rounded-full">
                        {cart.length} {t('order.items')}
                    </span>
                </div>

                <div className="mx-4 mt-4">
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
                            <span className="text-muted-foreground">{t('order.subtotal')}</span>
                            <span>{formatRupiah(calculateTotal())}</span>
                        </div>

                        <div className="flex justify-between text-lg font-bold border-t pt-2 mt-2">
                            <span>{t('order.total')}</span>
                            <span className="text-primary">{formatRupiah(calculateTotal())}</span>
                        </div>
                    </div>
                    <Button className="w-full h-12 text-lg " size="lg" disabled={cart.length === 0} onClick={handleCheckout}>
                        {t('order.charge')} {formatRupiah(calculateTotal())}
                    </Button>
                </div>
            </div>

            {/* Payment Dialog Component */}
            <PaymentDialog
                open={isPaymentOpen}
                onOpenChange={setIsPaymentOpen}
                orderId={createdOrderId}
                onPaymentSuccess={() => {
                    setCart([])
                    setCreatedOrderId(null)
                    // setIsPaymentOpen(false) // Handled by component usually, but parent state needs to sync? PaymentDialog calls onOpenChange(false) automatically.
                }}
            />

            {/* Variant Selection Dialog */}
            <Dialog open={variantSelectionOpen} onOpenChange={setVariantSelectionOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>{t('order.variant_dialog.title')}</DialogTitle>
                        <DialogDescription>
                            {t('order.variant_dialog.desc')} <span className="font-semibold">{productForVariantSelection?.name}</span>
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
                                    <span className="text-xs text-muted-foreground">{t('order.variant_dialog.original')}</span>
                                )}
                            </div>
                            <div className="flex flex-col">
                                <span className="font-medium text-sm">{t('order.variant_dialog.original')}</span>
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
                                            <span className="text-xs">{t('order.variant_dialog.no_image')}</span>
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
                        <Button variant="outline" onClick={() => setVariantSelectionOpen(false)}>{t('order.payment_dialog.cancel')}</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    )
}
