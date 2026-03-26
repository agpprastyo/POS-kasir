interface SettingsHeaderProps {
    t: any
}

export function SettingsHeader({ t }: SettingsHeaderProps) {
    return (
        <div>
            <h1 className="text-3xl font-bold tracking-tight">{t('settings.title')}</h1>
            <p className="text-muted-foreground">
                {t('settings.description')}
            </p>
        </div>
    )
}
