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

export type Entry = {
  id: string
  habit_id: string
  entry_date: string
  completed: boolean
  value?: number | null
  note?: string | null
  created_at: string
  updated_at: string
}

export type CreateEntryRequest = {
  habit_id: string
  date: string
  completed?: boolean
  value?: number | null
  note?: string | null
}

export async function getBoard(date: string): Promise<Board> {
  const res = await authFetch(`/api/v1/board/${date}`)
  if (!res.ok) {
    throw new Error(`Failed to fetch board: ${res.statusText}`)
  }
  return res.json() as Promise<Board>
}

export async function createEntry(req: CreateEntryRequest): Promise<Entry> {
  const res = await authFetch('/api/v1/entries', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    throw new Error(`Failed to create entry: ${res.statusText}`)
  }
  return res.json() as Promise<Entry>
}

export async function updateEntry(
  entryId: string,
  req: Partial<CreateEntryRequest>,
): Promise<Entry> {
  const res = await authFetch(`/api/v1/entries/${entryId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    throw new Error(`Failed to update entry: ${res.statusText}`)
  }
  return res.json() as Promise<Entry>
}

export async function deleteEntry(entryId: string): Promise<Entry> {
  const res = await authFetch(`/api/v1/entries/${entryId}`, {
    method: 'DELETE',
  })
  if (!res.ok) {
    throw new Error(`Failed to delete entry: ${res.statusText}`)
  }
  return res.json() as Promise<Entry>
}
