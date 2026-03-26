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
        <div className="min-h-screen flex items-center justify-center p-6 relative">
            <div className="absolute top-4 right-4">
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

