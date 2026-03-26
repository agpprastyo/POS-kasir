import { Link } from '@tanstack/react-router'
import { Zap, Menu } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'
import { SettingsPanel } from "@/components/SettingsPanel.tsx"
import { ShiftControl } from '@/components/dashboard/ShiftControl'
import { DashboardUserMenu } from '@/components/dashboard/DashboardUserMenu'

interface DashboardSidebarProps {
    t: any
    locale: string
    branding: any
    isSidebarCollapsed: boolean
    setIsSidebarCollapsed: (collapsed: boolean) => void
    filteredMenu: any[]
    user: any
    handleLogout: () => void
}

export function DashboardSidebar({
    t,
    locale,
    branding,
    isSidebarCollapsed,
    setIsSidebarCollapsed,
    filteredMenu,
    user,
    handleLogout
}: DashboardSidebarProps) {
    return (
        <div className="hidden md:block">
            <div className="flex h-full max-h-screen flex-col gap-2">
                <div className={cn("flex h-14 items-center px-4 lg:h-[60px]", isSidebarCollapsed ? "justify-center" : "lg:px-6 justify-between")}>
                    <Link
                        to="/$locale"
                        params={{ locale } as any}
                        className={cn("flex items-center gap-2 font-semibold", isSidebarCollapsed && "hidden")}
                    >
                        {branding?.app_logo ? (
                            <img src={branding.app_logo} alt={t('settings.branding.logo')} className="h-8 w-8 object-contain" />
                        ) : (
                            <Zap className="h-8 w-8" />
                        )}
                        <span className="text-2xl truncate">{branding?.app_name || t('dashboard.brand_name')}</span>
                    </Link>
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setIsSidebarCollapsed(!isSidebarCollapsed)}
                        className="h-8 w-8"
                    >
                        <Menu className="h-5 w-5" />
                        <span className="sr-only">Toggle Sidebar</span>
                    </Button>
                </div>
                <div className="flex-1">
                    <TooltipProvider delayDuration={0}>
                        <nav className="grid items-start px-2 text-sm font-medium lg:px-4 gap-1">
                            {filteredMenu.map((item) => (
                                <Tooltip key={item.to}>
                                    <TooltipTrigger asChild>
                                        <Link
                                            to={item.to}
                                            params={{ locale } as any}
                                            activeOptions={{ exact: item.to === '/$locale' }}
                                            className={cn(
                                                "flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary [&.active]:bg-muted [&.active]:text-primary",
                                                isSidebarCollapsed && "justify-center"
                                            )}
                                        >
                                            <item.icon className={cn("h-4 w-4", isSidebarCollapsed ? "h-5 w-5" : "")} />
                                            {!isSidebarCollapsed && <span>{item.label}</span>}
                                        </Link>
                                    </TooltipTrigger>
                                    {isSidebarCollapsed && (
                                        <TooltipContent side="right">
                                            {item.label}
                                        </TooltipContent>
                                    )}
                                </Tooltip>
                            ))}
                        </nav>
                    </TooltipProvider>
                </div>
                <div className={cn("mt-auto p-4", isSidebarCollapsed && "p-2")}>
                    {!isSidebarCollapsed && (
                        <>
                            <div className="hidden md:block mb-4">
                                <SettingsPanel />
                            </div>
                            <div className="hidden md:block mb-4">
                                <ShiftControl />
                            </div>
                        </>
                    )}

                    <DashboardUserMenu 
                        t={t}
                        user={user}
                        isSidebarCollapsed={isSidebarCollapsed}
                        handleLogout={handleLogout}
                    />
                </div>
            </div>
        </div>
    )
}
