import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { HabitsPage } from './HabitsPage'

describe('HabitsPage', () => {
  it('hides an archived habit from the active list', () => {
    render(<HabitsPage />)

    const archiveButton = screen.getByRole('button', { name: /archive/i })

    fireEvent.click(archiveButton)

    expect(document.body.textContent).not.toContain('Read 20 min')
  })

  it('asks for confirmation before deleting a habit', () => {
    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true)
    render(<HabitsPage />)

    fireEvent.click(screen.getByRole('button', { name: /more/i }))
    fireEvent.click(screen.getByRole('menuitem', { name: /delete/i }))

    expect(confirmSpy).toHaveBeenCalled()
  })
})
