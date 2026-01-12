import {createFileRoute, useNavigate} from '@tanstack/react-router'
import {useSuspenseQuery} from '@tanstack/react-query'
import {z} from 'zod'
import React, {useEffect, useState} from 'react'
import {
    useCreateUserMutation,
    useDeleteUserMutation,
    usersListQueryOptions,
    useToggleUserStatusMutation,
    useUpdateUserMutation
} from '@/lib/api/query/user'
import {
    POSKasirInternalDtoCreateUserRequest,
    POSKasirInternalDtoProfileResponse,
    POSKasirInternalDtoUpdateUserRequest,
    UsersGetRoleEnum,
    UsersGetStatusEnum
} from '@/lib/api/generated'
import {queryClient} from '@/lib/queryClient'

import {Button} from '@/components/ui/button'
import {Input} from '@/components/ui/input'
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from '@/components/ui/table'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue,} from "@/components/ui/select"
import {Badge} from '@/components/ui/badge'
import {Avatar, AvatarFallback, AvatarImage} from '@/components/ui/avatar'
import {Label} from '@/components/ui/label'
import {Loader2, MoreHorizontal, Pencil, Plus, Power, Search, Trash2} from 'lucide-react'
import {toast} from "sonner"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle
} from "@/components/ui/alert-dialog.tsx";
import {NewPagination} from "@/components/pagination.tsx";


const usersSearchSchema = z.object({
    page: z.number().catch(1),
    limit: z.number().catch(10),
    search: z.string().optional(),
    role: z.enum(UsersGetRoleEnum).optional(),
    status: z.enum(UsersGetStatusEnum).optional(),
})

// --- ROUTE DEFINITION ---
export const Route = createFileRoute('/_dashboard/users')({
    validateSearch: (search) => usersSearchSchema.parse(search),

    loaderDeps: ({search}) => ({
        page: search.page,
        limit: search.limit,
        search: search.search,
        role: search.role,
        status: search.status,
    }),

    loader: ({deps}) => {
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
    const navigate = useNavigate({from: Route.fullPath})
    const searchParams = Route.useSearch()

    const usersQuery = useSuspenseQuery(usersListQueryOptions(searchParams))

    const users = usersQuery.data.users || []
    const pagination = usersQuery.data.pagination

    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [selectedUser, setSelectedUser] = useState<POSKasirInternalDtoProfileResponse | null>(null)

    // Handlers
    const handleSearch = (term: string) => {
        navigate({
            search: (prev) => ({...prev, search: term || undefined, page: 1}),
            replace: true
        })
    }

    const handleFilterRole = (role: string) => {
        navigate({
            search: (prev) => ({...prev, role: role === 'all' ? undefined : role as UsersGetRoleEnum, page: 1})
        })
    }

    const handlePageChange = (newPage: number) => {
        navigate({search: (prev) => ({...prev, page: newPage})})
    }

    const openCreateModal = () => {
        setSelectedUser(null)
        setIsDialogOpen(true)
    }

    const openEditModal = (user: POSKasirInternalDtoProfileResponse) => {
        setSelectedUser(user)
        setIsDialogOpen(true)
    }

    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight">Users</h1>
                    <p className="text-muted-foreground">Manage user access and roles.</p>
                </div>
                <Button onClick={openCreateModal}>
                    <Plus className="mr-2 h-4 w-4"/> Add User
                </Button>
            </div>

            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div className="flex flex-1 items-center gap-2">
                    <div className="relative w-full md:w-[300px]">
                        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground"/>
                        <Input
                            type="search"
                            placeholder="Search user..."
                            className="pl-8"
                            defaultValue={searchParams.search}
                            onChange={(e) => handleSearch(e.target.value)}
                        />
                    </div>
                    <Select value={searchParams.role || 'all'} onValueChange={handleFilterRole}>
                        <SelectTrigger className="w-[150px]">
                            <SelectValue placeholder="Role"/>
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">All Roles</SelectItem>
                            <SelectItem value={UsersGetRoleEnum.Admin}>Admin</SelectItem>
                            <SelectItem value={UsersGetRoleEnum.Manager}>Manager</SelectItem>
                            <SelectItem value={UsersGetRoleEnum.Cashier}>Cashier</SelectItem>
                        </SelectContent>
                    </Select>
                </div>
            </div>

            <div className="rounded-md border bg-card">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-[80px]">Avatar</TableHead>
                            <TableHead>User</TableHead>
                            <TableHead>Role</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {users.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="h-24 text-center">No users found.</TableCell>
                            </TableRow>
                        ) : (
                            users.map((user) => (
                                <TableRow key={user.id}>
                                    <TableCell>
                                        <Avatar className="h-9 w-9">
                                            <AvatarImage src={user.avatar || undefined} alt={user.username}/>
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
                                        <Badge variant="outline" className="capitalize">{user.role}</Badge>
                                    </TableCell>
                                    <TableCell>
                                        <Badge
                                            variant={user.is_active ? 'default' : 'destructive'}
                                            className={user.is_active ? 'bg-green-500 hover:bg-green-600' : ''}
                                        >
                                            {user.is_active ? 'Active' : 'Inactive'}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="text-right">
                                        <UserActions user={user} onEdit={() => openEditModal(user)}/>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>


            <NewPagination pagination={pagination} onClick={() => handlePageChange((pagination.current_page || 1) - 1)}
                           onClick1={() => handlePageChange((pagination.current_page || 1) + 1)}/>

            <UserFormDialog
                open={isDialogOpen}
                onOpenChange={setIsDialogOpen}
                userToEdit={selectedUser}
            />
        </div>
    )
}

// --- ACTION DROPDOWN ---
function UserActions({user, onEdit}: { user: POSKasirInternalDtoProfileResponse, onEdit: () => void }) {
    const deleteMutation = useDeleteUserMutation()
    const toggleMutation = useToggleUserStatusMutation()

    const [showDeleteDialog, setShowDeleteDialog] = useState(false)

    const handleDelete = async (e: React.MouseEvent) => {
        e.preventDefault()

        deleteMutation.mutate(user.id!)
        setShowDeleteDialog(false)
    }

    return (
        <>
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">Open menu</span>
                        <MoreHorizontal className="h-4 w-4"/>
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuLabel>Actions</DropdownMenuLabel>
                    <DropdownMenuItem onClick={onEdit}>
                        <Pencil className="mr-2 h-4 w-4"/> Edit Details
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => toggleMutation.mutate(user.id!)}>
                        <Power className="mr-2 h-4 w-4"/>
                        {user.is_active ? 'Deactivate' : 'Activate'} Account
                    </DropdownMenuItem>
                    <DropdownMenuSeparator/>

                    <DropdownMenuItem
                        onSelect={(e) => {
                            e.preventDefault()
                            setShowDeleteDialog(true)
                        }}
                        className="text-red-600 focus:text-red-600 cursor-pointer"
                    >
                        <Trash2 className="mr-2 h-4 w-4"/> Delete User
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>


            <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                        <AlertDialogDescription>
                            This action cannot be undone. This will permanently delete the user account
                            <span className="font-bold text-foreground"> "{user.username}" </span>
                            and remove their data from our servers.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel disabled={deleteMutation.isPending}>Cancel</AlertDialogCancel>

                        <AlertDialogAction
                            onClick={handleDelete}
                            disabled={deleteMutation.isPending}
                            className="bg-red-600 focus:ring-red-600 hover:bg-red-700"
                        >
                            {deleteMutation.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                    Deleting...
                                </>
                            ) : (
                                "Delete"
                            )}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}

// --- FORM DIALOG (CREATE / EDIT) ---
function UserFormDialog({open, onOpenChange, userToEdit}: {
    open: boolean,
    onOpenChange: (open: boolean) => void,
    userToEdit: POSKasirInternalDtoProfileResponse | null
}) {
    const createMutation = useCreateUserMutation()
    const updateMutation = useUpdateUserMutation()

    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        role: UsersGetRoleEnum.Cashier as UsersGetRoleEnum
    })

    useEffect(() => {
        if (open) {
            if (userToEdit) {
                setFormData({
                    username: userToEdit.username || '',
                    email: userToEdit.email || '',
                    password: '',
                    role: (userToEdit.role as UsersGetRoleEnum) || UsersGetRoleEnum.Cashier
                })
            } else {

                setFormData({
                    username: '',
                    email: '',
                    password: '',
                    role: UsersGetRoleEnum.Cashier
                })
            }
        }
    }, [open, userToEdit])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        try {
            if (userToEdit) {
                const payload: POSKasirInternalDtoUpdateUserRequest = {
                    username: formData.username,
                    email: formData.email,
                    role: formData.role,

                }

                await updateMutation.mutateAsync({id: userToEdit.id!, body: payload})
                toast.success("User updated successfully")
            } else {
                // --- CREATE LOGIC ---
                const payload: POSKasirInternalDtoCreateUserRequest = {
                    username: formData.username,
                    email: formData.email,
                    password: formData.password,
                    role: formData.role,
                    is_active: true
                }
                await createMutation.mutateAsync(payload)
                toast.success("User created successfully")
            }
            onOpenChange(false)
        } catch (error) {
            console.error(error)
            toast.error("Failed to save user")
        }
    }

    const isSubmitting = createMutation.isPending || updateMutation.isPending

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>{userToEdit ? 'Edit User' : 'Create New User'}</DialogTitle>
                    <DialogDescription>
                        {userToEdit ? "Update user details. Password cannot be changed here." : "Add a new user to the system."}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        {/* Username */}
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="username" className="text-right">Username</Label>
                            <Input
                                id="username"
                                value={formData.username}
                                onChange={e => setFormData({...formData, username: e.target.value})}
                                className="col-span-3"
                                required
                            />
                        </div>
                        {/* Email */}
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="email" className="text-right">Email</Label>
                            <Input
                                id="email"
                                type="email"
                                value={formData.email}
                                onChange={e => setFormData({...formData, email: e.target.value})}
                                className="col-span-3"
                                required
                            />
                        </div>
                        {/* Role */}
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="role" className="text-right">Role</Label>
                            <div className="col-span-3">
                                <Select
                                    value={formData.role}
                                    onValueChange={(val: UsersGetRoleEnum) => setFormData({...formData, role: val})}
                                >
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select a role"/>
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value={UsersGetRoleEnum.Admin}>Admin</SelectItem>
                                        <SelectItem value={UsersGetRoleEnum.Manager}>Manager</SelectItem>
                                        <SelectItem value={UsersGetRoleEnum.Cashier}>Cashier</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>

                        {!userToEdit && (
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="password" className="text-right">Password</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    placeholder="Strong password"
                                    value={formData.password}
                                    onChange={e => setFormData({...formData, password: e.target.value})}
                                    className="col-span-3"
                                    required
                                    minLength={6}
                                />
                            </div>
                        )}
                    </div>
                    <DialogFooter>
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin"/>}
                            {userToEdit ? 'Save Changes' : 'Create User'}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}