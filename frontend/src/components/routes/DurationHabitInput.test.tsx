import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { DurationHabitInput } from './DurationHabitInput'
import { describe, it, expect } from 'vitest'
import { BoardHabit } from '@/api/board'
import type { ReactNode } from 'react'

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
  return ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}

describe('DurationHabitInput', () => {
  it('renders an input field', () => {
    const { container } = render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const inputs = container.querySelectorAll('input[type="number"]')
    expect(inputs.length).toBeGreaterThan(0)
  })

  it('displays placeholder text when empty', () => {
    const { container } = render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const input = container.querySelector('input[type="number"]') as HTMLInputElement
    expect(input?.placeholder).toBe('0 мин')
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
    const { container } = render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const input = container.querySelector('input[type="number"]') as HTMLInputElement
    expect(input?.inputMode).toBe('numeric')
  })

  it('is disabled when isEditable is false', () => {
    const { container } = render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={false}
      />,
      { wrapper: createWrapper() },
    )
    const input = container.querySelector('input[type="number"]') as HTMLInputElement
    expect(input?.disabled).toBe(true)
  })

  it('displays entry value when provided', async () => {
    const { container } = render(
      <DurationHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
        entryValue={15}
      />,
      { wrapper: createWrapper() },
    )
    const input = container.querySelector('input[type="number"]') as HTMLInputElement
    expect(input?.value).toBe('15')
  })
})
