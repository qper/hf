import { createFileRoute } from '@tanstack/react-router'
import { LoginPage } from '@/components/routes/LoginPage'

export const Route = createFileRoute('/login')({
  component: LoginPage,
})
