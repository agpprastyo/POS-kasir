import { Button } from '@/components/ui/button'
import { useShiftContext } from '@/context/ShiftContext'
import { Wallet, Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'

export function ShiftControl() {
    const { isShiftOpen, setOpenShiftModalOpen, setCloseShiftModalOpen, isLoading, shift } = useShiftContext()
    const { t } = useTranslation()

    if (isLoading) {
        return (
            <Button variant="outline" className="w-full justify-start gap-2" disabled>
                <Loader2 className="h-4 w-4 animate-spin" />
                <span className="truncate">{t('shift.loading', 'Loading shift...')}</span>
            </Button>
        )
    }

    if (isShiftOpen) {
        return (
            <Button
                variant="outline"
                className="w-full justify-between group hover:bg-destructive hover:text-destructive-foreground border-destructive/20"
                onClick={() => setCloseShiftModalOpen(true)}
            >
                <div className="flex items-center gap-2">
                    <Wallet className="h-4 w-4" />
                    <span className="truncate">{t('shift.register_open', 'Register Open')}</span>
                </div>
                <span className="text-xs opacity-75 font-mono group-hover:text-destructive-foreground/90">
                    {shift?.start_time ? new Date(shift.start_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : ''}
                </span>
            </Button>
        )
    }

    return (
        <Button
            variant="default"
            className="w-full justify-start gap-2 bg-primary hover:bg-primary/90 text-primary-foreground"
            onClick={() => setOpenShiftModalOpen(true)}
        >
            <Wallet className="h-4 w-4" />
            <span className="truncate">{t('shift.open_register', 'Open Register')}</span>
        </Button>
    )
}
