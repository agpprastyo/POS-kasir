import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface CustomersHeaderProps {
    t: any
    onCreateClick: () => void
}

export function CustomersHeader({ t, onCreateClick }: CustomersHeaderProps) {
    return (
        <div className="flex items-center justify-between">
            <div>
                <h1 className="text-2xl font-bold tracking-tight">{t('customers.title', 'Customers')}</h1>
                <p className="text-muted-foreground">{t('customers.description', 'Manage your customers directory.')}</p>
            </div>
            <Button onClick={onCreateClick}>
                <Plus className="mr-2 h-4 w-4" /> {t('customers.add_button', 'Add Customer')}
            </Button>
        </div>
    )
}
