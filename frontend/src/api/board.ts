import { authFetch } from './auth'

export type BoardProgress = {
  done: number
  total: number
}

export type BoardHabit = {
  id: string
  user_id: string
  category_id?: string | null
  name: string
  description?: string | null
  color?: string | null
  type: 'boolean' | 'numeric' | 'duration'
  frequency: string
  target_value?: number | null
  unit?: string | null
  sort_order: number
  is_completed: boolean
  streak: number
}

export type Board = {
  date: string
  is_editable: boolean
  progress: BoardProgress
  habits: BoardHabit[]
}

export async function getBoard(date: string): Promise<Board> {
  const res = await authFetch(`/api/v1/board/${date}`)
  if (!res.ok) {
    throw new Error(`Failed to fetch board: ${res.statusText}`)
  }
  return res.json() as Promise<Board>
}
