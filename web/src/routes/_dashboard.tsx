import { useState } from 'react'
import { createFileRoute, getRouteApi, Link, Outlet, redirect, RegisteredRouter, useRouter } from '@tanstack/react-router'
import { FileText, LayoutDashboard, LogOut, Menu, Package, Settings, ShoppingCart, User as UserIcon, Zap, Receipt } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { useAuth } from '@/lib/auth/AuthContext'
import { meQueryOptions } from '@/lib/api/query/auth'
import { queryClient } from '@/lib/queryClient'
import { POSKasirInternalRepositoryUserRole } from '@/lib/api/generated/models/poskasir-internal-repository-user-role'





export const Route = createFileRoute('/_dashboard')({
    loader: async () => {
        try {
            return await queryClient.ensureQueryData(meQueryOptions())
        } catch (error: any) {
            const status = error?.response?.status ?? error?.status
            if (status === 401) {
                throw redirect({ to: '/login' })
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

    const user = auth.user


    const userRole = user?.role
    const userAvatar = user?.avatar
    const userName = user?.username ?? 'User'

    const handleLogout = async () => {
        if (isLoggingOut) return
        setIsLoggingOut(true)

        try {

            await auth.logout()

            await router.navigate({ to: '/login', replace: true })
        } catch (error) {
            console.error("Logout UI error:", error)
        } finally {
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
            label: 'POS',
            icon: ShoppingCart,
            to: '/order',
            allowedRoles: [
                POSKasirInternalRepositoryUserRole.UserRoleAdmin,
                POSKasirInternalRepositoryUserRole.UserRoleManager,
                POSKasirInternalRepositoryUserRole.UserRoleCashier
            ]
        },
        {
            label: 'Transactions',
            icon: Receipt,
            to: '/transactions',
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
                , POSKasirInternalRepositoryUserRole.UserRoleManager
            ]
        },
        {
            label: 'Settings',
            icon: Settings,
            to: '/settings',
            allowedRoles: [
                POSKasirInternalRepositoryUserRole.UserRoleAdmin,
                POSKasirInternalRepositoryUserRole.UserRoleManager,
                POSKasirInternalRepositoryUserRole.UserRoleCashier
            ]
        },
        {
            label: 'Account',
            icon: UserIcon,
            to: '/account',
            allowedRoles: [
                POSKasirInternalRepositoryUserRole.UserRoleAdmin,
                POSKasirInternalRepositoryUserRole.UserRoleManager,
                POSKasirInternalRepositoryUserRole.UserRoleCashier
            ]
        },
    ]

    const filteredMenu = menuItems.filter(item =>
        userRole && item.allowedRoles.includes(userRole)
    )

    return (
        <div className="grid h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr] overflow-hidden">
            {/* --- DESKTOP SIDEBAR --- */}
            <div className="hidden md:block">
                <div className="flex h-full max-h-screen flex-col gap-2">
                    <div className="flex h-14 items-center  px-4 lg:h-[60px] lg:px-6">
                        <Link to="/" className="flex items-center gap-2 font-semibold">
                            <Zap className="h-8 w-8" />
                            <span className="text-2xl">Acme Inc</span>
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
                    <div className="mt-auto p-4">

                        <div className="rounded-2xl w-full flex items-center justify-between px-2 gap-2 pl-4 aspect-auto h-12 border">


                            <div className="flex items-center gap-4 cursor-default">
                                <Avatar className="h-8 w-8">
                                    <AvatarImage src={userAvatar || undefined} alt={userName} />
                                    <AvatarFallback><UserIcon className="h-4 w-4" /></AvatarFallback>
                                </Avatar>
                                <div className="flex flex-col items-start truncate text-sm">
                                    <span className="font-semibold">{userName}</span>
                                    <span className="text-xs text-muted-foreground">{userRole}</span>
                                </div>
                            </div>


                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleLogout();
                                }}
                                className="h-8 w-8 hover:bg-destructive/10"
                            >
                                <LogOut className="h-4 w-4 text-muted-foreground hover:text-background transition-colors" />
                                <span className="sr-only">Logout</span>
                            </Button>
                        </div>
                    </div>
                </div>
            </div>

            {/* --- MAIN CONTENT AREA --- */}
            <div className="flex flex-col border rounded-xl m-2 bg-background overflow-hidden h-[calc(100vh-1rem)]">
                {/* HEADER / TOPBAR */}
                <main className="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-4 relative overflow-y-auto">
                    <Sheet>
                        <SheetTrigger asChild>
                            <Button
                                variant="outline"
                                size="icon"
                                className="shrink-0 md:hidden absolute left-4 top-4 z-10"
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

                                <div className="mt-auto">
                                    <div className="flex items-center gap-4 px-2 py-4">
                                        <Avatar className="h-8 w-8">
                                            <AvatarImage src={userAvatar || undefined} alt={userName} />
                                            <AvatarFallback><UserIcon className="h-4 w-4" /></AvatarFallback>
                                        </Avatar>
                                        <div className="flex flex-col">
                                            <span className="font-semibold text-sm">{userName}</span>
                                            <span className="text-xs text-muted-foreground">{userRole}</span>
                                        </div>
                                        <Button variant="ghost" size="icon" onClick={handleLogout} className="ml-auto text-destructive">
                                            <LogOut className="h-5 w-5" />
                                        </Button>
                                    </div>
                                </div>
                            </nav>
                        </SheetContent>
                    </Sheet>

                    <div className="mt-12 md:mt-0 flex-1">
                        <Outlet />
                    </div>
                </main>
            </div>
        </div>
    )
}