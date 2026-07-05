import { createFileRoute } from '@tanstack/react-router'
import { HabitsPage } from '@/components/routes/HabitsPage'

export const Route = createFileRoute('/habits')({
  component: HabitsPage,
})
