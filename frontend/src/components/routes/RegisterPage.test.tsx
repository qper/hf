import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import '@/i18n'
import * as auth from '@/api/auth'
import { RegisterPage } from './RegisterPage'

const navigateMock = vi.fn()

vi.mock('@tanstack/react-router', () => ({
  useNavigate: () => navigateMock,
}))

describe('RegisterPage', () => {
  beforeEach(() => {
    navigateMock.mockReset()
    vi.restoreAllMocks()
  })

  it('shows recovery codes even if auto-login fails after registration', async () => {
    vi.spyOn(auth, 'login').mockRejectedValueOnce(new Error('login failed'))

    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        recovery_codes: ['REAL-111', 'REAL-222'],
        user: { id: 'u1', username: 'demo', email: 'demo@example.com' },
      }),
    }) as unknown as typeof fetch

    const user = userEvent.setup()
    render(<RegisterPage />)

    await user.type(screen.getByLabelText(/логин/i), 'demo')
    await user.type(screen.getByLabelText(/пароль/i), 'StrongPass1')
    await user.type(screen.getByLabelText(/подтвердите пароль/i), 'StrongPass1')
    await user.click(screen.getByRole('button', { name: /зарегистрироваться/i }))

    const code = await screen.findByText('REAL-111')
    expect(code).toBeTruthy()
  })
})
