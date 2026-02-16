import { useEffect } from 'react'
import { useBrandingSettingsQuery } from '@/lib/api/query/settings'
import { useTheme } from 'next-themes'

export function ThemeManager() {
    const { data: branding } = useBrandingSettingsQuery()
    const { resolvedTheme } = useTheme()

    useEffect(() => {
        if (!branding) return

        const root = document.documentElement
        const isDark = resolvedTheme === 'dark'

        // Determine the color to use
        let activeColor = branding.theme_color

        if (isDark && branding.theme_color_dark) {
            activeColor = branding.theme_color_dark
        } else if (isDark && !branding.theme_color_dark) {
            activeColor = branding.theme_color
        }

        if (activeColor) {
            // Update primary color
            root.style.setProperty('--primary', activeColor)
            root.style.setProperty('--sidebar-primary', activeColor)
            root.style.setProperty('--ring', activeColor)
        }

        if (branding.app_logo) {
            const link = (document.querySelector("link[rel*='icon']") || document.createElement('link')) as HTMLLinkElement;
            link.type = 'image/x-icon';
            link.rel = 'shortcut icon';
            link.href = branding.app_logo;
            document.getElementsByTagName('head')[0].appendChild(link);
        }

    }, [branding, resolvedTheme])

    return null
}
