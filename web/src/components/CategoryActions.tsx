import { Category, useDeleteCategoryMutation } from "@/lib/api/query/categories"
import { AlertDialog, AlertDialogContent, AlertDialogTitle, AlertDialogDescription, AlertDialogCancel, AlertDialogAction } from "@radix-ui/react-alert-dialog"
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuLabel, DropdownMenuItem } from "@radix-ui/react-dropdown-menu"
import { MoreHorizontal, Pencil, Trash2, Loader2 } from "lucide-react"
import { useState } from "react"

import { useTranslation } from "react-i18next"
import { AlertDialogHeader, AlertDialogFooter } from "./ui/alert-dialog"
import { Button } from "./ui/button"

export function CategoryActions({ category, onEdit }: { category: Category, onEdit: () => void }) {
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
                        className="text-destructive focus:text-destructive cursor-pointer"
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
                            className="bg-destructive hover:bg-destructive/90 focus:ring-destructive"
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
