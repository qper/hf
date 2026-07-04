import { createFileRoute } from '@tanstack/react-router'
import { HomePage } from '@/components/routes/HomePage'

export const Route = createFileRoute('/')({
  component: HomePage,
})
