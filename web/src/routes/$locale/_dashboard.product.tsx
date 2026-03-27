import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useSuspenseQuery, useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { useState, useEffect } from 'react'
import { useDebounce } from '@/hooks/use-debounce'
import { type Product, productsListQueryOptions, useTrashProductsListQuery, useRestoreProductMutation, useCreateProductMutation, useUpdateProductMutation, useDeleteProductMutation } from '@/lib/api/query/products'
import { categoriesListQueryOptions } from '@/lib/api/query/categories'
import { Tabs, TabsContent } from '@/components/ui/tabs'
import { NewPagination } from "@/components/pagination.tsx";
import { useTranslation } from 'react-i18next';
import { ProductHeader } from '@/components/products/ProductHeader'
import { ProductFilters } from '@/components/products/ProductFilters'
import { ProductTable } from '@/components/products/ProductTable'
import { ProductGridView } from '@/components/products/ProductGridView'
import { ProductFormDialog } from '@/components/products/ProductFormDialog'

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

    loader: ({ context: { queryClient }, deps }) => {
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
    const createMutation = useCreateProductMutation()

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

    const trashProductsQuery = useTrashProductsListQuery({
        limit: searchParams.limit,
        page: searchParams.page,
        search: searchParams.search,
        category: searchParams.category,
    })
    const canReadTrash = trashProductsQuery.isAllowed
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

    const canCreate = createMutation.isAllowed
    const canRestore = restoreMutation.isAllowed

    const { isAllowed: canEdit } = useUpdateProductMutation()
    const { isAllowed: canDelete } = useDeleteProductMutation()
    const hasActions = canEdit || canDelete

    const onRestoreHandler = canRestore ? (p: Product) => p.id && restoreMutation.mutate(p.id) : undefined

    return (
        <div className="flex flex-col gap-6">
            <ProductHeader 
                canCreate={canCreate}
                openCreateModal={openCreateModal}
                t={t}
            />

            <Tabs value={currentTab} onValueChange={(v) => {
                navigate({ search: (prev) => ({ ...prev, tab: v as 'active' | 'trash', page: 1 }) })
            }} className="w-full">
                <ProductFilters 
                trashCount={trashPagination.total_data || 0}
                viewMode={viewMode}
                onViewModeChange={setViewMode}
                searchTerm={searchTerm}
                onSearchChange={setSearchTerm}
                category={searchParams.category}
                categories={categories}
                onCategoryChange={handleFilterCategory}
                t={t}
                canReadTrash={canReadTrash}
            />

            <TabsContent value="active" className="space-y-4">
                {viewMode === 'list' ? (
                    <ProductTable 
                        products={products}
                        emptyMessage={t('products.table.empty')}
                        onEdit={openEditModal}
                        t={t}
                        hasActions={hasActions}
                    />
                ) : (
                    <ProductGridView 
                        products={products}
                        emptyMessage={t('products.table.empty')}
                        onEdit={openEditModal}
                        t={t}
                        hasActions={hasActions}
                    />
                )}

                {products.length > 0 && (
                    <NewPagination pagination={pagination} onClickPrev={() => handlePageChange((pagination.current_page || 1) - 1)}
                        onClickNext={() => handlePageChange((pagination.current_page || 1) + 1)} />
                )}
            </TabsContent>

            {canReadTrash && (
                <TabsContent value="trash" className="space-y-4">
                    {viewMode === 'list' ? (
                    <ProductTable 
                        products={trashProducts}
                        emptyMessage={t('products.table.empty_trash')}
                        isTrash={true}
                        canRestore={canRestore}
                        onRestore={onRestoreHandler}
                        t={t}
                    />
                ) : (
                    <ProductGridView 
                        products={trashProducts}
                        emptyMessage={t('products.table.empty_trash')}
                        onRestore={onRestoreHandler}
                        t={t}
                    />
                )}

                    {trashProducts.length > 0 && (
                        <NewPagination pagination={trashPagination} onClickPrev={() => handlePageChange((trashPagination.current_page || 1) - 1)}
                            onClickNext={() => handlePageChange((trashPagination.current_page || 1) + 1)} />
                    )}
                </TabsContent>
            )}
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



