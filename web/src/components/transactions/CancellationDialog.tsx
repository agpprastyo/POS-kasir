import { Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"

interface CancellationDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    order: any | null
    cancellationReasons: any[] | undefined
    cancelReasonId: string
    setCancelReasonId: (id: string) => void
    cancelNotes: string
    setCancelNotes: (notes: string) => void
    onConfirm: () => void
    isPending: boolean
    t: any
}

export function CancellationDialog({
    open, onOpenChange, order, cancellationReasons, cancelReasonId,
    setCancelReasonId, cancelNotes, setCancelNotes, onConfirm, isPending, t
}: CancellationDialogProps) {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{t('transactions.cancellation_dialog.title')}</DialogTitle>
                    <DialogDescription>
                        {t('transactions.cancellation_dialog.description')} <b>#{order?.queue_number}</b>
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 py-4">
                    <div className="space-y-2">
                        <label className="text-sm font-medium">{t('transactions.cancellation_dialog.reason_label')}</label>
                        <Select value={cancelReasonId} onValueChange={setCancelReasonId}>
                            <SelectTrigger>
                                <SelectValue placeholder={t('transactions.cancellation_dialog.reason_placeholder')} />
                            </SelectTrigger>
                            <SelectContent>
                                {cancellationReasons?.map((reason) => (
                                    <SelectItem key={reason.id} value={String(reason.id)}>
                                        {reason.reason}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-medium">{t('transactions.cancellation_dialog.notes_label')}</label>
                        <Textarea
                            value={cancelNotes}
                            onChange={(e) => setCancelNotes(e.target.value)}
                            placeholder={t('transactions.cancellation_dialog.notes_placeholder')}
                        />
                    </div>
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>{t('transactions.cancellation_dialog.cancel')}</Button>
                    <Button variant="destructive" onClick={onConfirm} disabled={!cancelReasonId || isPending}>
                        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        {t('transactions.cancellation_dialog.confirm')}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
