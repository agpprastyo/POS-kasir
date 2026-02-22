import { InternalUserProfileResponse } from "@/lib/api/generated";
import { useTranslation } from "react-i18next";
import { useDeleteUserMutation, useToggleUserStatusMutation } from "@/lib/api/query/user.ts";
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
import { Loader2, MoreHorizontal, Pencil, Power, Trash2 } from "lucide-react";
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

export function UserActions({ user, onEdit }: { user: InternalUserProfileResponse, onEdit: () => void }) {
    const { t } = useTranslation()
    const deleteMutation = useDeleteUserMutation()
    const toggleMutation = useToggleUserStatusMutation()

    const [showDeleteDialog, setShowDeleteDialog] = useState(false)

    const handleDelete = async (e: React.MouseEvent) => {
        e.preventDefault()

        deleteMutation.mutate(user.id!)
        setShowDeleteDialog(false)
    }

    return (
        <>
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">{t('users.actions.open_menu')}</span>
                        <MoreHorizontal className="h-4 w-4" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuLabel>{t('users.table.actions')}</DropdownMenuLabel>
                    <DropdownMenuItem onClick={onEdit}>
                        <Pencil className="mr-2 h-4 w-4" /> {t('users.actions.edit')}
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => toggleMutation.mutate(user.id!)}>
                        <Power className="mr-2 h-4 w-4" />
                        {user.is_active ? t('users.actions.deactivate') : t('users.actions.activate')}
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />

                    <DropdownMenuItem
                        onSelect={(e) => {
                            e.preventDefault()
                            setShowDeleteDialog(true)
                        }}
                        className="text-destructive focus:text-destructive cursor-pointer"
                    >
                        <Trash2 className="mr-2 h-4 w-4" /> {t('users.actions.delete')}
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>


            <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>{t('users.actions.delete_title')}</AlertDialogTitle>
                        <AlertDialogDescription>
                            {t('users.actions.delete_desc')}
                            <span className="font-bold text-foreground"> "{user.username}" </span>
                            {t('users.actions.delete_desc_2')}
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel
                            disabled={deleteMutation.isPending}>{t('users.actions.cancel')}</AlertDialogCancel>

                        <AlertDialogAction
                            onClick={handleDelete}
                            disabled={deleteMutation.isPending}
                            className="bg-destructive focus:ring-destructive hover:bg-destructive/90"
                        >
                            {deleteMutation.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    {t('users.actions.deleting')}
                                </>
                            ) : (
                                t('users.actions.delete_confirm')
                            )}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}