import { useEffect } from 'react'
import { useBrandingSettingsQuery } from '@/lib/api/query/settings'
import { useTheme } from 'next-themes'

export function ThemeManager() {
    const { data: branding } = useBrandingSettingsQuery()
    const { resolvedTheme } = useTheme()

    useEffect(() => {
        if (!branding) return

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
