import { Product, useDeleteProductMutation, useUpdateProductMutation } from "@/lib/api/query/products.ts";
import React, { useState } from "react";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Loader2, MoreHorizontal, Pencil, Trash2 } from "lucide-react";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle
} from "@/components/ui/alert-dialog.tsx";
import { useTranslation } from 'react-i18next';

export function ProductActions({ product, onEdit }: { product: Product, onEdit?: () => void }) {
    const { t } = useTranslation();
    const deleteMutation = useDeleteProductMutation()
    const updateMutation = useUpdateProductMutation()

    const [showDeleteDialog, setShowDeleteDialog] = useState(false)

    const canEdit = updateMutation.isAllowed
    const canDelete = deleteMutation.isAllowed

    if (!canEdit && !canDelete) return null

    const handleDelete = (e: React.MouseEvent) => {
        e.preventDefault()
        if (product.id) {
            deleteMutation.mutate(product.id, {
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
                    <DropdownMenuLabel>{t('products.table.actions')}</DropdownMenuLabel>

                    {canEdit && (
                        <DropdownMenuItem onClick={onEdit}>
                            <Pencil className="mr-2 h-4 w-4" /> {t('products.actions.edit')}
                        </DropdownMenuItem>
                    )}

                    {canEdit && canDelete && <DropdownMenuSeparator />}

                    {canDelete && (
                        <DropdownMenuItem
                            onSelect={(e) => {
                                e.preventDefault()
                                setShowDeleteDialog(true)
                            }}
                            className="text-destructive focus:text-destructive cursor-pointer"
                        >
                            <Trash2 className="mr-2 h-4 w-4" /> {t('products.actions.delete')}
                        </DropdownMenuItem>
                    )}
                </DropdownMenuContent>
            </DropdownMenu>

            <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>{t('products.actions.delete_title')}</AlertDialogTitle>
                        <AlertDialogDescription>
                            {t('products.actions.delete_desc')} <span
                                className="font-semibold text-foreground">"{product.name}"</span>?
                            {t('products.actions.delete_desc_2')}
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel disabled={deleteMutation.isPending}>{t('products.actions.cancel')}</AlertDialogCancel>
                        <AlertDialogAction
                            onClick={handleDelete}
                            disabled={deleteMutation.isPending}
                            className="bg-destructive hover:bg-destructive/90 focus:ring-destructive"
                        >
                            {deleteMutation.isPending ? (
                                <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> {t('products.actions.deleting')}</>
                            ) : (
                                t('products.actions.delete_confirm')
                            )}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}