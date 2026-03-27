import { Filter, Search, LayoutGrid, LayoutList } from 'lucide-react'
import { Input } from '@/components/ui/input'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'

interface ProductFiltersProps {
    trashCount: number
    viewMode: 'list' | 'grid'
    onViewModeChange: (mode: 'list' | 'grid') => void
    searchTerm: string
    onSearchChange: (term: string) => void
    category: number | undefined
    categories: any[]
    onCategoryChange: (value: string) => void
    t: any
    canReadTrash?: boolean
}

export function ProductFilters({
    trashCount, viewMode, onViewModeChange,
    searchTerm, onSearchChange, category, categories, onCategoryChange, t, canReadTrash = true
}: ProductFiltersProps) {
    return (
        <div className="w-full">
            <div className="flex items-center justify-between mb-4">
                <TabsList>
                    <TabsTrigger value="active">{t('products.tabs.active')}</TabsTrigger>
                    {canReadTrash && (
                        <TabsTrigger value="trash">{t('products.tabs.trash', { count: trashCount })}</TabsTrigger>
                    )}
                </TabsList>

                <Tabs value={viewMode} onValueChange={(v) => onViewModeChange(v as 'list' | 'grid')} className="w-[80px]">
                    <TabsList className="grid w-full grid-cols-2">
                        <TabsTrigger value="list" className='px-2' title={t('products.tabs.list')}><LayoutList className="h-4 w-4" /></TabsTrigger>
                        <TabsTrigger value="grid" className='px-2' title={t('products.tabs.grid')}><LayoutGrid className="h-4 w-4" /></TabsTrigger>
                    </TabsList>
                </Tabs>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-4 md:flex-row md:items-center mb-6">
                <div className="relative flex-1 md:max-w-sm">
                    <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                        type="search"
                        placeholder={t('products.filters.search_placeholder')}
                        className="pl-8"
                        value={searchTerm}
                        onChange={(e) => onSearchChange(e.target.value)}
                    />
                </div>

                <Select
                    value={category ? String(category) : 'all'}
                    onValueChange={onCategoryChange}
                >
                    <SelectTrigger className="w-[180px]">
                        <Filter className="mr-2 h-4 w-4" />
                        <SelectValue placeholder={t('products.filters.category')} />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">{t('products.filters.category_all')}</SelectItem>
                        {categories.map((cat: any) => (
                            <SelectItem key={cat.id} value={String(cat.id)}>
                                {cat.name}
                            </SelectItem>
                        ))}
                    </SelectContent>
                </Select>
            </div>
        </div>
    )
}
