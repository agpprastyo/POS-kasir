import { LogOut, User as UserIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'

interface DashboardUserMenuProps {
    t: any
    user: any
    isSidebarCollapsed?: boolean
    handleLogout: () => void
    isMobile?: boolean
}

export function DashboardUserMenu({
    t,
    user,
    isSidebarCollapsed,
    handleLogout,
    isMobile = false
}: DashboardUserMenuProps) {
    const userRole = user?.role
    const userAvatar = user?.avatar
    const userName = user?.username ?? 'User'

    return (
        <div className={cn(
            "rounded-2xl w-full flex items-center border",
            isMobile ? "flex items-center gap-4 px-2 py-4 border-t" : (
                isSidebarCollapsed 
                    ? "justify-center p-2 h-auto flex-col gap-4" 
                    : "justify-between px-2 gap-2 pl-4 h-12"
            )
        )}>
            <div className="flex items-center gap-4 cursor-default">
                <Avatar className={cn("h-8 w-8", isSidebarCollapsed && "h-10 w-10")}>
                    <AvatarImage src={userAvatar || undefined} alt={userName} />
                    <AvatarFallback><UserIcon className="h-4 w-4" /></AvatarFallback>
                </Avatar>
                {(!isSidebarCollapsed || isMobile) && (
                    <div className="flex flex-col items-start truncate text-sm">
                        <span className="font-semibold">{userName}</span>
                        <span className="text-xs text-muted-foreground">{userRole}</span>
                    </div>
                )}
            </div>

            <Button
                variant="ghost"
                size="icon"
                onClick={(e) => {
                    if (!isMobile) e.stopPropagation();
                    handleLogout();
                }}
                className={cn("h-8 w-8 hover:bg-destructive/10", isMobile && "ml-auto text-destructive")}
                title={t('common.logout')}
            >
                <LogOut className={cn("h-4 w-4", isMobile ? "h-5 w-5" : "text-muted-foreground hover:text-background transition-colors")} />
                <span className="sr-only">{t('common.logout')}</span>
            </Button>
        </div>
    )
}
