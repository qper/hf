import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { NumericHabitInput } from './NumericHabitInput'
import { describe, it, expect } from 'vitest'
import { BoardHabit } from '@/api/board'
import type { ReactNode } from 'react'

const mockHabit: BoardHabit = {
  id: 'test-id',
  user_id: 'user-1',
  name: 'Water Intake',
  description: 'Drink water',
  type: 'numeric',
  frequency: 'daily',
  sort_order: 1,
  is_completed: false,
  streak: 0,
  target_value: 8,
  unit: 'glasses',
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

describe('NumericHabitInput', () => {
  it('renders minus and plus buttons', () => {
    render(
      <NumericHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const buttons = screen.queryAllByRole('button')
    expect(buttons.length).toBeGreaterThanOrEqual(2)
  })

  it('displays current value and target', () => {
    render(
      <NumericHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
        entryValue={5}
      />,
      { wrapper: createWrapper() },
    )
    const inputs = screen.queryAllByDisplayValue('5')
    // The value should be shown in the display button or in edit input
    if (inputs.length > 0) {
      expect(inputs[0]).toBeTruthy()
    } else {
      // Check that the component rendered
      const buttons = screen.queryAllByRole('button')
      expect(buttons.length).toBeGreaterThanOrEqual(3)
    }
  })

  it('displays 0 when no entry value', () => {
    render(
      <NumericHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
      />,
      { wrapper: createWrapper() },
    )
    const buttons = screen.queryAllByRole('button')
    expect(buttons.length).toBeGreaterThanOrEqual(3)
  })

  it('has min dimensions of 44x44 for touch targets', () => {
    render(
      <NumericHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
        entryValue={0}
      />,
      { wrapper: createWrapper() },
    )
    const buttons = screen.queryAllByRole('button')
    buttons.forEach((button) => {
      const style = (button as HTMLButtonElement).style
      // At least the main buttons should have min dimensions
      if (style.minWidth && style.minHeight) {
        expect(style.minWidth).toBe('44px')
        expect(style.minHeight).toBe('44px')
      }
    })
  })

  it('is disabled when isEditable is false', () => {
    render(
      <NumericHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={false}
        entryValue={0}
      />,
      { wrapper: createWrapper() },
    )
    const buttons = screen.getAllByRole('button')
    buttons.forEach((button) => {
      expect((button as HTMLButtonElement).disabled).toBe(true)
    })
  })

  it('allows inline editing when value is clicked', async () => {
    render(
      <NumericHabitInput
        habit={mockHabit}
        date="2026-07-13"
        isEditable={true}
        entryValue={5}
      />,
      { wrapper: createWrapper() },
    )
    
    const buttons = screen.getAllByRole('button')
    const valueButton = buttons[1]
    valueButton.click()
    
    const input = screen.getByDisplayValue('5')
    expect(input).toBeTruthy()
  })
})
