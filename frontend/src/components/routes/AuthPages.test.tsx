import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import '@/i18n'
import i18n from '@/i18n'
import App from '@/App'

function renderWithRouter(initialEntry: string) {
  window.history.pushState({}, '', initialEntry)
  return render(<App />)
}

const fetchMock = vi.fn()

vi.stubGlobal('fetch', fetchMock)

describe('Auth pages', () => {
  beforeEach(async () => {
    fetchMock.mockReset()
    fetchMock.mockResolvedValue(new Response(null, { status: 401 }))
    await i18n.changeLanguage('ru')
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })
  it('shows error on wrong password', async () => {
    renderWithRouter('/login')

    await screen.findByLabelText(/username/i)

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/пароль|password/i), {
      target: { value: 'wrong' },
    })
    fireEvent.submit(screen.getByRole('button', { name: /войти/i }))

    expect(await screen.findByText(/неверный логин или пароль/i)).toBeTruthy()
  })

  it('navigates to the board after a successful login', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify({ access_token: 'abc' }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    renderWithRouter('/login')

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/пароль|password/i), {
      target: { value: 'Password1!' },
    })
    fireEvent.submit(screen.getByRole('button', { name: /войти/i }))

    await waitFor(() => {
      expect(window.location.pathname).toMatch(/^\/board\/\d{4}-\d{2}-\d{2}$/)
    })
  })

  it('shows inline error for password mismatch on register page', async () => {
    renderWithRouter('/register')

    await screen.findByLabelText(/подтверд|confirm/i)

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/пароль|password/i), {
      target: { value: 'Password1!' },
    })
    fireEvent.input(screen.getByLabelText(/подтверд|confirm/i), {
      target: { value: 'Password2!' },
    })
    fireEvent.submit(
      screen.getByRole('button', { name: /зарегистрироваться/i }),
    )

    expect(await screen.findByText(/пароли не совпадают/i)).toBeTruthy()
  })

  it('shows a password policy error for weak passwords', async () => {
    renderWithRouter('/register')

    await screen.findByLabelText(/пароль|password/i)

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/пароль|password/i), {
      target: { value: 'short' },
    })
    fireEvent.input(screen.getByLabelText(/подтверд|confirm/i), {
      target: { value: 'short' },
    })
    fireEvent.submit(
      screen.getByRole('button', { name: /зарегистрироваться/i }),
    )

    expect(await screen.findByText(/минимум 8 символов/i)).toBeTruthy()
  })
})
