import { Search } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { UsersGetRoleEnum } from '@/lib/api/generated'

interface UsersFiltersProps {
    t: any
    search: string | undefined
    onSearch: (term: string) => void
    role: string | undefined
    onRoleChange: (role: string) => void
}

export function UsersFilters({ t, search, onSearch, role, onRoleChange }: UsersFiltersProps) {
    return (
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="flex flex-1 items-center gap-2">
                <div className="relative w-full md:w-[300px]">
                    <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                        type="search"
                        placeholder={t('users.search_placeholder')}
                        className="pl-8"
                        defaultValue={search}
                        onChange={(e) => onSearch(e.target.value)}
                    />
                </div>
                <Select value={role || 'all'} onValueChange={onRoleChange}>
                    <SelectTrigger className="w-[150px]">
                        <SelectValue placeholder={t('users.role_filter')} />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">{t('users.role_all')}</SelectItem>
                        <SelectItem value={UsersGetRoleEnum.Admin}>{t('users.roles.admin')}</SelectItem>
                        <SelectItem value={UsersGetRoleEnum.Manager}>{t('users.roles.manager')}</SelectItem>
                        <SelectItem value={UsersGetRoleEnum.Cashier}>{t('users.roles.cashier')}</SelectItem>
                    </SelectContent>
                </Select>
            </div>
        </div>
    )
}
