import { createRouter } from '@tanstack/react-router'
import { NotFoundPage } from '@/components/NotFoundPage'
import { ErrorPage } from '@/components/ErrorPage'

// Import the generated route tree
import { routeTree } from './routeTree.gen'

// Create a new router instance
export const getRouter = () => {
  return createRouter({
    routeTree,
    defaultNotFoundComponent: NotFoundPage,
    defaultErrorComponent: ({ error, reset }) => <ErrorPage error={error} reset={reset} />,
    scrollRestoration: true,
    defaultPreloadStaleTime: 0,
  })
}
