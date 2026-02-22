import { useRouter } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { AlertCircle, RefreshCw, Home } from 'lucide-react'
import { useTranslation } from 'react-i18next'

export function ErrorPage({ error, reset }: { error?: Error; reset?: () => void }) {
    const router = useRouter()
    const { t } = useTranslation()

    return (
        <div className="flex h-screen w-full flex-col items-center justify-center gap-6 bg-background px-4 text-center">
            <div className="rounded-full bg-destructive/10 p-4">
                <AlertCircle className="h-12 w-12 text-destructive" />
            </div>
            <div className="space-y-2">
                <h1 className="text-4xl font-bold tracking-tighter text-destructive sm:text-5xl">{t('common.error.title')}</h1>
                <h2 className="text-2xl font-semibold tracking-tight">{t('common.error.subtitle')}</h2>
                <p className="max-w-[600px] text-muted-foreground">
                    {error?.message || t('common.error.default_message')}
                </p>
            </div>

            <div className="flex gap-4">
                <Button
                    variant="outline"
                    size="lg"
                    onClick={() => {
                        // Attempt to recover by invalidating router context or just reloading
                        router.invalidate()
                        // If provided a reset function (like from ErrorBoundary), call it
                        if (reset) reset()
                    }}
                    className="gap-2"
                >
                    <RefreshCw className="h-4 w-4" />
                    {t('common.error.try_again')}
                </Button>
                <Button asChild variant="default" size="lg" className="gap-2">
                    <a href="/">
                        <Home className="h-4 w-4" />
                        {t('common.error.go_home')}
                    </a>
                </Button>
            </div>
        </div>
    )
}
