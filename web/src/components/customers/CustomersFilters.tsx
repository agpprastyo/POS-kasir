import { Search } from 'lucide-react'
import { Input } from '@/components/ui/input'

interface CustomersFiltersProps {
    t: any
    search: string | undefined
    onSearch: (term: string) => void
}

export function CustomersFilters({ t, search, onSearch }: CustomersFiltersProps) {
    return (
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="flex flex-1 items-center gap-2">
                <div className="relative w-full md:w-[300px]">
                    <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                        type="search"
                        placeholder={t('customers.search_placeholder', 'Search customers...')}
                        className="pl-8"
                        defaultValue={search}
                        onChange={(e) => onSearch(e.target.value)}
                    />
                </div>
            </div>
        </div>
    )
}
