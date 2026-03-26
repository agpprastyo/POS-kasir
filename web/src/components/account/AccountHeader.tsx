interface AccountHeaderProps {
    t: any
}

export function AccountHeader({ t }: AccountHeaderProps) {
    return (
        <div>
            <h1 className="text-3xl font-bold tracking-tight">{t('account.title')}</h1>
            <p className="text-muted-foreground">
                {t('account.subtitle')}
            </p>
        </div>
    )
}
