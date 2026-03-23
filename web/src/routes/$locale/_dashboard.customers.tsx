import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { useState } from 'react'
import { customersListQueryOptions, CustomersListParams } from '@/lib/api/query/customers'
import { InternalCustomersCustomerResponse } from '@/lib/api/generated'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow, } from '@/components/ui/table'
import { Plus, Search } from 'lucide-react'
import { NewPagination } from "@/components/pagination.tsx";
import { useTranslation } from 'react-i18next'
import { CustomerActions } from "@/components/customerActions.tsx";
import { CustomerFormDialog } from "@/components/customerFormDialog.tsx";

const customersSearchSchema = z.object({
    page: z.number().catch(1),
    limit: z.number().catch(10),
    search: z.string().optional(),
})

export const Route = createFileRoute('/$locale/_dashboard/customers')({
    validateSearch: (search) => customersSearchSchema.parse(search),

    loaderDeps: ({ search }) => ({
        page: search.page,
        limit: search.limit,
        search: search.search,
    }),

    loader: ({ context: { queryClient }, deps }) => {
        const typedDeps = deps as CustomersListParams;
        return queryClient.ensureQueryData(customersListQueryOptions({
            page: typedDeps.page,
            limit: typedDeps.limit,
            search: typedDeps.search,
        }))
    },

    component: CustomersPage,
})

function CustomersPage() {
    const { t } = useTranslation()
    const navigate = useNavigate({ from: Route.fullPath })
    const searchParams = Route.useSearch() as CustomersListParams

    const customersQuery = useSuspenseQuery(customersListQueryOptions(searchParams))

    const customers = customersQuery.data.customers || []
    const pagination = customersQuery.data.pagination

    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [selectedCustomer, setSelectedCustomer] = useState<InternalCustomersCustomerResponse | null>(null)

    const handleSearch = (term: string) => {
        navigate({
            search: (prev) => ({ ...prev, search: term || undefined, page: 1 }),
            replace: true
        })
    }

    const handlePageChange = (newPage: number) => {
        navigate({ search: (prev) => ({ ...prev, page: newPage }) })
    }

    const openCreateModal = () => {
        setSelectedCustomer(null)
        setIsDialogOpen(true)
    }

    const openEditModal = (customer: InternalCustomersCustomerResponse) => {
        setSelectedCustomer(customer)
        setIsDialogOpen(true)
    }

    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">{t('customers.title', 'Customers')}</h1>
                    <p className="text-muted-foreground">{t('customers.description', 'Manage your customers directory.')}</p>
                </div>
                <Button onClick={openCreateModal}>
                    <Plus className="mr-2 h-4 w-4" /> {t('customers.add_button', 'Add Customer')}
                </Button>
            </div>

            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div className="flex flex-1 items-center gap-2">
                    <div className="relative w-full md:w-[300px]">
                        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                        <Input
                            type="search"
                            placeholder={t('customers.search_placeholder', 'Search customers...')}
                            className="pl-8"
                            defaultValue={searchParams.search}
                            onChange={(e) => handleSearch(e.target.value)}
                        />
                    </div>
                </div>
            </div>

            <div className="rounded-md border bg-card">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>{t('customers.table.name', 'Name')}</TableHead>
                            <TableHead>{t('customers.table.phone', 'Phone')}</TableHead>
                            <TableHead>{t('customers.table.email', 'Email')}</TableHead>
                            <TableHead>{t('customers.table.address', 'Address')}</TableHead>
                            <TableHead className="text-right">{t('customers.table.actions', 'Actions')}</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {customers.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="h-24 text-center">{t('customers.table.empty', 'No customers found.')}</TableCell>
                            </TableRow>
                        ) : (
                            customers.map((customer) => (
                                <TableRow key={customer.id}>
                                    <TableCell className="font-medium">{customer.name}</TableCell>
                                    <TableCell>{customer.phone || '-'}</TableCell>
                                    <TableCell>{customer.email || '-'}</TableCell>
                                    <TableCell>{customer.address || '-'}</TableCell>
                                    <TableCell className="text-right">
                                        <CustomerActions customer={customer} onEdit={() => openEditModal(customer)} />
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>

            {pagination && (
                <NewPagination
                    pagination={pagination}
                    onClickPrev={() => handlePageChange((pagination.current_page || 1) - 1)}
                    onClickNext={() => handlePageChange((pagination.current_page || 1) + 1)}
                />
            )}

            <CustomerFormDialog
                open={isDialogOpen}
                onOpenChange={setIsDialogOpen}
                customerToEdit={selectedCustomer}
            />
        </div>
    )
}
