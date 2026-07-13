import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { DurationHabitInput } from './DurationHabitInput'
import { describe, it, expect } from 'vitest'
import { BoardHabit } from '@/api/board'

const mockHabit: BoardHabit = {
  id: 'test-id',
  user_id: 'user-1',
  name: 'Meditation',
  description: 'Daily meditation',
  type: 'duration',
  frequency: 'daily',
  sort_order: 1,
  is_completed: false,
  streak: 0,
  unit: 'мин',
}

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}

describe('DurationHabitInput', () => {
  it('renders an input field', () => {
    render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const input = screen.getByRole('textbox') as HTMLInputElement
    expect(input).toBeTruthy()
    expect(input.type).toBe('number')
  })

  it('displays placeholder text when empty', () => {
    render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const input = screen.getByRole('textbox') as HTMLInputElement
    expect(input.placeholder).toBe('0 мин')
  })

  it('displays unit label', () => {
    render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    expect(screen.getByText('мин')).toBeTruthy()
  })

  it('has inputmode numeric', () => {
    render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const input = screen.getByRole('textbox') as HTMLInputElement
    expect(input.inputMode).toBe('numeric')
  })

  it('is disabled when isEditable is false', () => {
    render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={false}
      />,
      { wrapper: createWrapper() },
    )
    const input = screen.getByRole('textbox') as HTMLInputElement
    expect(input.disabled).toBe(true)
  })

  it('displays entry value when provided', async () => {
    render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
        entryValue={15}
      />,
      { wrapper: createWrapper() },
    )
    const input = screen.getByRole('textbox') as HTMLInputElement
    expect(input.value).toBe('15')
  })
})
