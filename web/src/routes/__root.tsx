import {HeadContent, Scripts, createRootRoute} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'

import appCss from '../styles.css?url'
import React from 'react'

import { QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider } from '@/lib/auth/AuthContext'
import {queryClient} from "@/lib/queryClient.ts";
import {Toaster} from "@/components/ui/sonner.tsx";



export const Route = createRootRoute({

    head: () => ({
        meta: [
            { charSet: 'utf-8' },
            { name: 'viewport', content: 'width=device-width, initial-scale=1' },
            { title: 'POS Kasir' },
        ],
        links: [{ rel: 'stylesheet', href: appCss }],
    }),
    shellComponent: RootDocument,
    notFoundComponent: NotFound,
    errorComponent: RootError,

} as any)


function RootDocument({ children }: any) {
    return (
        <html lang="en">
        <head>
            <HeadContent />
        </head>
        <body>
        <QueryClientProvider client={queryClient}>
            <AuthProvider>
                {children}
            </AuthProvider>
            <TanStackDevtools
                config={{ position: 'bottom-right' }}
                plugins={[
                    {
                        name: 'Tanstack Router',
                        render: () => <TanStackRouterDevtoolsPanel />,
                    },
                ]}
            />
        </QueryClientProvider>
        <Scripts />
        <Toaster />
        </body>
        </html>
    )
}


function RootError({ error }: { error: any }) {
    const message =
        error?.message ??
        error?.response?.data?.message ??
        'Something went wrong'

    return (
        <div className="min-h-screen flex items-center justify-center p-6 ">
            <div className="text-center max-w-md">
                <h1 className="text-4xl font-extrabold text-white mb-2">Oops…</h1>
                <p className="text-gray-300 mb-4">
                    An unexpected error occurred while loading this page.
                </p>
                <pre className="text-sm text-red-200 bg-black/40 rounded-md p-3 overflow-auto">
          {String(message)}
        </pre>
                <a
                    href="/"
                    className="mt-6 inline-block px-6 py-3 bg-primary text-primary-foreground rounded-md hover:opacity-95"
                >
                    Go home
                </a>
            </div>
        </div>
    )
}

function NotFound() {
    return (
        <div className="min-h-screen flex items-center justify-center p-6">
            <div className="w-full max-w-md p-8 text-center">
                <h1 className="text-6xl font-extrabold text-foreground">404</h1>
                <p className="mt-4 text-sm text-muted-foreground">Page not found</p>
                <a
                    href="/"
                    className="mt-6 inline-flex items-center justify-center rounded-md bg-primary px-6 py-3 text-sm font-semibold text-primary-foreground shadow-sm hover:opacity-95"
                >
                    Go home
                </a>
            </div>
        </div>
    )
}