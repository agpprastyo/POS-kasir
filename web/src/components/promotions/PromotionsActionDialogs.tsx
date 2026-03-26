import { Loader2 } from "lucide-react"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle
} from "@/components/ui/alert-dialog"

interface PromotionsActionDialogsProps {
    deleteId: string | null
    setDeleteId: (id: string | null) => void
    confirmDelete: () => void
    isDeleting: boolean
    restoreId: string | null
    setRestoreId: (id: string | null) => void
    confirmRestore: () => void
    isRestoring: boolean
    t: any
}

export function PromotionsActionDialogs({
    deleteId,
    setDeleteId,
    confirmDelete,
    isDeleting,
    restoreId,
    setRestoreId,
    confirmRestore,
    isRestoring,
    t
}: PromotionsActionDialogsProps) {
    return (
        <>
            <AlertDialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>{t('promotions.delete_title', 'Delete Promotion')}</AlertDialogTitle>
                        <AlertDialogDescription>
                            {t('promotions.delete_confirm', 'Are you sure you want to delete this promotion?')}
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel disabled={isDeleting}>{t('common.cancel', 'Cancel')}</AlertDialogCancel>
                        <AlertDialogAction
                            onClick={confirmDelete}
                            disabled={isDeleting}
                            className="bg-destructive hover:bg-destructive/90 focus:ring-destructive"
                        >
                            {isDeleting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                            {t('common.delete', 'Delete')}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>

            <AlertDialog open={!!restoreId} onOpenChange={(open) => !open && setRestoreId(null)}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>{t('promotions.restore_title', 'Restore Promotion')}</AlertDialogTitle>
                        <AlertDialogDescription>
                            {t('promotions.restore_confirm', 'Are you sure you want to restore this promotion?')}
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel disabled={isRestoring}>{t('common.cancel', 'Cancel')}</AlertDialogCancel>
                        <AlertDialogAction
                            onClick={confirmRestore}
                            disabled={isRestoring}
                        >
                            {isRestoring ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                            {t('common.restore', 'Restore')}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}
