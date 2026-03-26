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

interface RefundDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    order: any | null
    refundReason: string
    setRefundReason: (reason: string) => void
    onConfirm: () => void
    isPending: boolean
    t: any
}

export function RefundDialog({
    open, onOpenChange, order, refundReason, setRefundReason, onConfirm, isPending, t
}: RefundDialogProps) {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{t('transactions.refund_dialog.title', 'Refund Order')}</DialogTitle>
                    <DialogDescription>
                        {t('transactions.refund_dialog.description', 'Are you sure you want to refund order')} <b>#{order?.queue_number}</b>? {t('transactions.refund_dialog.warning', 'This will revert items and totals.')}
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 py-4">
                    <div className="space-y-2">
                        <label className="text-sm font-medium">{t('transactions.refund_dialog.reason_label', 'Reason for Refund')}</label>
                        <Textarea
                            value={refundReason}
                            onChange={(e) => setRefundReason(e.target.value)}
                            placeholder={t('transactions.refund_dialog.reason_placeholder', 'Enter refund reason details...')}
                        />
                    </div>
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>{t('transactions.cancellation_dialog.cancel')}</Button>
                    <Button variant="destructive" onClick={onConfirm} disabled={!refundReason || isPending}>
                        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        {t('transactions.actions_button.refund', 'Refund')}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
