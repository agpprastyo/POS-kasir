import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useSuspenseQuery, useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import React, { useState, useEffect } from 'react'
import { useDebounce } from '@/hooks/use-debounce'
import { type Product, productsListQueryOptions, trashProductsListQueryOptions, useRestoreProductMutation } from '@/lib/api/query/products'
import { categoriesListQueryOptions } from '@/lib/api/query/categories'
import { queryClient } from '@/lib/queryClient'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow, } from '@/components/ui/table'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue, } from "@/components/ui/select"
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Filter, Package, Plus, Search, LayoutGrid, LayoutList } from 'lucide-react'
import { formatRupiah } from "@/lib/utils.ts";
import { NewPagination } from "@/components/pagination.tsx";
import { ProductFormDialog } from "@/components/ProductFormDialog.tsx";
import { ProductActions } from "@/components/ProductActions.tsx";
import { ProductCard } from '@/components/ProductCard'

import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { useTranslation } from 'react-i18next';


const productsSearchSchema = z.object({
    page: z.number().catch(1),
    limit: z.number().min(1).catch(10),
    search: z.string().optional(),
    category: z.number().optional(),
    tab: z.enum(['active', 'trash']).catch('active'),
})


export const Route = createFileRoute('/$locale/_dashboard/product')({
    validateSearch: (search) => productsSearchSchema.parse(search),

    loaderDeps: ({ search }) => ({
        page: search.page,
        limit: search.limit,
        search: search.search,
        category: search.category,
        tab: search.tab,
    }),

    loader: ({ deps }) => {
        return queryClient.ensureQueryData(productsListQueryOptions({
            limit: deps.limit,
            page: deps.page,
            search: deps.search,
            category: deps.category
        }))
    },

    component: ProductsPage,
})

function ProductsPage() {
    const { t } = useTranslation();
    const navigate = useNavigate({ from: Route.fullPath })
    const searchParams = Route.useSearch()
    const productsQuery = useQuery(productsListQueryOptions({
        limit: searchParams.limit,
        page: searchParams.page,
        search: searchParams.search,
        category: searchParams.category,
    }))


    const { data: categoriesData } = useSuspenseQuery(categoriesListQueryOptions())
    const categories = Array.isArray(categoriesData) ? categoriesData : (categoriesData as any)?.data || []

    const products = productsQuery.data?.products || []
    const pagination = productsQuery.data?.pagination || {}

    const trashProductsQuery = useQuery(trashProductsListQueryOptions({
        limit: searchParams.limit,
        page: searchParams.page,
        search: searchParams.search,
        category: searchParams.category,
    }))
    const trashProducts = trashProductsQuery.data?.products || []
    const trashPagination = trashProductsQuery.data?.pagination || {}

    const restoreMutation = useRestoreProductMutation()

    const currentTab = searchParams.tab || 'active'

    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
    const [viewMode, setViewMode] = useState<'list' | 'grid'>('grid')

    const [searchTerm, setSearchTerm] = useState(searchParams.search || '')
    const debouncedSearch = useDebounce(searchTerm, 500)

    useEffect(() => {
        if (debouncedSearch !== (searchParams.search || '')) {
            navigate({
                to: '.',
                search: (prev) => ({ ...prev, search: debouncedSearch || undefined, page: 1 }),
                replace: true
            })
        }
    }, [debouncedSearch, navigate, searchParams.search])


    const handleFilterCategory = (value: string) => {
        navigate({
            to: '.',
            search: (prev) => ({
                ...prev,
                category: value === 'all' ? undefined : Number(value),
                page: 1
            })
        })
    }

    const handlePageChange = (newPage: number) => {
        navigate({ to: '.', search: (prev) => ({ ...prev, page: newPage }) })
    }

    const openCreateModal = () => {
        setSelectedProduct(null)
        setIsDialogOpen(true)
    }

    const openEditModal = (product: Product) => {
        setSelectedProduct(product)
        setIsDialogOpen(true)
    }

    return (
        <div className="flex flex-col gap-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">{t('products.title')}</h1>
                    <p className="text-muted-foreground">{t('products.description')}</p>
                </div>
                <div className="flex items-center gap-2">
                    <Button onClick={openCreateModal}>
                        <Plus className="mr-2 h-4 w-4" /> {t('products.add_button')}
                    </Button>
                </div>
            </div>

            <Tabs value={currentTab} onValueChange={(v) => {
                navigate({ search: (prev) => ({ ...prev, tab: v as 'active' | 'trash', page: 1 }) })
            }} className="w-full">
                <div className="flex items-center justify-between mb-4">
                    <TabsList>
                        <TabsTrigger value="active">{t('products.tabs.active')}</TabsTrigger>
                        <TabsTrigger value="trash">{t('products.tabs.trash', { count: trashPagination.total_data || 0 })}</TabsTrigger>
                    </TabsList>

                    <Tabs value={viewMode} onValueChange={(v) => setViewMode(v as 'list' | 'grid')} className="w-[80px]">
                        <TabsList className="grid w-full grid-cols-2">
                            <TabsTrigger value="list" className='px-2' title={t('products.tabs.list')}><LayoutList className="h-4 w-4" /></TabsTrigger>
                            <TabsTrigger value="grid" className='px-2' title={t('products.tabs.grid')}><LayoutGrid className="h-4 w-4" /></TabsTrigger>
                        </TabsList>
                    </Tabs>
                </div>

                {/* Filters */}
                <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
                    <div className="relative flex-1 md:max-w-sm">
                        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                        <Input
                            type="search"
                            placeholder={t('products.filters.search_placeholder')}
                            className="pl-8"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                    </div>

                    <Select
                        value={searchParams.category ? String(searchParams.category) : 'all'}
                        onValueChange={handleFilterCategory}
                    >
                        <SelectTrigger className="w-[180px]">
                            <Filter className="mr-2 h-4 w-4" />
                            <SelectValue placeholder={t('products.filters.category')} />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">{t('products.filters.category_all')}</SelectItem>
                            {categories.map((cat: any) => (
                                <SelectItem key={cat.id} value={String(cat.id)}>
                                    {cat.name}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>

                <TabsContent value="active" className="space-y-4">
                    {/* Content */}
                    {viewMode === 'list' ? (
                        <div className="rounded-md border bg-card">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead className="w-[80px] hidden sm:table-cell">{t('products.table.image')}</TableHead>
                                        <TableHead>{t('products.table.name')}</TableHead>
                                        <TableHead className="hidden md:table-cell">{t('products.table.category')}</TableHead>
                                        <TableHead>{t('products.table.price')}</TableHead>
                                        <TableHead>{t('products.table.stock')}</TableHead>
                                        <TableHead className="text-right">{t('products.table.actions')}</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {products.length === 0 ? (
                                        <TableRow>
                                            <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                                                {t('products.table.empty')}
                                            </TableCell>
                                        </TableRow>
                                    ) : (
                                        products.map((product: Product) => (
                                            <TableRow key={product.id}>
                                                <TableCell className="hidden sm:table-cell">
                                                    <Avatar className="h-10 w-10 rounded-md border">
                                                        <AvatarImage src={product.image_url} alt={product.name}
                                                            className="object-cover" />
                                                        <AvatarFallback className="rounded-md">
                                                            <Package className="h-5 w-5 text-muted-foreground" />
                                                        </AvatarFallback>
                                                    </Avatar>
                                                </TableCell>
                                                <TableCell className="font-medium">
                                                    {product.name}
                                                </TableCell>
                                                <TableCell className="hidden md:table-cell">
                                                    <Badge variant="outline">{product.category_name || '-'}</Badge>
                                                </TableCell>
                                                <TableCell>
                                                    {formatRupiah(product.price || 0)}
                                                </TableCell>
                                                <TableCell>
                                                    <Badge
                                                        variant={product.stock && product.stock > 0 ? 'secondary' : 'destructive'}>
                                                        {product.stock}
                                                    </Badge>
                                                </TableCell>
                                                <TableCell className="text-right">
                                                    <ProductActions product={product} onEdit={() => openEditModal(product)} />
                                                </TableCell>
                                            </TableRow>
                                        ))
                                    )}
                                </TableBody>
                            </Table>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
                            {products.length === 0 ? (
                                <div className="col-span-full h-24 flex items-center justify-center text-muted-foreground border rounded-md border-dashed">
                                    {t('products.table.empty')}
                                </div>
                            ) : (
                                products.map((product: Product) => (
                                    <ProductCard
                                        key={product.id}
                                        product={product}
                                        onEdit={openEditModal}
                                    />
                                ))
                            )}
                        </div>
                    )}

                    {products.length > 0 && (
                        <NewPagination pagination={pagination} onClick={() => handlePageChange((pagination.current_page || 1) - 1)}
                            onClick1={() => handlePageChange((pagination.current_page || 1) + 1)} />
                    )}
                </TabsContent>

                <TabsContent value="trash" className="space-y-4">
                    {/* Content Trash */}
                    {viewMode === 'list' ? (
                        <div className="rounded-md border bg-card">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead className="w-[80px] hidden sm:table-cell">{t('products.table.image')}</TableHead>
                                        <TableHead>{t('products.table.name')}</TableHead>
                                        <TableHead className="hidden md:table-cell">{t('products.table.category')}</TableHead>
                                        <TableHead>{t('products.table.price')}</TableHead>
                                        <TableHead>{t('products.table.stock')}</TableHead>
                                        <TableHead className="text-right">{t('products.table.actions')}</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {trashProducts.length === 0 ? (
                                        <TableRow>
                                            <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                                                {t('products.table.empty_trash')}
                                            </TableCell>
                                        </TableRow>
                                    ) : (
                                        trashProducts.map((product: Product) => (
                                            <TableRow key={product.id} className="opacity-70 bg-muted/30">
                                                <TableCell className="hidden sm:table-cell">
                                                    <Avatar className="h-10 w-10 rounded-md border grayscale">
                                                        <AvatarImage src={product.image_url} alt={product.name}
                                                            className="object-cover" />
                                                        <AvatarFallback className="rounded-md">
                                                            <Package className="h-5 w-5 text-muted-foreground" />
                                                        </AvatarFallback>
                                                    </Avatar>
                                                </TableCell>
                                                <TableCell className="font-medium text-muted-foreground">
                                                    {product.name}
                                                </TableCell>
                                                <TableCell className="hidden md:table-cell">
                                                    <Badge variant="outline" className="opacity-50">{product.category_name || '-'}</Badge>
                                                </TableCell>
                                                <TableCell className="text-muted-foreground">
                                                    {formatRupiah(product.price || 0)}
                                                </TableCell>
                                                <TableCell>
                                                    <span className="text-muted-foreground">{product.stock}</span>
                                                </TableCell>
                                                <TableCell className="text-right">
                                                    <Button
                                                        size="sm"
                                                        variant="outline"
                                                        onClick={() => product.id && restoreMutation.mutate(product.id)}
                                                    >
                                                        {t('products.actions.restore')}
                                                    </Button>
                                                </TableCell>
                                            </TableRow>
                                        ))
                                    )}
                                </TableBody>
                            </Table>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
                            {trashProducts.length === 0 ? (
                                <div className="col-span-full h-24 flex items-center justify-center text-muted-foreground border rounded-md border-dashed">
                                    {t('products.table.empty_trash')}
                                </div>
                            ) : (
                                trashProducts.map((product: Product) => (
                                    <ProductCard
                                        key={product.id}
                                        product={product}
                                        onRestore={(p) => p.id && restoreMutation.mutate(p.id)}
                                    />
                                ))
                            )}
                        </div>
                    )}

                    {trashProducts.length > 0 && (
                        <NewPagination pagination={trashPagination} onClick={() => handlePageChange((trashPagination.current_page || 1) - 1)}
                            onClick1={() => handlePageChange((trashPagination.current_page || 1) + 1)} />
                    )}
                </TabsContent>
            </Tabs>

            <ProductFormDialog
                open={isDialogOpen}
                onOpenChange={setIsDialogOpen}
                productToEdit={selectedProduct}
                categories={categories}
            />
        </div>
    )
}



