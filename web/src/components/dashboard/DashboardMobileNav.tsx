import { Link } from '@tanstack/react-router'
import { Zap, Menu } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'
import { SettingsPanel } from "@/components/SettingsPanel.tsx"
import { ShiftControl } from '@/components/dashboard/ShiftControl'
import { DashboardUserMenu } from './DashboardUserMenu'

interface DashboardMobileNavProps {
    t: any
    locale: string
    branding: any
    filteredMenu: any[]
    user: any
    handleLogout: () => void
}

export function DashboardMobileNav({
    t,
    locale,
    branding,
    filteredMenu,
    user,
    handleLogout
}: DashboardMobileNavProps) {
    return (
        <Sheet>
            <SheetTrigger asChild>
                <Button
                    variant="outline"
                    size="icon"
                    className="shrink-0 md:hidden absolute left-4 top-4 z-10"
                >
                    <Menu className="h-5 w-5" />
                    <span className="sr-only">{t('dashboard.toggle_nav')}</span>
                </Button>
            </SheetTrigger>
            <SheetContent side="left" className="flex flex-col">
                <nav className="grid gap-2 text-lg font-medium">
                    <Link
                        to="/$locale"
                        params={{ locale } as any}
                        className="flex items-center gap-2 text-lg font-semibold mb-4"
                    >
                        {branding?.app_logo ? (
                            <img src={branding.app_logo} alt={t('settings.branding.logo')} className="h-6 w-6 object-contain" />
                        ) : (
                            <Zap className="h-6 w-6" />
                        )}
                        <span className="sr-only">{branding?.app_name || t('dashboard.brand_name')}</span>
                    </Link>
                    <div className="mb-4">
                        <SettingsPanel />
                    </div>
                    <div className="mb-4">
                        <ShiftControl />
                    </div>
                    {/* Mobile Menu Filtered */}
                    {filteredMenu.map((item) => (
                        <Link
                            key={item.to}
                            to={item.to}
                            params={{ locale } as any}
                            activeOptions={{ exact: item.to === '/$locale' }}
                            className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground [&.active]:bg-muted [&.active]:text-foreground"
                        >
                            <item.icon className="h-5 w-5" />
                            {item.label}
                        </Link>
                    ))}

                    <div className="mt-auto">
                        <DashboardUserMenu 
                            t={t}
                            user={user}
                            handleLogout={handleLogout}
                            isMobile={true}
                        />
                    </div>
                </nav>
            </SheetContent>
        </Sheet>
    )
}
