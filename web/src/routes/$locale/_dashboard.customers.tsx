import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { useState } from 'react'
import { customersListQueryOptions, CustomersListParams } from '@/lib/api/query/customers'
import { InternalCustomersCustomerResponse } from '@/lib/api/generated'
import { NewPagination } from "@/components/pagination.tsx";
import { useTranslation } from 'react-i18next'
import { CustomersHeader } from "@/components/customers/CustomersHeader"
import { CustomersFilters } from "@/components/customers/CustomersFilters"
import { CustomersTable } from "@/components/customers/CustomersTable"
import { CustomerActions } from "@/components/customers/CustomerActions"
import { CustomerFormDialog } from "@/components/customers/CustomerFormDialog"

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
            <CustomersHeader 
                t={t}
                onCreateClick={openCreateModal}
            />

            <CustomersFilters 
                t={t}
                search={searchParams.search}
                onSearch={handleSearch}
            />

            <CustomersTable 
                customers={customers}
                t={t}
                renderActions={(customer) => (
                    <CustomerActions customer={customer} onEdit={() => openEditModal(customer)} />
                )}
            />

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
