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
        <Card className="w-full max-w-md">
            <CardHeader>
                <CardTitle>{t('auth.welcome_back')}</CardTitle>
                <CardDescription>{t('auth.sign_in_subtitle')}</CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    form.handleSubmit();
                }} className="space-y-4">
                    <form.Field
                        name="email"
                        children={(field) => (
                            <div>
                                <Label className="mb-1" htmlFor={field.name}>
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
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <em role="alert" className="text-[0.8rem] font-medium text-destructive mt-1">
                                        {field.state.meta.errors.map(err => typeof err === 'object' ? ((err as any).message ?? JSON.stringify(err)) : String(err)).join(', ')}
                                    </em>
                                )}
                            </div>
                        )}
                    />

                    <form.Field
                        name="password"
                        children={(field) => (
                            <div>
                                <Label className="mb-1" htmlFor={field.name}>
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
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <em role="alert" className="text-[0.8rem] font-medium text-destructive mt-1">
                                        {field.state.meta.errors.map(err => typeof err === 'object' ? ((err as any).message ?? JSON.stringify(err)) : String(err)).join(', ')}
                                    </em>
                                )}
                            </div>
                        )}
                    />

                    {serverError && (
                        <div className="text-destructive text-sm">
                            {serverError}
                        </div>
                    )}
                    <form.Subscribe
                        selector={(state) => [state.canSubmit, state.isSubmitting]}
                        children={([canSubmit, isSubmitting]) => (
                            <div className="pt-2">
                                <Button
                                    type="submit"
                                    className="w-full"
                                    disabled={!canSubmit || isSubmitting || auth.isLoading}
                                >
                                    {isSubmitting || auth.isLoading ? t('auth.signing_in') : t('auth.sign_in')}
                                </Button>
                            </div>
                        )}
                    />
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
    )
}
