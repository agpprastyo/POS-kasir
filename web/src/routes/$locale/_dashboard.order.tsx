import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { useProductsListQuery, Product, productDetailQueryOptions } from '@/lib/api/query/products'
import { ShoppingCart } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { PaymentDialog } from "@/components/payment/PaymentDialog"
import { useCreateOrderMutation } from '@/lib/api/query/orders'
import { toast } from 'sonner'
import { InternalProductsProductOptionResponse, POSKasirInternalOrdersRepositoryOrderType } from '@/lib/api/generated'
import { useCustomersListQuery } from '@/lib/api/query/customers'
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'
import { CartContent, CartItem } from '@/components/orders/CartContent'
import { VariantSelectionDialog } from '@/components/orders/VariantSelectionDialog'
import { ProductSearch } from '@/components/orders/ProductSearch'
import { ProductGrid } from '@/components/orders/ProductGrid'

const orderSearchSchema = z.object({
    category: z.string().optional().catch('all'),
})

export const Route = createFileRoute('/$locale/_dashboard/order')({
    validateSearch: (search) => orderSearchSchema.parse(search),
    loaderDeps: ({ search }) => ({ category: search.category }),
    component: OrderPage,
})

function OrderPage() {
    const { t } = useTranslation()
    const queryClient = useQueryClient()
    const navigate = useNavigate({ from: Route.fullPath })
    const searchParams = Route.useSearch()

    const [searchTerm, setSearchTerm] = useState('')
    const { data: productsData } = useProductsListQuery({ limit: 100, search: searchTerm })
    const products = productsData?.products || []

    const { data: customersData } = useCustomersListQuery({ limit: 100 })
    const customers = customersData?.customers || []

    const [selectedCustomerId, setSelectedCustomerId] = useState<string | null>(null)

    const selectedCategory = searchParams.category || "all"
    const [selectedOrderType, setSelectedOrderType] = useState<POSKasirInternalOrdersRepositoryOrderType>(POSKasirInternalOrdersRepositoryOrderType.OrderTypeDineIn)

    const categories = useMemo(() => {
        const unique = new Map<string, string>()
        products.forEach(p => {
            p.categories?.forEach(cat => {
                if (cat.id && cat.name) {
                    unique.set(String(cat.id), cat.name)
                }
            })
        })
        return Array.from(unique.entries())
            .map(([id, name]) => ({ id, name }))
            .sort((a, b) => a.name.localeCompare(b.name))
    }, [products])

    const filteredProducts = useMemo(() => {
        if (selectedCategory === "all") return products
        return products.filter(p => p.categories?.some(cat => String(cat.id) === selectedCategory))
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
    const canCreateOrder = createOrderMutation.isAllowed

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
            type: selectedOrderType,
            customer_id: selectedCustomerId || undefined
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
        <div className="flex h-[calc(100vh-4rem)] gap-2 relative">
            {/* Left: Product Grid */}
            <div className="flex-1 flex flex-col gap-4 overflow-hidden bg-background min-h-0">
                <ProductSearch 
                    searchTerm={searchTerm}
                    onSearchChange={setSearchTerm}
                    selectedCategory={selectedCategory}
                    onCategoryChange={(v) => navigate({ search: (prev) => ({ ...prev, category: v }) })}
                    categories={categories}
                    t={t}
                />

                <div className="flex-1 min-h-0">
                    <ProductGrid 
                        inStockProducts={inStockProducts}
                        outOfStockProducts={outOfStockProducts}
                        isLoadingDetailsId={isLoadingDetailsId}
                        onAddToCart={addToCart}
                        t={t}
                    />
                </div>
            </div>

            {/* Cart Mobile Toggle */}
            <div className="md:hidden fixed bottom-6 right-6 z-50">
                <Sheet>
                    <SheetTrigger asChild>
                        <Button size="lg" className="rounded-full h-14 w-14 shadow-lg p-0 relative">
                            <ShoppingCart className="h-6 w-6" />
                            {cart.length > 0 && (
                                <span className="absolute -top-1 -right-1 bg-destructive text-destructive-foreground text-xs font-bold h-5 w-5 rounded-full flex items-center justify-center border-2 border-background">
                                    {cart.length}
                                </span>
                            )}
                        </Button>
                    </SheetTrigger>
                    <SheetContent side="right" className="p-0 w-full sm:max-w-md border-l">
                        <CartContent 
                            cart={cart} 
                            t={t} 
                            customers={customers} 
                            selectedCustomerId={selectedCustomerId} 
                            setSelectedCustomerId={setSelectedCustomerId}
                            selectedOrderType={selectedOrderType}
                            setSelectedOrderType={setSelectedOrderType}
                            updateQuantity={updateQuantity}
                            removeFromCart={removeFromCart}
                            calculateTotal={calculateTotal}
                            handleCheckout={handleCheckout}
                            canCheckout={canCreateOrder}
                        />
                    </SheetContent>
                </Sheet>
            </div>

            {/* Right: Cart Sidebar (Desktop) */}
            <div className="hidden md:flex w-[350px] flex-col rounded-xl border bg-card overflow-hidden min-h-0">
                <CartContent 
                    cart={cart} 
                    t={t} 
                    customers={customers} 
                    selectedCustomerId={selectedCustomerId} 
                    setSelectedCustomerId={setSelectedCustomerId}
                    selectedOrderType={selectedOrderType}
                    setSelectedOrderType={setSelectedOrderType}
                    updateQuantity={updateQuantity}
                    removeFromCart={removeFromCart}
                    calculateTotal={calculateTotal}
                    handleCheckout={handleCheckout}
                    canCheckout={canCreateOrder}
                />
            </div>

            {/* Payment Dialog Component */}
            <PaymentDialog
                open={isPaymentOpen}
                onOpenChange={setIsPaymentOpen}
                orderId={createdOrderId}
                onPaymentSuccess={() => {
                    setCart([])
                    setCreatedOrderId(null)
                    setSelectedCustomerId(null)
                }}
            />

            {/* Variant Selection Dialog */}
            <VariantSelectionDialog 
                open={variantSelectionOpen}
                onOpenChange={setVariantSelectionOpen}
                product={productForVariantSelection}
                onSelect={addCartItem}
                t={t}
            />
        </div>
    )
}
