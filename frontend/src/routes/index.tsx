import { createFileRoute, Navigate } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: IndexRoute,
})

function IndexRoute() {
  const today = new Date().toISOString().split('T')[0]
  return <Navigate to={`/board/${today}`} replace />
}
