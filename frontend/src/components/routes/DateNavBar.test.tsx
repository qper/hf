import { render, screen } from '@testing-library/react'
import { DateNavBar } from './DateNavBar'
import { I18nextProvider } from 'react-i18next'
import i18n from '@/i18n'
import { describe, it, expect, vi } from 'vitest'

describe('DateNavBar', () => {
  const today = new Date().toISOString().split('T')[0]
  const yesterday = new Date(new Date().setDate(new Date().getDate() - 1))
    .toISOString()
    .split('T')[0]

  const renderDateNavBar = (date: string, onDateChange?: (date: string) => void) => {
    const mockOnDateChange = onDateChange || vi.fn()
    return render(
      <I18nextProvider i18n={i18n}>
        <DateNavBar
          date={date}
          onDateChange={mockOnDateChange}
          progress={{ done: 5, total: 7 }}
        />
      </I18nextProvider>,
    )
  }

  it('should display the progress counter', () => {
    renderDateNavBar(today)
    const progressTexts = screen.queryAllByText('5 / 7')
    expect(progressTexts.length).toBeGreaterThan(0)
  })

  it('should disable forward arrow when on today', () => {
    renderDateNavBar(today)
    const buttons = screen.getAllByRole('button')
    const forwardButton = buttons[buttons.length - 1]
    expect((forwardButton as HTMLButtonElement).disabled).toBe(true)
  })

  it('should enable forward arrow when not on today', () => {
    renderDateNavBar(yesterday)
    const buttons = screen.getAllByRole('button')
    const forwardButton = buttons[buttons.length - 1]
    expect((forwardButton as HTMLButtonElement).disabled).toBe(false)
  })

  it('should render two navigation buttons', () => {
    renderDateNavBar(today)
    const buttons = screen.getAllByRole('button')
    expect(buttons.length).toBeGreaterThanOrEqual(2)
  })
})
