import { createFileRoute } from '@tanstack/react-router'
import { BoardPageWithRoute } from '@/components/routes/BoardPage'

export const Route = createFileRoute('/board/$date')({
  component: BoardPageWithRoute,
})
