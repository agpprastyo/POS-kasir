import React, { useState } from 'react'
import {createFileRoute, RegisteredRouter, useRouter} from '@tanstack/react-router'
import { redirect } from '@tanstack/react-router'


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



import {meQueryOptions} from "@/lib/api/query/auth.ts";
import {useAuth} from "@/lib/auth/AuthContext.tsx";
import {queryClient} from "@/lib/queryClient.ts";
import {FileRouteByToPath} from "@tanstack/router-core/src/routeInfo.ts";

export const Route = createFileRoute('/login' as FileRouteByToPath<any, any>)({
    ssr: false,
    loader: async () => {
        try {
            const me = await queryClient.ensureQueryData(meQueryOptions())
            if (me) {
                throw redirect({ to: '/' } as RegisteredRouter)
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
            await router.navigate({ to: '/', replace: true })

        } catch (error: any) {
            console.error('Login Failed:', error)
            const msg = error?.response?.data?.message ?? error?.message ?? 'Login failed. Please check your credentials.'
            setServerError(msg)
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center p-6">
            <Card className="w-full max-w-md">
                <CardHeader>
                    <CardTitle>Welcome back</CardTitle>
                    <CardDescription>Sign in to your account to continue.</CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <Label className="mb-1" htmlFor="email">
                                Email
                            </Label>
                            <Input
                                id="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                                placeholder="you@example.com"
                                type="email"
                                required
                            />
                        </div>

                        <div>
                            <Label className="mb-1" htmlFor="password">
                                Password
                            </Label>
                            <Input
                                id="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                placeholder="Your password"
                                type="password"
                                required
                            />
                        </div>

                        {serverError && (
                            <div className="text-destructive text-sm">
                                {serverError}
                            </div>
                        )}

                        <div className="flex items-center justify-between">
                            <a
                                href="/forgot"
                                className="text-sm text-muted-foreground hover:underline"
                            >
                                Forgot password?
                            </a>
                        </div>

                        <div className="pt-2">
                            <Button
                                type="submit"
                                className="w-full"
                                disabled={isSubmitting}
                            >
                                {isSubmitting ? 'Signing in...' : 'Sign in'}
                            </Button>
                        </div>
                    </form>
                </CardContent>
                <CardFooter>
                    <div className="text-sm text-muted-foreground">
                        Don&apos;t have an account?{' '}
                        <a
                            href="/register"
                            className="text-primary hover:underline"
                        >
                            Register
                        </a>
                    </div>
                </CardFooter>
            </Card>
        </div>
    )
}
