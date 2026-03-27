import { Settings } from 'lucide-react'

interface SettingsHeaderProps {
    t: any
}

export function SettingsHeader({ t }: SettingsHeaderProps) {
    return (
        <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center">
                <Settings className="h-5 w-5 text-primary" />
            </div>
            <div>
                <h1 className="text-2xl font-bold tracking-tight font-heading">{t('settings.title')}</h1>
                <p className="text-sm text-muted-foreground">
                    {t('settings.description')}
                </p>
            </div>
        </div>
    )
}
