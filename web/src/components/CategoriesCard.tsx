import React, { useEffect, useState } from 'react'
import { useSuspenseQuery } from '@tanstack/react-query'
import {
    categoriesListQueryOptions,
    type Category,
    useCreateCategoryMutation,
    useDeleteCategoryMutation,
    useUpdateCategoryMutation
} from '@/lib/api/query/categories'
import { POSKasirInternalDtoCreateCategoryRequest } from '@/lib/api/generated'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

import { Card, CardContent, CardDescription, CardHeader, CardTitle, } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow, } from '@/components/ui/table'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
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

import { Loader2, MoreHorizontal, Package, Pencil, Plus, Tag, Trash2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'

// --- MAIN COMPONENT ---
export function CategoriesCard() {
    const { t } = useTranslation()
    const { data: categories } = useSuspenseQuery(categoriesListQueryOptions())
    const categoriesList = Array.isArray(categories) ? categories : (categories as any)?.data || []

    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [selectedCategory, setSelectedCategory] = useState<Category | null>(null)

    const openCreateModal = () => {
        setSelectedCategory(null)
        setIsDialogOpen(true)
    }

    const openEditModal = (category: Category) => {
        setSelectedCategory(category)
        setIsDialogOpen(true)
    }

    return (
        <>
            <Card>
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <div>
                            <CardTitle className="flex items-center gap-2">
                                <Tag className="h-5 w-5" /> {t('settings.category.title')}
                            </CardTitle>
                            <CardDescription>
                                {t('settings.category.description')}
                            </CardDescription>
                        </div>
                        <Button onClick={openCreateModal} size="sm">
                            <Plus className="mr-2 h-4 w-4" /> {t('settings.category.add_button')}
                        </Button>
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="rounded-md border">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>{t('settings.category.table.name')}</TableHead>
                                    {/* Kolom Description Dihapus */}
                                    <TableHead className="w-[150px] text-right">{t('settings.category.table.created_at')}</TableHead>
                                    <TableHead className="w-[80px] text-right">{t('settings.category.table.actions')}</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {categoriesList.length === 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={3} className="h-24 text-center text-muted-foreground">
                                            {t('settings.category.table.empty')}
                                        </TableCell>
                                    </TableRow>
                                ) : (
                                    categoriesList.map((category: Category) => (
                                        <TableRow key={category.id}>
                                            <TableCell className="font-medium">
                                                <div className="flex items-center gap-2">
                                                    <Package className="h-4 w-4 text-muted-foreground" />
                                                    {category.name}
                                                </div>
                                            </TableCell>
                                            {/* Menampilkan tanggal pembuatan sebagai ganti deskripsi */}
                                            <TableCell className="text-right text-muted-foreground text-sm">
                                                {category.created_at ? new Date(category.created_at).toLocaleDateString() : '-'}
                                            </TableCell>
                                            <TableCell className="text-right">
                                                <CategoryActions
                                                    category={category}
                                                    onEdit={() => openEditModal(category)}
                                                />
                                            </TableCell>
                                        </TableRow>
                                    ))
                                )}
                            </TableBody>
                        </Table>
                    </div>
                </CardContent>
            </Card>

            <CategoryFormDialog
                open={isDialogOpen}
                onOpenChange={setIsDialogOpen}
                categoryToEdit={selectedCategory}
            />
        </>
    )
}

// --- SUB-COMPONENT: Actions ---
function CategoryActions({ category, onEdit }: { category: Category, onEdit: () => void }) {
    const { t } = useTranslation()
    const deleteMutation = useDeleteCategoryMutation()
    const [showDeleteDialog, setShowDeleteDialog] = useState(false)

    const handleDelete = (e: React.MouseEvent) => {
        e.preventDefault()
        if (category.id) {
            deleteMutation.mutate(category.id, {
                onSuccess: () => setShowDeleteDialog(false)
            })
        }
    }

    return (
        <>
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">Open menu</span>
                        <MoreHorizontal className="h-4 w-4" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuLabel>{t('settings.category.table.actions')}</DropdownMenuLabel>
                    <DropdownMenuItem onClick={onEdit}>
                        <Pencil className="mr-2 h-4 w-4" /> {t('settings.category.actions.edit')}
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onSelect={(e) => {
                            e.preventDefault()
                            setShowDeleteDialog(true)
                        }}
                        className="text-red-600 focus:text-red-600 cursor-pointer"
                    >
                        <Trash2 className="mr-2 h-4 w-4" /> {t('settings.category.actions.delete')}
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>

            <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>{t('settings.category.actions.delete_title')}</AlertDialogTitle>
                        <AlertDialogDescription>
                            {t('settings.category.actions.delete_confirm')} <span
                                className="font-semibold text-foreground">"{category.name}"</span>?
                            {t('settings.category.actions.delete_warning')}
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel disabled={deleteMutation.isPending}>{t('settings.category.actions.cancel')}</AlertDialogCancel>
                        <AlertDialogAction
                            onClick={handleDelete}
                            disabled={deleteMutation.isPending}
                            className="bg-red-600 hover:bg-red-700 focus:ring-red-600"
                        >
                            {deleteMutation.isPending ? (
                                <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> {t('settings.category.actions.deleting')}</>
                            ) : (
                                t('settings.category.actions.delete_button')
                            )}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}

// --- SUB-COMPONENT: Form Dialog ---
function CategoryFormDialog({ open, onOpenChange, categoryToEdit }: {
    open: boolean,
    onOpenChange: (open: boolean) => void,
    categoryToEdit: Category | null
}) {
    const { t } = useTranslation()
    const createMutation = useCreateCategoryMutation()
    const updateMutation = useUpdateCategoryMutation()

    const [formData, setFormData] = useState({
        name: ''
    })

    useEffect(() => {
        if (open) {

            if (categoryToEdit) {
                setFormData({
                    name: categoryToEdit.name || ''
                })
            } else {
                setFormData({
                    name: ''
                })
            }
        }
    }, [open, categoryToEdit])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const payload: POSKasirInternalDtoCreateCategoryRequest = {
            name: formData.name
        }

        try {
            if (categoryToEdit && categoryToEdit.id) {
                await updateMutation.mutateAsync({ id: categoryToEdit.id, body: payload })
            } else {
                await createMutation.mutateAsync(payload)
            }
            onOpenChange(false)
            setFormData({ name: '' })
        } catch (error) {
            console.error(error)
        }
    }

    const isSubmitting = createMutation.isPending || updateMutation.isPending

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>{categoryToEdit ? t('settings.category.form.title_edit') : t('settings.category.form.title_add')}</DialogTitle>
                    <DialogDescription>
                        {categoryToEdit ? t('settings.category.form.desc_edit') : t('settings.category.form.desc_add')}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        <div className="grid gap-2">
                            <Label htmlFor="name">{t('settings.category.form.name_label')}</Label>
                            <Input
                                id="name"
                                value={formData.name}
                                onChange={(e) => setFormData({ name: e.target.value })}
                                placeholder={t('settings.category.form.name_placeholder')}
                                required
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? (
                                <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> {t('settings.category.form.saving')}</>
                            ) : (
                                t('settings.category.form.save')
                            )}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}