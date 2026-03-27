import { createFileRoute, useRouter, useParams } from '@tanstack/react-router'
import { redirect } from '@tanstack/react-router'
import { meQueryOptions } from "@/lib/api/query/auth.ts";
import { useAuth } from "@/context/AuthContext";
import { useTranslation } from 'react-i18next'
import { SettingsPanel } from "@/components/SettingsPanel.tsx";
import { LoginForm } from "@/components/auth/LoginForm"

export const Route = createFileRoute('/$locale/login')({
    ssr: false,
    loader: async ({ context: { queryClient }, params }) => {
        try {
            const me = await queryClient.ensureQueryData(meQueryOptions())
            if (me) {
                throw redirect({
                    to: '/$locale',
                    params: { locale: params.locale }
                })
            }
        } catch (error: any) {
            const status = error?.response?.status ?? error?.status ?? error?.cause?.status
            if (status === 401) {
                return
            }
        }
    },
    component: LoginPage,
})

function LoginPage() {
    const { locale } = useParams({ from: '/$locale/login' })
    const { t } = useTranslation()
    const auth = useAuth()
    const router = useRouter()

    return (
        <div className="min-h-screen flex items-center justify-center p-6 relative bg-linear-to-br from-primary/5 via-background to-amber/5">
            {/* Decorative elements */}
            <div className="absolute inset-0 overflow-hidden pointer-events-none">
                <div className="absolute -top-24 -left-24 w-96 h-96 rounded-full bg-primary/5 blur-3xl" />
                <div className="absolute -bottom-24 -right-24 w-96 h-96 rounded-full bg-amber/5 blur-3xl" />
                <div className="absolute top-1/4 right-1/4 w-64 h-64 rounded-full bg-primary/3 blur-2xl" />
            </div>

            <div className="absolute top-4 right-4 z-10">
                <SettingsPanel />
            </div>
            <LoginForm 
                t={t}
                auth={auth}
                onSubmitSuccess={async () => {
                    await router.invalidate()
                    await router.navigate({
                        to: '/$locale',
                        params: { locale },
                        replace: true
                    })
                }}
            />
        </div>
    )
}

