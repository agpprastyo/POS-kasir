import { User as UserIcon, Receipt } from 'lucide-react'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"

interface TransactionFiltersProps {
    users: any[]
    selectedUserId: string
    onUserChange: (value: string) => void
    t: any
}

export function TransactionFilters({
    users, selectedUserId, onUserChange, t
}: TransactionFiltersProps) {
    return (
        <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center">
                    <Receipt className="h-5 w-5 text-primary" />
                </div>
                <div>
                    <h1 className="text-2xl font-bold tracking-tight font-heading">{t('transactions.title')}</h1>
                    <p className="text-sm text-muted-foreground">{t('transactions.subtitle')}</p>
                </div>
            </div>

            {/* User Filter */}
            <div className="flex items-center gap-2">
                <UserIcon className="h-4 w-4 text-muted-foreground" />
                <Select value={selectedUserId} onValueChange={onUserChange}>
                    <SelectTrigger className="w-[180px] rounded-xl">
                        <SelectValue placeholder={t('transactions.filter_user')} />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">{t('transactions.filter_user')}</SelectItem>
                        {users.map(user => (
                            <SelectItem key={user.id} value={user.id || ''}>{user.username}</SelectItem>
                        ))}
                    </SelectContent>
                </Select>
            </div>
        </div>
    )
}
