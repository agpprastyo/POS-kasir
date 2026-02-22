import { Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { FileQuestion } from 'lucide-react'
import { useTranslation } from 'react-i18next'

export function NotFoundPage() {
    const { t } = useTranslation()

    return (
        <div className="flex h-screen w-full flex-col items-center justify-center gap-4 bg-background px-4 text-center">
            <div className="rounded-full bg-muted p-4">
                <FileQuestion className="h-10 w-10 text-muted-foreground" />
            </div>
            <div className="space-y-2">
                <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl">404</h1>
                <h2 className="text-2xl font-semibold tracking-tight">{t('common.not_found.title')}</h2>
                <p className="max-w-[500px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                    {t('common.not_found.description')}
                </p>
            </div>
            <Button asChild variant="default" size="lg">
                <Link to="/">{t('common.not_found.go_home')}</Link>
            </Button>
        </div>
    )
}
