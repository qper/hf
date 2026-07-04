import { createFileRoute } from '@tanstack/react-router'
import { RegisterPage } from '@/components/routes/RegisterPage'

export const Route = createFileRoute('/register')({
  component: RegisterPage,
})
