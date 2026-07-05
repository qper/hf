import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { RouterProvider, createMemoryHistory, createRouter } from '@tanstack/react-router'
import { routeTree } from '@/routeTree.gen'

import { LoginPage } from '@/components/routes/LoginPage'
import { RegisterPage } from '@/components/routes/RegisterPage'

const router = createRouter({
  routeTree,
  history: createMemoryHistory({ initialEntries: ['/login'] }),
})

function renderWithRouter(node: JSX.Element) {
  return render(<RouterProvider router={router}>{node}</RouterProvider>)
}

vi.stubGlobal('fetch', vi.fn(async () => ({
  ok: false,
  status: 401,
} as Response)))

describe('Auth pages', () => {
  it('shows error on wrong password', async () => {
    render(<LoginPage />)

    fireEvent.input(screen.getByLabelText(/username/i), {
      target: { value: 'alice' },
    })
    fireEvent.input(screen.getByLabelText(/password/i), {
      target: { value: 'wrong' },
    })
    fireEvent.submit(screen.getByRole('button', { name: /войти/i }))

    expect(await screen.findByText(/неверный логин или пароль/i)).toBeInTheDocument()
  })

  it('shows inline error for password mismatch on register page', async () => {
    render(<RegisterPage />)

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

    expect(await screen.findByText(/пароли не совпадают/i)).toBeInTheDocument()
  })
})
