import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import App from '@/App'

describe('router', () => {
  it('redirects home to board with today\'s date', async () => {
    render(<App />)

    await waitFor(() => {
      expect(window.location.pathname).toContain('/board/')
    })
  })

  it('navigates to login', async () => {
    render(<App />)

    await waitFor(() => {
      expect(window.location.pathname).toContain('/board/')
    })

    const loginLinks = screen.getAllByRole('link', { name: /login/i })
    await userEvent.click(loginLinks[0])

    await waitFor(() => expect(window.location.pathname).toBe('/login'))
  })
})
