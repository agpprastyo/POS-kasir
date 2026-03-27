import { useState, useEffect } from 'react'
import { createFileRoute, Outlet, redirect, useRouter, useParams } from '@tanstack/react-router'
import { OpenShiftModal } from "@/components/modals/OpenShiftModal"
import { CloseShiftModal } from "@/components/modals/CloseShiftModal"
import { cn } from '@/lib/utils'
import { useAuth } from '@/context/AuthContext'
import { meQueryOptions } from '@/lib/api/query/auth'
import { useTranslation } from 'react-i18next'
import { useBrandingSettingsQuery } from '@/lib/api/query/settings'
import { useNavigationMenu } from '@/hooks/useNavigationMenu'
import { DashboardSidebar } from "@/components/dashboard/DashboardSidebar"
import { DashboardMobileNav } from "@/components/dashboard/DashboardMobileNav"
import { useAppWebSocket } from '@/hooks/useAppWebSocket'
import { useQueryClient } from '@tanstack/react-query'

export const Route = createFileRoute('/$locale/_dashboard')({
    loader: async ({ context: { queryClient }, params }) => {
        try {
            return await queryClient.ensureQueryData(meQueryOptions())
        } catch (error: any) {
            const status = error?.response?.status ?? error?.status
            if (status === 401) {
                throw redirect({
                    to: '/$locale/login',

                    params: { locale: params.locale },
                })
            }
            throw error
        }
    },
    component: DashboardLayout,
})

function DashboardLayout() {
    const { locale } = useParams({ from: '/$locale/_dashboard' })
    const { t } = useTranslation()
    const auth = useAuth()
    const router = useRouter()
    const [isLoggingOut, setIsLoggingOut] = useState(false)
    const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false)
    const { data: branding } = useBrandingSettingsQuery()
    const queryClient = useQueryClient()

    useAppWebSocket((event) => {
        if (event.type === 'ORDER_CREATED' || event.type === 'ORDER_UPDATED') {
            queryClient.invalidateQueries({ queryKey: ['orders'] })
        }
    })

    useEffect(() => {
        if (branding?.app_name) {
            document.title = branding.app_name
        }
    }, [branding?.app_name])

    const handleLogout = async () => {
        if (isLoggingOut) return
        setIsLoggingOut(true)

        try {
            await auth.logout()
            await router.navigate({
                to: '/$locale/login',
                params: { locale },
                replace: true
            })
        } catch (error) {
            console.error("Logout UI error:", error)
        } finally {
            setIsLoggingOut(false)
        }
    }

    const { filteredMenu } = useNavigationMenu(auth.user?.role)

    return (
        <div className={cn(
            "grid h-screen w-full overflow-hidden transition-all duration-300",
            isSidebarCollapsed ? "md:grid-cols-[80px_1fr]" : "md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]"
        )}>
            <DashboardSidebar
                t={t}
                locale={locale}
                branding={branding}
                isSidebarCollapsed={isSidebarCollapsed}
                setIsSidebarCollapsed={setIsSidebarCollapsed}
                filteredMenu={filteredMenu}
                user={auth.user}
                handleLogout={handleLogout}
            />

            <div className='min-w-0 overflow-hidden'>
                <div className="flex flex-col bg-card/50 backdrop-blur-sm overflow-hidden h-[calc(100vh)] shadow-sm">
                    <main className="flex flex-1 flex-col gap-4 p-4 sm:p-5 lg:gap-6 lg:p-6 relative overflow-y-auto overflow-x-hidden min-w-0">
                        <DashboardMobileNav
                            t={t}
                            locale={locale}
                            branding={branding}
                            filteredMenu={filteredMenu}
                            user={auth.user}
                            handleLogout={handleLogout}
                        />

                        <div className="mt-12 md:mt-0 flex-1 min-w-0">
                            <Outlet />
                        </div>
                    </main>
                </div>
                <OpenShiftModal />
                <CloseShiftModal />
            </div>
        </div>
    )
}
