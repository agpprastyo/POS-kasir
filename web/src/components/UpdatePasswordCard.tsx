import React, {ReactNode, useState} from "react";
import {useUpdatePasswordMutation} from "@/lib/api/query/auth.ts";
import {POSKasirInternalDtoUpdatePasswordRequest} from "@/lib/api/generated";
import {Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {KeyRound, Loader2, Save} from "lucide-react";
import {Alert, AlertDescription} from "@/components/ui/alert.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";

export function UpdatePasswordCard() {
    const [formData, setFormData] = useState({
        old_password: '',
        new_password: '',
        confirm_password: ''
    })
    const [message, setMessage] = useState<{ type: 'success' | 'error', text: string } | null>(null)

    const mutation = useUpdatePasswordMutation()

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {name, value} = e.target
        setFormData(prev => ({...prev, [name]: value}))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setMessage(null)

        if (formData.new_password !== formData.confirm_password) {
            setMessage({type: 'error', text: 'New passwords do not match.'})
            return
        }

        try {
            const payload: POSKasirInternalDtoUpdatePasswordRequest = {
                old_password: formData.old_password,
                new_password: formData.new_password
            }

            await mutation.mutateAsync(payload)
            setMessage({type: 'success', text: 'Password updated successfully!'})
            setFormData({old_password: '', new_password: '', confirm_password: ''})
        } catch (error: any) {
            const msg = error?.response?.data?.message ?? 'Failed to update password.'
            setMessage({type: 'error', text: msg})
        }
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <KeyRound className="h-5 w-5"/> Password & Security
                </CardTitle>
                <CardDescription>
                    Change your password to keep your account secure.
                </CardDescription>
            </CardHeader>
            <form onSubmit={handleSubmit}>
                <CardContent className="grid gap-4">
                    {message && (
                        <Alert
                            variant={(message.type === 'error' ? 'destructive' : 'default') as "default" | "destructive"}
                            className={message.type === 'success' ? 'border-green-500 text-green-500' : ''}
                        >
                            <AlertDescription>{message.text}</AlertDescription>
                        </Alert> as ReactNode
                    )}

                    <div className="grid gap-2">
                        <Label htmlFor="old_password">Current Password</Label>
                        <Input
                            id="old_password"
                            name="old_password"
                            type="password"
                            value={formData.old_password}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="grid gap-2">
                        <Label htmlFor="new_password">New Password</Label>
                        <Input
                            id="new_password"
                            name="new_password"
                            type="password"
                            value={formData.new_password}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="grid gap-2">
                        <Label htmlFor="confirm_password">Confirm New Password</Label>
                        <Input
                            id="confirm_password"
                            name="confirm_password"
                            type="password"
                            value={formData.confirm_password}
                            onChange={handleChange}
                            required
                        />
                    </div>
                </CardContent>
                <CardFooter className="justify-end border-t bg-muted/20 px-6 py-4">
                    <Button type="submit" disabled={mutation.isPending}>
                        {mutation.isPending ? (
                            <Loader2 className="mr-2 h-4 w-4 animate-spin"/> as ReactNode
                        ) : (
                            <Save className="mr-2 h-4 w-4"/> as ReactNode
                        )}
                        Change Password
                    </Button>
                </CardFooter>
            </form>
        </Card>
    )
}