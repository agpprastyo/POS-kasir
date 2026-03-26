import { useState } from 'react'
import { useSuspenseQuery } from '@tanstack/react-query'
import {
    categoriesListQueryOptions,
    type Category
} from '@/lib/api/query/categories'
import { Button } from '@/components/ui/button'

import { Card, CardContent, CardDescription, CardHeader, CardTitle, } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow, } from '@/components/ui/table'

import { Package, Plus, Tag } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { CategoryFormDialog } from './CategoryFormDialog'
import { CategoryActions } from './CategoryActions'

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


