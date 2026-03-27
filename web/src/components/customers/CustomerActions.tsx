import { InternalCustomersCustomerResponse } from "@/lib/api/generated";
import { useTranslation } from "react-i18next";
import { useDeleteCustomerMutation } from "@/lib/api/query/customers.ts";
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

export function CustomerActions({ customer, onEdit, canEdit = true, canDelete = true }: { customer: InternalCustomersCustomerResponse, onEdit: () => void, canEdit?: boolean, canDelete?: boolean }) {
    const { t } = useTranslation()
    const deleteMutation = useDeleteCustomerMutation()

    const [showDeleteDialog, setShowDeleteDialog] = useState(false)

    const handleDelete = async (e: React.MouseEvent) => {
        e.preventDefault()

        deleteMutation.mutate(customer.id!)
        setShowDeleteDialog(false)
    }

    return (
        <>
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">{t('customers.actions.open_menu', 'Open menu')}</span>
                        <MoreHorizontal className="h-4 w-4" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuLabel>{t('customers.table.actions', 'Actions')}</DropdownMenuLabel>
                    {canEdit && (
                        <DropdownMenuItem onClick={onEdit}>
                            <Pencil className="mr-2 h-4 w-4" /> {t('customers.actions.edit', 'Edit')}
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
                            <Trash2 className="mr-2 h-4 w-4" /> {t('customers.actions.delete', 'Delete')}
                        </DropdownMenuItem>
                    )}
                </DropdownMenuContent>
            </DropdownMenu>

            <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>{t('customers.actions.delete_title', 'Delete Customer')}</AlertDialogTitle>
                        <AlertDialogDescription>
                            {t('customers.actions.delete_desc', 'Are you sure you want to delete customer')}:
                            <span className="font-bold text-foreground"> "{customer.name}"</span>?
                            {t('customers.actions.delete_desc_2', " This action cannot be undone.")}
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel
                            disabled={deleteMutation.isPending}>{t('customers.actions.cancel', 'Cancel')}</AlertDialogCancel>

                        <AlertDialogAction
                            onClick={handleDelete}
                            disabled={deleteMutation.isPending}
                            className="bg-destructive focus:ring-destructive hover:bg-destructive/90"
                        >
                            {deleteMutation.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    {t('customers.actions.deleting', 'Deleting...')}
                                </>
                            ) : (
                                t('customers.actions.delete_confirm', 'Delete')
                            )}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}
