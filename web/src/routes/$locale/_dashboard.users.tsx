import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { z } from 'zod'
import { useState } from 'react'
import { usersListQueryOptions } from '@/lib/api/query/user'
import { InternalUserProfileResponse, UsersGetRoleEnum, UsersGetStatusEnum } from '@/lib/api/generated'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow, } from '@/components/ui/table'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue, } from "@/components/ui/select"
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Plus, Search } from 'lucide-react'
import { NewPagination } from "@/components/pagination.tsx";
import { useTranslation } from 'react-i18next'
import { UserActions } from "@/components/userActions.tsx";
import { UserFormDialog } from "@/components/userFormDialog.tsx";


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
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">{t('users.title')}</h1>
                    <p className="text-muted-foreground">{t('users.description')}</p>
                </div>
                <Button onClick={openCreateModal}>
                    <Plus className="mr-2 h-4 w-4" /> {t('users.add_button')}
                </Button>
            </div>

            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div className="flex flex-1 items-center gap-2">
                    <div className="relative w-full md:w-[300px]">
                        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                        <Input
                            type="search"
                            placeholder={t('users.search_placeholder')}
                            className="pl-8"
                            defaultValue={searchParams.search}
                            onChange={(e) => handleSearch(e.target.value)}
                        />
                    </div>
                    <Select value={searchParams.role || 'all'} onValueChange={handleFilterRole}>
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

            <div className="rounded-md border bg-card">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-20">{t('users.table.avatar')}</TableHead>
                            <TableHead>{t('users.table.user')}</TableHead>
                            <TableHead>{t('users.table.role')}</TableHead>
                            <TableHead>{t('users.table.status')}</TableHead>
                            <TableHead className="text-right">{t('users.table.actions')}</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {users.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="h-24 text-center">{t('users.table.empty')}</TableCell>
                            </TableRow>
                        ) : (
                            users.map((user) => (
                                <TableRow key={user.id}>
                                    <TableCell>
                                        <Avatar className="h-9 w-9">
                                            <AvatarImage src={user.avatar || undefined} alt={user.username} />
                                            <AvatarFallback>{user.username?.slice(0, 2).toUpperCase()}</AvatarFallback>
                                        </Avatar>
                                    </TableCell>
                                    <TableCell>
                                        <div className="flex flex-col">
                                            <span className="font-medium">{user.username}</span>
                                            <span className="text-xs text-muted-foreground">{user.email}</span>
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant="outline" className="capitalize">{t(`users.roles.${user.role}`)}</Badge>
                                    </TableCell>
                                    <TableCell>
                                        <Badge
                                            variant={user.is_active ? 'default' : 'destructive'}
                                        >
                                            {user.is_active ? t('users.status.active') : t('users.status.inactive')}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="text-right">
                                        <UserActions user={user} onEdit={() => openEditModal(user)} />
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

            <UserFormDialog
                open={isDialogOpen}
                onOpenChange={setIsDialogOpen}
                userToEdit={selectedUser}
            />
        </div>
    )
}

