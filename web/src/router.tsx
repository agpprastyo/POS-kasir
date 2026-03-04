import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { NotFoundPage } from '@/components/NotFoundPage'
import { ErrorPage } from '@/components/ErrorPage'

import { routeTree } from './routeTree.gen'

import { makeQueryClient } from './lib/queryClient'

import { QueryClientProvider } from '@tanstack/react-query'

export function createRouter() {
  const queryClient = makeQueryClient()

  const router = createTanStackRouter({
    routeTree,
    context: {
      queryClient,
    },
    defaultNotFoundComponent: NotFoundPage,
    defaultErrorComponent: ({ error, reset }) => <ErrorPage error={error} reset={reset} />,
    scrollRestoration: true,
    defaultPreloadStaleTime: 0,
    Wrap: ({ children }) => (
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    ),
  })

  return router
}

export const getRouter = createRouter

declare module '@tanstack/react-router' {
  interface Register {
    router: ReturnType<typeof createRouter>
  }
}
