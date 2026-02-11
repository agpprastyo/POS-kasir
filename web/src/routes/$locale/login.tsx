import React, { useState } from 'react'
import { createFileRoute, useRouter, useParams } from '@tanstack/react-router'
import { redirect } from '@tanstack/react-router'
import {
    Card,
    CardHeader,
    CardTitle,
    CardDescription,
    CardContent, CardFooter,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { meQueryOptions } from "@/lib/api/query/auth.ts";
import { useAuth } from "@/context/AuthContext";
import { queryClient } from "@/lib/queryClient.ts";
import { useTranslation } from 'react-i18next'
import { SettingsPanel } from "@/components/SettingsPanel.tsx";

export const Route = createFileRoute('/$locale/login')({
    ssr: false,
    loader: async ({ params }) => {
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

    const [serverError, setServerError] = useState<string | null>(null)

    const [formData, setFormData] = useState({
        email: '',
        password: '',
    })

    const isSubmitting = auth.isLoading

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target
        setFormData((prev) => ({ ...prev, [name]: value }))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setServerError(null)
        try {
            await auth.login({
                email: formData.email,
                password: formData.password,
            })
            await router.invalidate()
            await router.navigate({
                to: '/$locale',
                params: { locale },
                replace: true
            })

        } catch (error: any) {
            console.error('Login Failed:', error)
            const msg = error?.response?.data?.message ?? error?.message ?? t('auth.login_failed')
            setServerError(msg)
        }
    }

    const demoAccount = {
        email: 'admin@example.com',
        password: 'passwordrahasia'
    }

    return (
        <div className="min-h-screen flex items-center justify-center p-6 relative">
            <div className="absolute top-4 right-4">
                <SettingsPanel />
            </div>
            <Card className="w-full max-w-md">
                <CardHeader>
                    <CardTitle>{t('auth.welcome_back')}</CardTitle>
                    <CardDescription>{t('auth.sign_in_subtitle')}</CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <Label className="mb-1" htmlFor="email">
                                {t('auth.email')}
                            </Label>
                            <Input
                                id="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                                placeholder={t('auth.email_placeholder')}
                                type="email"
                                required
                            />
                        </div>

                        <div>
                            <Label className="mb-1" htmlFor="password">
                                {t('auth.password')}
                            </Label>
                            <Input
                                id="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                placeholder={t('auth.password_placeholder')}
                                type="password"
                                required
                            />
                        </div>

                        {serverError && (
                            <div className="text-destructive text-sm">
                                {serverError}
                            </div>
                        )}
                        <div className="pt-2">
                            <Button
                                type="submit"
                                className="w-full"
                                disabled={isSubmitting}
                            >
                                {isSubmitting ? t('auth.signing_in') : t('auth.sign_in')}
                            </Button>
                        </div>
                    </form>
                </CardContent>

                <CardFooter>
                    <div className="w-full text-center text-sm text-muted-foreground">
                        <p>{t('auth.demo_account')}:</p>
                        <p>{t('auth.email')}: <code>{demoAccount.email}</code></p>
                        <p>{t('auth.password')}: <code>{demoAccount.password}</code></p>
                    </div>
                </CardFooter>
            </Card>
        </div>
    )
}

