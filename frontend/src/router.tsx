import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import {
  RouterProvider,
  RouterContextProvider,
  createBrowserHistory,
  createRouter,
} from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { routeTree } from './routeTree.gen'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 2,
    },
  },
})

const router = createRouter({
  routeTree,
  history: createBrowserHistory(),
  defaultPreload: 'intent',
})

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

export function AppRouter() {
  const isDevelopment = import.meta.env.MODE === 'development'

  return (
    <QueryClientProvider client={queryClient}>
      <RouterContextProvider router={router}>
        <RouterProvider router={router} />
        {isDevelopment ? <ReactQueryDevtools initialIsOpen={false} /> : null}
        {isDevelopment ? (
          <TanStackRouterDevtools position="bottom-right" />
        ) : null}
      </RouterContextProvider>
    </QueryClientProvider>
  )
}
