import { useState } from 'react'
import {
    Card,
    CardHeader,
    CardTitle,
    CardDescription,
    CardContent,
    CardFooter,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { useForm } from '@tanstack/react-form'
import * as z from 'zod'
import { Zap } from 'lucide-react'

interface LoginFormProps {
    t: any
    auth: any
    onSubmitSuccess: () => void
}

export function LoginForm({ t, auth, onSubmitSuccess }: LoginFormProps) {
    const [serverError, setServerError] = useState<string | null>(null)

    const form = useForm({
        defaultValues: {
            email: '',
            password: '',
        },
        validators: {
            onChange: z.object({
                email: z.string().email(),
                password: z.string().min(1)
            })
        },
        onSubmit: async ({ value }) => {
            setServerError(null)
            try {
                await auth.login({
                    email: value.email,
                    password: value.password,
                })
                onSubmitSuccess()
            } catch (error: any) {
                console.error('Login Failed:', error)
                const msg = error?.response?.data?.message ?? error?.message ?? t('auth.login_failed')
                setServerError(msg)
            }
        }
    })

    const demoAccount = {
        email: 'admin@example.com',
        password: 'passwordrahasia'
    }

    return (
        <Card className="w-full max-w-md border-0 shadow-2xl shadow-primary/10 relative z-10">
            <CardHeader className="text-center pb-2 pt-8">
                <div className="flex justify-center mb-4">
                    <div className="h-14 w-14 rounded-2xl bg-primary flex items-center justify-center shadow-lg shadow-primary/30">
                        <Zap className="h-7 w-7 text-primary-foreground" />
                    </div>
                </div>
                <CardTitle className="text-2xl font-heading font-bold">{t('auth.welcome_back')}</CardTitle>
                <CardDescription className="text-sm">{t('auth.sign_in_subtitle')}</CardDescription>
            </CardHeader>
            <CardContent className="px-8 pb-4">
                <form onSubmit={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    form.handleSubmit();
                }} className="space-y-5">
                    <form.Field
                        name="email"
                        children={(field) => (
                            <div className="space-y-2">
                                <Label htmlFor={field.name} className="text-sm font-medium">
                                    {t('auth.email')}
                                </Label>
                                <Input
                                    id={field.name}
                                    name={field.name}
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) => field.handleChange(e.target.value)}
                                    placeholder={t('auth.email_placeholder')}
                                    type="email"
                                    className="h-12 rounded-xl"
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <em role="alert" className="text-sm font-medium text-destructive">
                                        {field.state.meta.errors.map(err => typeof err === 'object' ? ((err as any).message ?? JSON.stringify(err)) : String(err)).join(', ')}
                                    </em>
                                )}
                            </div>
                        )}
                    />

                    <form.Field
                        name="password"
                        children={(field) => (
                            <div className="space-y-2">
                                <Label htmlFor={field.name} className="text-sm font-medium">
                                    {t('auth.password')}
                                </Label>
                                <Input
                                    id={field.name}
                                    name={field.name}
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) => field.handleChange(e.target.value)}
                                    placeholder={t('auth.password_placeholder')}
                                    type="password"
                                    className="h-12 rounded-xl"
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <em role="alert" className="text-sm font-medium text-destructive">
                                        {field.state.meta.errors.map(err => typeof err === 'object' ? ((err as any).message ?? JSON.stringify(err)) : String(err)).join(', ')}
                                    </em>
                                )}
                            </div>
                        )}
                    />

                    {serverError && (
                        <div className="text-destructive text-sm bg-destructive/10 p-3 rounded-xl text-center">
                            {serverError}
                        </div>
                    )}
                    <form.Subscribe
                        selector={(state) => [state.canSubmit, state.isSubmitting]}
                        children={([canSubmit, isSubmitting]) => (
                            <div className="pt-1">
                                <Button
                                    type="submit"
                                    className="w-full h-12 rounded-xl text-sm font-semibold shadow-lg shadow-primary/25"
                                    disabled={!canSubmit || isSubmitting || auth.isLoading}
                                >
                                    {isSubmitting || auth.isLoading ? t('auth.signing_in') : t('auth.sign_in')}
                                </Button>
                            </div>
                        )}
                    />
                </form>
            </CardContent>

            <CardFooter className="px-8 pb-8">
                <div className="w-full text-center text-sm text-muted-foreground bg-muted/50 rounded-xl p-4">
                    <p className="font-medium mb-1">{t('auth.demo_account')}:</p>
                    <p>{t('auth.email')}: <code className="text-sm bg-background px-1.5 py-0.5 rounded-md">{demoAccount.email}</code></p>
                    <p>{t('auth.password')}: <code className="text-sm bg-background px-1.5 py-0.5 rounded-md">{demoAccount.password}</code></p>
                </div>
            </CardFooter>
        </Card>
    )
}
