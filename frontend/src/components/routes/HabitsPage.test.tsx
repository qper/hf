import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { HabitsPage } from './HabitsPage'
import { reorderHabits } from './habitReorder'

describe('HabitsPage', () => {
  it('hides an archived habit from the active list', () => {
    render(<HabitsPage />)

    const archiveButton = screen.getAllByRole('button', { name: /archive/i })[0]

    fireEvent.click(archiveButton)

    expect(document.body.textContent).not.toContain('Read 20 min')
  })

  it('asks for confirmation before deleting a habit', () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true)
    render(<HabitsPage />)

    fireEvent.click(screen.getAllByRole('button', { name: /more/i })[0])
    fireEvent.click(screen.getAllByRole('menuitem', { name: /delete/i })[0])

    expect(confirmSpy).toHaveBeenCalled()
  })

  it('moves a habit to the front of the list', () => {
    const habits = [
      { id: '1', name: 'First' },
      { id: '2', name: 'Second' },
      { id: '3', name: 'Third' },
    ]

    const reordered = reorderHabits(habits, '3', '1')

    expect(reordered.map((habit) => habit.id)).toEqual(['3', '1', '2'])
  })
})
