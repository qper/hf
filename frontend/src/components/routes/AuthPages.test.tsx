import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { RouterProvider, createMemoryHistory, createRouter } from '@tanstack/react-router'
import { routeTree } from '@/routeTree.gen'

import { LoginPage } from '@/components/routes/LoginPage'
import { RegisterPage } from '@/components/routes/RegisterPage'

function createTestRouter(initialEntry: string) {
  return createRouter({
    routeTree,
    history: createMemoryHistory({ initialEntries: [initialEntry] }),
  })
}

function renderWithRouter(node: JSX.Element, initialEntry: string) {
  const router = createTestRouter(initialEntry)
  return render(<RouterProvider router={router}>{node}</RouterProvider>)
}

vi.stubGlobal('fetch', vi.fn(async () => ({
  ok: false,
  status: 401,
} as Response)))

describe('Auth pages', () => {
  it('shows error on wrong password', async () => {
    renderWithRouter(<LoginPage />, '/login')

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/password/i), {
      target: { value: 'wrong' },
    })
    fireEvent.submit(screen.getByRole('button', { name: /войти/i }))

    expect(await screen.findByText(/неверный логин или пароль/i)).toBeTruthy()
  })

  it('shows inline error for password mismatch on register page', async () => {
    renderWithRouter(<RegisterPage />, '/register')

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/^password$/i), {
      target: { value: 'Password1!' },
    })
    fireEvent.input(screen.getByLabelText(/confirm password/i), {
      target: { value: 'Password2!' },
    })
    fireEvent.submit(screen.getByRole('button', { name: /зарегистрироваться/i }))

    expect(await screen.findByText(/пароли не совпадают/i)).toBeTruthy()
  })
})
