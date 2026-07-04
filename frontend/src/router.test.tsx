import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import App from '@/App'

describe('router', () => {
  it('navigates to login', async () => {
    render(<App />)

    await waitFor(() => expect(screen.getByText(/home route/i)).toBeTruthy())
    expect(window.location.pathname).toBe('/')

    await userEvent.click(screen.getByRole('link', { name: /login/i }))

    await waitFor(() => expect(window.location.pathname).toBe('/login'))
  })
})
