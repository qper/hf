import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import App from '@/App'

function renderWithRouter(initialEntry: string) {
  window.history.pushState({}, '', initialEntry)
  return render(<App />)
}

const fetchMock = vi.fn()

vi.stubGlobal('fetch', fetchMock)

describe('Auth pages', () => {
  beforeEach(() => {
    fetchMock.mockReset()
    fetchMock.mockResolvedValue(new Response(null, { status: 401 }))
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
    fireEvent.input(screen.getByLabelText(/password/i), {
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
    fireEvent.input(screen.getByLabelText(/password/i), {
      target: { value: 'Password1!' },
    })
    fireEvent.submit(screen.getByRole('button', { name: /войти/i }))

    await waitFor(() => {
      expect(window.location.pathname).toMatch(/^\/board\/\d{4}-\d{2}-\d{2}$/)
    })
  })

  it('shows inline error for password mismatch on register page', async () => {
    renderWithRouter('/register')

    await screen.findByLabelText(/confirm password/i)

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/^password$/i), {
      target: { value: 'Password1!' },
    })
    fireEvent.input(screen.getByLabelText(/confirm password/i), {
      target: { value: 'Password2!' },
    })
    fireEvent.submit(
      screen.getByRole('button', { name: /зарегистрироваться/i }),
    )

    expect(await screen.findByText(/пароли не совпадают/i)).toBeTruthy()
  })

  it('shows a password policy error for weak passwords', async () => {
    renderWithRouter('/register')

    await screen.findByLabelText(/^password$/i)

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/^password$/i), {
      target: { value: 'short' },
    })
    fireEvent.input(screen.getByLabelText(/confirm password/i), {
      target: { value: 'short' },
    })
    fireEvent.submit(
      screen.getByRole('button', { name: /зарегистрироваться/i }),
    )

    expect(await screen.findByText(/минимум 8 символов/i)).toBeTruthy()
  })
})
