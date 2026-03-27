import { Search } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"

interface ProductSearchProps {
    searchTerm: string
    onSearchChange: (value: string) => void
    selectedCategory: string
    onCategoryChange: (category: string) => void
    categories: { id: string, name: string }[]
    t: any
}

export function ProductSearch({
    searchTerm, onSearchChange, selectedCategory, onCategoryChange, categories, t
}: ProductSearchProps) {
    return (
        <div className="flex flex-col gap-4">
            <div className="relative">
                <Search className="absolute left-2.5 top-3 h-6 w-6 text-muted-foreground" />
                <Input
                    type="search"
                    placeholder={t('order.search_placeholder')}
                    className="pl-12 py-6"
                    value={searchTerm}
                    onChange={(e) => onSearchChange(e.target.value)}
                />
            </div>

            <div className="px-0">
                <Tabs value={selectedCategory} onValueChange={onCategoryChange} className="w-full ">
                    <TabsList className="w-full justify-start overflow-x-auto h-auto p-1 no-scrollbar text-sm border-none">
                        <TabsTrigger
                            value="all"
                            className="rounded-full px-4 data-[state=active]:bg-primary data-[state=active]:text-primary-foreground border"
                        >
                            {t('order.all_categories')}
                        </TabsTrigger>
                        {categories.map(category => (
                            <TabsTrigger
                                key={category.id}
                                value={category.id}
                                className="rounded-full px-4 data-[state=active]:bg-primary data-[state=active]:text-primary-foreground border"
                            >
                                {category.name}
                            </TabsTrigger>
                        ))}
                    </TabsList>
                </Tabs>
            </div>
        </div>
    )
}
