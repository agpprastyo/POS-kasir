import {useState} from 'react'
import {createFileRoute, getRouteApi, Link, Outlet, redirect, RegisteredRouter, useRouter} from '@tanstack/react-router'
import {FileText, LayoutDashboard, LogOut, Menu, Package, Settings, User as UserIcon, Zap} from 'lucide-react'
import {Button} from '@/components/ui/button'
import {Sheet, SheetContent, SheetTrigger} from '@/components/ui/sheet'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {Avatar, AvatarFallback, AvatarImage} from '@/components/ui/avatar'
import {useAuth} from '@/lib/auth/AuthContext'
import {meQueryOptions} from '@/lib/api/query/auth'
import {queryClient} from '@/lib/queryClient'
import {POSKasirInternalRepositoryUserRole} from '@/lib/api/generated/models/poskasir-internal-repository-user-role'
import {FileRouteByToPath} from "@tanstack/router-core/src/routeInfo.ts";

// 1. Definisikan Route API helper
const routeApi = getRouteApi('/_dashboard')

export const Route = createFileRoute('/_dashboard' as FileRouteByToPath<any, any>)({
    loader: async () => {
        try {
            return await queryClient.ensureQueryData(meQueryOptions())
        } catch (error: any) {
            const status =
                error?.response?.status ??
                error?.status ??
                error?.cause?.status
            if (status === 401) {
                throw redirect({ to: '/login' } as RegisteredRouter)
            }
            throw error
        }
    },
    component: DashboardLayout,
})

function DashboardLayout() {
    const auth = useAuth()
    const router = useRouter()

    const [isLoggingOut, setIsLoggingOut] = useState(false)

    const user = routeApi.useLoaderData()

    const profile = user as any
    const userRole = (profile?.role ?? profile?.data?.role) as POSKasirInternalRepositoryUserRole | undefined
    const userAvatar = profile?.avatar ?? profile?.data?.avatar
    const userName = profile?.username ?? profile?.data?.username ?? 'User'

    console.log('[Dashboard] User Data (Loader):', user)
    console.log('[Dashboard] User Role:', userRole)


    const handleLogout = async () => {
        if (isLoggingOut) return
        setIsLoggingOut(true)

        try {

            await auth.logout()
        } catch (error) {
            console.error("Gagal logout di server (mungkin sudah expired):", error)
        } finally {
            queryClient.clear()

            await router.navigate({ to: '/login', replace: true })
            setIsLoggingOut(false)
        }
    }

    type DashboardMenuItem = {
        label: string
        icon: any
        to: string
        allowedRoles: POSKasirInternalRepositoryUserRole[]
    }

    const menuItems: DashboardMenuItem[] = [
        {
            label: 'Summary',
            icon: LayoutDashboard,
            to: '/',
            allowedRoles: [
                POSKasirInternalRepositoryUserRole.UserRoleAdmin,
                POSKasirInternalRepositoryUserRole.UserRoleManager,
                POSKasirInternalRepositoryUserRole.UserRoleCashier
            ]
        },
        {
            label: 'Product',
            icon: Package,
            to: '/product',
            allowedRoles: [
                POSKasirInternalRepositoryUserRole.UserRoleAdmin,
                POSKasirInternalRepositoryUserRole.UserRoleManager
            ]
        },
        {
            label: 'Reports',
            icon: FileText,
            to: '/reports',
            allowedRoles: [
                POSKasirInternalRepositoryUserRole.UserRoleAdmin
            ]
        },
        {
            label: 'Users',
            icon: UserIcon,
            to: '/users',
            allowedRoles: [
                POSKasirInternalRepositoryUserRole.UserRoleAdmin
                ,POSKasirInternalRepositoryUserRole.UserRoleManager
            ]
        },
        {
            label: 'Settings',
            icon: Settings,
            to: '/settings',
            allowedRoles: [
                POSKasirInternalRepositoryUserRole.UserRoleAdmin
            ]
        },
    ]

    const filteredMenu = menuItems.filter(item =>
        userRole && item.allowedRoles.includes(userRole)
    )

    return (
        <div className="grid min-h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]">
            {/* --- DESKTOP SIDEBAR --- */}
            <div className="hidden border-r bg-muted/40 md:block">
                <div className="flex h-full max-h-screen flex-col gap-2">
                    <div className="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
                        <Link to="/" className="flex items-center gap-2 font-semibold">
                            <Zap className="h-6 w-6" />
                            <span className="">Acme Inc</span>
                        </Link>
                    </div>
                    <div className="flex-1">
                        <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
                            {filteredMenu.map((item) => (
                                <Link
                                    key={item.to}
                                    to={item.to}
                                    className="flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary [&.active]:bg-muted [&.active]:text-primary"
                                >
                                    <item.icon className="h-4 w-4" />
                                    {item.label}
                                </Link>
                            ))}
                        </nav>
                    </div>
                </div>
            </div>

            {/* --- MAIN CONTENT AREA --- */}
            <div className="flex flex-col">
                {/* HEADER / TOPBAR */}
                <header className="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6">
                    <Sheet>
                        <SheetTrigger asChild>
                            <Button
                                variant="outline"
                                size="icon"
                                className="shrink-0 md:hidden"
                            >
                                <Menu className="h-5 w-5" />
                                <span className="sr-only">Toggle navigation menu</span>
                            </Button>
                        </SheetTrigger>
                        <SheetContent side="left" className="flex flex-col">
                            <nav className="grid gap-2 text-lg font-medium">
                                <Link
                                    to="/"
                                    className="flex items-center gap-2 text-lg font-semibold mb-4"
                                >
                                    <Zap className="h-6 w-6" />
                                    <span className="sr-only">Acme Inc</span>
                                </Link>
                                {/* Mobile Menu Filtered */}
                                {filteredMenu.map((item) => (
                                    <Link
                                        key={item.to}
                                        to={item.to}
                                        className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground [&.active]:bg-muted [&.active]:text-foreground"
                                    >
                                        <item.icon className="h-5 w-5" />
                                        {item.label}
                                    </Link>
                                ))}
                            </nav>
                        </SheetContent>
                    </Sheet>

                    <div className="w-full flex-1">
                        {/* Search bar could go here */}
                    </div>

                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="secondary" size="icon" className="rounded-full">
                                <Avatar className="h-8 w-8">
                                    {/* Menggunakan userAvatar jika ada, fallback ke default github jika kosong */}
                                    <AvatarImage src={userAvatar || "https://github.com/shadcn.png"} alt={userName} />
                                    <AvatarFallback><UserIcon className="h-4 w-4"/></AvatarFallback>
                                </Avatar>
                                <span className="sr-only">Toggle user menu</span>
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuLabel>My Account ({userRole})</DropdownMenuLabel>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem>Settings</DropdownMenuItem>
                            <DropdownMenuItem>Support</DropdownMenuItem>
                            <DropdownMenuSeparator />

                            {/* Gunakan state lokal isLoggingOut untuk UI feedback */}
                            <DropdownMenuItem
                                onClick={handleLogout}
                                disabled={isLoggingOut}
                                className="text-red-600 focus:text-red-600 cursor-pointer"
                            >
                                <LogOut className="mr-2 h-4 w-4"/>
                                {isLoggingOut ? 'Logging out...' : 'Logout'}
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </header>

                <main className="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
                    <Outlet />
                </main>
            </div>
        </div>
    )
}