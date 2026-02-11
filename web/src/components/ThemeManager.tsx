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
        // Fallback: If dark mode color is not set, stick to the main theme color (or vice versa? user expectation implies separate control)
        // If theme_color is set but theme_color_dark is empty, maybe we should use theme_color for both?
        // Or let it be empty/default?
        // For now: prioritize specific setting, fallback to light mode color if dark is missing (to maintain previous behavior)

        let activeColor = branding.theme_color

        if (isDark && branding.theme_color_dark) {
            activeColor = branding.theme_color_dark
        } else if (isDark && !branding.theme_color_dark) {
            // Fallback for dark mode if only light mode color is set?
            // User requested explicit separate control. If they don't set it, maybe default behavior is safer.
            // But let's assume they want the main color if dark is not specified.
            activeColor = branding.theme_color
        }

        if (!activeColor) return

        // Update primary color
        root.style.setProperty('--primary', activeColor)
        root.style.setProperty('--sidebar-primary', activeColor)
        root.style.setProperty('--ring', activeColor)

    }, [branding, resolvedTheme])

    return null
}
