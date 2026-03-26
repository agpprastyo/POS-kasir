import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface UsersHeaderProps {
    t: any
    onCreateClick: () => void
}

export function UsersHeader({ t, onCreateClick }: UsersHeaderProps) {
    return (
        <div className="flex items-center justify-between">
            <div>
                <h1 className="text-2xl font-bold tracking-tight">{t('users.title')}</h1>
                <p className="text-muted-foreground">{t('users.description')}</p>
            </div>
            <Button onClick={onCreateClick}>
                <Plus className="mr-2 h-4 w-4" /> {t('users.add_button')}
            </Button>
        </div>
    )
}
