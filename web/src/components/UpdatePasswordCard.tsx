import { ReactNode, useState } from "react";
import { useUpdatePasswordMutation } from "@/lib/api/query/auth.ts";
import { InternalUserUpdatePasswordRequest } from "@/lib/api/generated";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { KeyRound, Loader2, Save } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert.tsx";
import { Label } from "@/components/ui/label.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useTranslation } from 'react-i18next';
import { useForm } from '@tanstack/react-form';
import * as z from 'zod';

export function UpdatePasswordCard() {
    const { t } = useTranslation();
    const [message, setMessage] = useState<{ type: 'success' | 'error', text: string } | null>(null)

    const mutation = useUpdatePasswordMutation()

    const form = useForm({
        defaultValues: {
            old_password: '',
            new_password: '',
            confirm_password: ''
        },
        validators: {
            onChange: z.object({
                old_password: z.string().min(1),
                new_password: z.string().min(6),
                confirm_password: z.string().min(1)
            }).refine((data) => data.new_password === data.confirm_password, {
                message: t('account.password.error_match'),
                path: ["confirm_password"],
            })
        },
        onSubmit: async ({ value }) => {
            setMessage(null)

            try {
                const payload: InternalUserUpdatePasswordRequest = {
                    old_password: value.old_password,
                    new_password: value.new_password
                }

                await mutation.mutateAsync(payload)
                setMessage({ type: 'success', text: t('account.password.success') })
                form.reset()
            } catch (error: any) {
                const msg = error?.response?.data?.message ?? t('account.password.error_fail')
                setMessage({ type: 'error', text: msg })
            }
        }
    })

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <KeyRound className="h-5 w-5" /> {t('account.password.title')}
                </CardTitle>
                <CardDescription>
                    {t('account.password.description')}
                </CardDescription>
            </CardHeader>
            <form onSubmit={(e) => {
                e.preventDefault();
                e.stopPropagation();
                form.handleSubmit();
            }}>
                <CardContent className="grid gap-4">
                    {message && (
                        <Alert
                            variant={(message.type === 'error' ? 'destructive' : 'default') as "default" | "destructive"}
                            className={message.type === 'success' ? 'border-primary text-primary' : ''}
                        >
                            <AlertDescription>{message.text}</AlertDescription>
                        </Alert> as ReactNode
                    )}

                    <form.Field
                        name="old_password"
                        children={(field) => (
                            <div className="grid gap-2">
                                <Label htmlFor={field.name}>{t('account.password.current_password')}</Label>
                                <Input
                                    id={field.name}
                                    name={field.name}
                                    type="password"
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) => field.handleChange(e.target.value)}
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <em role="alert" className="text-[0.8rem] font-medium text-destructive">
                                        {field.state.meta.errors.join(', ')}
                                    </em>
                                )}
                            </div>
                        )}
                    />

                    <form.Field
                        name="new_password"
                        children={(field) => (
                            <div className="grid gap-2">
                                <Label htmlFor={field.name}>{t('account.password.new_password')}</Label>
                                <Input
                                    id={field.name}
                                    name={field.name}
                                    type="password"
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) => field.handleChange(e.target.value)}
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <em role="alert" className="text-[0.8rem] font-medium text-destructive">
                                        {field.state.meta.errors.join(', ')}
                                    </em>
                                )}
                            </div>
                        )}
                    />

                    <form.Field
                        name="confirm_password"
                        children={(field) => (
                            <div className="grid gap-2">
                                <Label htmlFor={field.name}>{t('account.password.confirm_password')}</Label>
                                <Input
                                    id={field.name}
                                    name={field.name}
                                    type="password"
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) => field.handleChange(e.target.value)}
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <em role="alert" className="text-[0.8rem] font-medium text-destructive">
                                        {field.state.meta.errors.join(', ')}
                                    </em>
                                )}
                            </div>
                        )}
                    />
                </CardContent>
                <CardFooter className="justify-end border-t bg-muted/20 px-6 py-4">
                    <form.Subscribe
                        selector={(state) => [state.canSubmit, state.isSubmitting]}
                        children={([canSubmit, isSubmitting]) => (
                            <Button type="submit" disabled={!canSubmit || isSubmitting || mutation.isPending}>
                                {(isSubmitting || mutation.isPending) ? (
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" /> as ReactNode
                                ) : (
                                    <Save className="mr-2 h-4 w-4" /> as ReactNode
                                )}
                                {t('account.password.button')}
                            </Button>
                        )}
                    />
                </CardFooter>
            </form>
        </Card>
    )
}