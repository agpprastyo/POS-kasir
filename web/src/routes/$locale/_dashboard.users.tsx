import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { useState } from 'react'
import { usersListQueryOptions } from '@/lib/api/query/user'
import { InternalUserProfileResponse, UsersGetRoleEnum, UsersGetStatusEnum } from '@/lib/api/generated'
import { NewPagination } from "@/components/pagination.tsx";
import { useTranslation } from 'react-i18next'
import { UsersHeader } from "@/components/users/UsersHeader"
import { UsersFilters } from "@/components/users/UsersFilters"
import { UsersTable } from "@/components/users/UsersTable"
import { UserActions } from "@/components/users/UserActions"
import { UserFormDialog } from "@/components/users/UserFormDialog"


const usersSearchSchema = z.object({
    page: z.number().catch(1),
    limit: z.number().catch(10),
    search: z.string().optional(),
    role: z.enum(UsersGetRoleEnum).optional(),
    status: z.enum(UsersGetStatusEnum).optional(),
})

// --- ROUTE DEFINITION ---
export const Route = createFileRoute('/$locale/_dashboard/users')({
    validateSearch: (search) => usersSearchSchema.parse(search),

    loaderDeps: ({ search }) => ({
        page: search.page,
        limit: search.limit,
        search: search.search,
        role: search.role,
        status: search.status,
    }),

    loader: ({ context: { queryClient }, deps }) => {
        return queryClient.ensureQueryData(usersListQueryOptions({
            page: deps.page,
            limit: deps.limit,
            search: deps.search,
            role: deps.role,
            status: deps.status,
        }))
    },

    component: UsersPage,
})


// --- MAIN COMPONENT ---
function UsersPage() {
    const { t } = useTranslation()
    const navigate = useNavigate({ from: Route.fullPath })
    const searchParams = Route.useSearch()

    const usersQuery = useSuspenseQuery(usersListQueryOptions(searchParams))

    const users = usersQuery.data.users || []
    const pagination = usersQuery.data.pagination

    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [selectedUser, setSelectedUser] = useState<InternalUserProfileResponse | null>(null)

    // Handlers
    const handleSearch = (term: string) => {
        navigate({
            search: (prev) => ({ ...prev, search: term || undefined, page: 1 }),
            replace: true
        })
    }

    const handleFilterRole = (role: string) => {
        navigate({
            search: (prev) => ({ ...prev, role: role === 'all' ? undefined : role as UsersGetRoleEnum, page: 1 })
        })
    }

    const handlePageChange = (newPage: number) => {
        navigate({ search: (prev) => ({ ...prev, page: newPage }) })
    }

    const openCreateModal = () => {
        setSelectedUser(null)
        setIsDialogOpen(true)
    }

    const openEditModal = (user: InternalUserProfileResponse) => {
        setSelectedUser(user)
        setIsDialogOpen(true)
    }

    return (
        <div className="flex flex-col gap-6">
            <UsersHeader 
                t={t}
                onCreateClick={openCreateModal}
            />

            <UsersFilters 
                t={t}
                search={searchParams.search}
                onSearch={handleSearch}
                role={searchParams.role}
                onRoleChange={handleFilterRole}
            />

            <UsersTable 
                users={users}
                t={t}
                renderActions={(user) => (
                    <UserActions user={user} onEdit={() => openEditModal(user)} />
                )}
            />

            {pagination && (
                <NewPagination
                    pagination={pagination}
                    onClickPrev={() => handlePageChange((pagination.current_page || 1) - 1)}
                    onClickNext={() => handlePageChange((pagination.current_page || 1) + 1)}
                />
            )}

            <UserFormDialog
                open={isDialogOpen}
                onOpenChange={setIsDialogOpen}
                userToEdit={selectedUser}
            />
        </div>
    )
}

