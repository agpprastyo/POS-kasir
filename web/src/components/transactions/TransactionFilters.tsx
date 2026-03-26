import { User as UserIcon } from 'lucide-react'
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
            <div>
                <h1 className="text-3xl font-bold tracking-tight">{t('transactions.title')}</h1>
                <p className="text-muted-foreground">{t('transactions.subtitle')}</p>
            </div>

            {/* User Filter */}
            <div className="flex items-center gap-2">
                <UserIcon className="h-4 w-4 text-muted-foreground" />
                <Select value={selectedUserId} onValueChange={onUserChange}>
                    <SelectTrigger className="w-[180px]">
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
