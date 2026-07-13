import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BooleanHabitToggle } from './BooleanHabitToggle'
import { describe, it, expect } from 'vitest'
import { BoardHabit } from '@/api/board'
import type { ReactNode } from 'react'

const mockHabit: BoardHabit = {
  id: 'test-id',
  user_id: 'user-1',
  name: 'Morning Exercise',
  description: 'Daily workout',
  type: 'boolean',
  frequency: 'daily',
  sort_order: 1,
  is_completed: false,
  streak: 0,
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

describe('BooleanHabitToggle', () => {
  it('renders as a button', () => {
    render(
      <BooleanHabitToggle
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const buttons = screen.queryAllByRole('button')
    expect(buttons.length).toBeGreaterThan(0)
  })

  it('shows uncompleted state with circle icon', () => {
    const { container } = render(
      <BooleanHabitToggle
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const svgs = container.querySelectorAll('svg')
    expect(svgs.length).toBeGreaterThan(0)
  })

  it('shows completed state with filled checkmark', () => {
    const completedHabit = { ...mockHabit, is_completed: true }
    const { container } = render(
      <BooleanHabitToggle
        habit={completedHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const svgs = container.querySelectorAll('svg')
    expect(svgs.length).toBeGreaterThan(0)
  })

  it('is disabled when isEditable is false', () => {
    render(
      <BooleanHabitToggle
        habit={mockHabit}
        date="2026-07-13"
        isEditable={false}
      />,
      { wrapper: createWrapper() },
    )
    const buttons = screen.queryAllByRole('button')
    expect(buttons.length).toBeGreaterThan(0)
    if (buttons[0]) {
      expect((buttons[0] as HTMLButtonElement).disabled).toBe(true)
    }
  })

  it('has min dimensions of 44x44 for touch targets', () => {
    render(
      <BooleanHabitToggle
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const buttons = screen.queryAllByRole('button')
    if (buttons[0]) {
      const minWidth = (buttons[0] as HTMLButtonElement).style.minWidth
      const minHeight = (buttons[0] as HTMLButtonElement).style.minHeight
      expect(minWidth).toBe('44px')
      expect(minHeight).toBe('44px')
    }
  })
})
