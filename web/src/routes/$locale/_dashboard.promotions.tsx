import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { format } from 'date-fns'

import { Button } from '@/components/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { MoreHorizontal } from "lucide-react"

import {
    promotionsListQueryOptions,
    useDeletePromotionMutation,
    Promotion
} from '@/lib/api/query/promotions'
import { queryClient } from '@/lib/queryClient'
import { PromotionFormDialog } from '@/components/PromotionFormDialog'
import { NewPagination } from "@/components/pagination"
import { formatRupiah } from "@/lib/utils"
import { POSKasirInternalRepositoryDiscountType } from '@/lib/api/generated'

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useRestorePromotionMutation } from '@/lib/api/query/promotions'
import { RotateCcw } from 'lucide-react'

const promotionsSearchSchema = z.object({
    page: z.number().catch(1),
    limit: z.number().min(1).catch(10),
})

export const Route = createFileRoute('/$locale/_dashboard/promotions')({
    validateSearch: (search) => promotionsSearchSchema.parse(search),
    loaderDeps: ({ search }) => ({
        page: search.page,
        limit: search.limit,
    }),
    loader: ({ deps }) => {
        return queryClient.ensureQueryData(promotionsListQueryOptions({
            limit: deps.limit,
            page: deps.page,
        }))
    },
    component: PromotionsPage,
})



function PromotionsPage() {
    const { t } = useTranslation()
    const navigate = useNavigate({ from: Route.fullPath })
    const searchParams = Route.useSearch()

    const [activeTab, setActiveTab] = useState("active")

    const promotionsQuery = useQuery(promotionsListQueryOptions({
        limit: searchParams.limit,
        page: searchParams.page,
        trash: activeTab === "trash"
    }))

    const promotions = promotionsQuery.data?.promotions || []
    const pagination = promotionsQuery.data?.pagination || {}

    const deleteMutation = useDeletePromotionMutation()
    const restoreMutation = useRestorePromotionMutation()

    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [selectedPromotion, setSelectedPromotion] = useState<Promotion | null>(null)

    const handlePageChange = (newPage: number) => {
        navigate({ search: (prev) => ({ ...prev, page: newPage }) })
    }

    const openCreateModal = () => {
        setSelectedPromotion(null)
        setIsDialogOpen(true)
    }

    const openEditModal = (promo: Promotion) => {
        setSelectedPromotion(promo)
        setIsDialogOpen(true)
    }

    const handleDelete = (id: string) => {
        if (confirm(t('common.confirm') + '?')) {
            deleteMutation.mutate(id)
        }
    }

    const handleRestore = (id: string) => {
        if (confirm(t('common.confirm') + '?')) {
            restoreMutation.mutate(id)
        }
    }

    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">{t('promotions.title')}</h1>
                    <p className="text-muted-foreground">{t('promotions.description')}</p>
                </div>
                <div className="flex items-center gap-2">
                    <Button onClick={openCreateModal}>
                        <Plus className="mr-2 h-4 w-4" /> {t('promotions.add_button')}
                    </Button>
                </div>
            </div>

            <Tabs defaultValue="active" onValueChange={setActiveTab} className="w-full">
                <TabsList>
                    <TabsTrigger value="active">Active</TabsTrigger>
                    <TabsTrigger value="trash">Trash</TabsTrigger>
                </TabsList>

                <TabsContent value="active" className="mt-4">
                    <PromotionsTable
                        promotions={promotions}
                        t={t}
                        onEdit={openEditModal}
                        onDelete={handleDelete}
                        isTrash={false}
                    />
                </TabsContent>

                <TabsContent value="trash" className="mt-4">
                    <PromotionsTable
                        promotions={promotions}
                        t={t}
                        onEdit={openEditModal}
                        onDelete={handleDelete}
                        onRestore={handleRestore}
                        isTrash={true}
                    />
                </TabsContent>
            </Tabs>

            {promotions.length > 0 && pagination && (
                <NewPagination
                    pagination={pagination}
                    onClick={() => handlePageChange(((pagination as any).current_page || 1) - 1)}
                    onClick1={() => handlePageChange(((pagination as any).current_page || 1) + 1)}
                />
            )}

            <PromotionFormDialog
                open={isDialogOpen}
                onOpenChange={setIsDialogOpen}
                promotionToEdit={selectedPromotion}
            />
        </div>
    )
}

function PromotionsTable({ promotions, t, onEdit, onDelete, onRestore, isTrash }: any) {
    return (
        <div className="rounded-md border bg-card">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>{t('promotions.table.name')}</TableHead>
                        <TableHead>{t('promotions.table.scope')}</TableHead>
                        <TableHead>{t('promotions.table.discount')}</TableHead>
                        <TableHead>{t('promotions.table.period')}</TableHead>
                        <TableHead>{t('promotions.table.status')}</TableHead>
                        <TableHead className="text-right">{t('promotions.table.actions')}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {promotions.length === 0 ? (
                        <TableRow>
                            <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                                {t('promotions.table.empty')}
                            </TableCell>
                        </TableRow>
                    ) : (
                        promotions.map((promo: Promotion) => (
                            <TableRow key={promo.id}>
                                <TableCell className="font-medium">
                                    <div className="flex flex-col">
                                        <span>{promo.name}</span>
                                        {promo.description && <span className="text-xs text-muted-foreground">{promo.description}</span>}
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant="outline">{promo.scope}</Badge>
                                </TableCell>
                                <TableCell>
                                    <div className="flex flex-col">
                                        <span className="font-bold">
                                            {promo.discount_type === POSKasirInternalRepositoryDiscountType.DiscountTypePercentage
                                                ? `${promo.discount_value}%`
                                                : formatRupiah(promo.discount_value)}
                                        </span>
                                        {promo.max_discount_amount && promo.max_discount_amount > 0 && (
                                            <span className="text-xs text-muted-foreground">Max: {formatRupiah(promo.max_discount_amount)}</span>
                                        )}
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <div className="text-sm">
                                        {format(new Date(promo.start_date), 'dd MMM yyyy')} - {format(new Date(promo.end_date), 'dd MMM yyyy')}
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant={promo.is_active ? 'default' : 'secondary'}>
                                        {promo.is_active ? 'Active' : 'Inactive'}
                                    </Badge>
                                </TableCell>
                                <TableCell className="text-right">
                                    <DropdownMenu>
                                        <DropdownMenuTrigger asChild>
                                            <Button variant="ghost" className="h-8 w-8 p-0">
                                                <span className="sr-only">Open menu</span>
                                                <MoreHorizontal className="h-4 w-4" />
                                            </Button>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent align="end">
                                            <DropdownMenuLabel>{t('common.actions')}</DropdownMenuLabel>
                                            {!isTrash && (
                                                <>
                                                    <DropdownMenuItem onClick={() => onEdit(promo)}>
                                                        <Pencil className="mr-2 h-4 w-4" /> {t('common.edit')}
                                                    </DropdownMenuItem>
                                                    <DropdownMenuSeparator />
                                                    <DropdownMenuItem onClick={() => onDelete(promo.id)} className="text-destructive">
                                                        <Trash2 className="mr-2 h-4 w-4" /> {t('common.delete')}
                                                    </DropdownMenuItem>
                                                </>
                                            )}
                                            {isTrash && (
                                                <DropdownMenuItem onClick={() => onRestore(promo.id)}>
                                                    <RotateCcw className="mr-2 h-4 w-4" /> Restore
                                                </DropdownMenuItem>
                                            )}
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </TableCell>
                            </TableRow>
                        ))
                    )}
                </TableBody>
            </Table>
        </div>
    )
}
