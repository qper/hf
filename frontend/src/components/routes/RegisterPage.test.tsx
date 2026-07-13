import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { RegisterPage } from './RegisterPage'

// Mock router
vi.mock('@tanstack/react-router', async () => {
  const actual = await vi.importActual('@tanstack/react-router')
  return {
    ...actual,
    useNavigate: () => vi.fn(),
  }
})
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render registration form', () => {
    render(<RegisterPage />)
    expect(screen.getByText('Register')).toBeDefined()
    expect(screen.getByLabelText('Username')).toBeDefined()
    expect(screen.getByLabelText('Password')).toBeDefined()
    expect(screen.getByLabelText('Confirm password')).toBeDefined()
  })

  it('should validate username minimum length', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    const usernameInput = screen.getByLabelText('Username') as HTMLInputElement
    const passwordInput = screen.getByLabelText('Password') as HTMLInputElement
    const confirmInput = screen.getByLabelText('Confirm password') as HTMLInputElement
    const submitButton = screen.getByRole('button', { name: /Зарегистрироваться/i })

    await user.type(usernameInput, 'ab')
    await user.type(passwordInput, 'password123')
    await user.type(confirmInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Минимум 3 символа')).toBeDefined()
    })
  })

  it('should validate password minimum length', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    const usernameInput = screen.getByLabelText('Username') as HTMLInputElement
    const passwordInput = screen.getByLabelText('Password') as HTMLInputElement
    const confirmInput = screen.getByLabelText('Confirm password') as HTMLInputElement
    const submitButton = screen.getByRole('button', { name: /Зарегистрироваться/i })

    await user.type(usernameInput, 'testuser')
    await user.type(passwordInput, 'pass')
    await user.type(confirmInput, 'pass')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Минимум 8 символов')).toBeDefined()
    })
  })

  it('should validate password confirmation match', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    const usernameInput = screen.getByLabelText('Username') as HTMLInputElement
    const passwordInput = screen.getByLabelText('Password') as HTMLInputElement
    const confirmInput = screen.getByLabelText('Confirm password') as HTMLInputElement
    const submitButton = screen.getByRole('button', { name: /Зарегистрироваться/i })

    await user.type(usernameInput, 'testuser')
    await user.type(passwordInput, 'password123')
    await user.type(confirmInput, 'password456')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Пароли не совпадают')).toBeDefined()
    })
  })

  it('should show recovery codes dialog on successful registration', async () => {
    const user = userEvent.setup()
    vi.stubGlobal('fetch', vi.fn().mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        recovery_codes: ['CODE1', 'CODE2', 'CODE3', 'CODE4', 'CODE5', 'CODE6', 'CODE7', 'CODE8'],
      }),
    }))

    render(<RegisterPage />)

    const usernameInput = screen.getByLabelText('Username') as HTMLInputElement
    const passwordInput = screen.getByLabelText('Password') as HTMLInputElement
    const confirmInput = screen.getByLabelText('Confirm password') as HTMLInputElement
    const submitButton = screen.getByRole('button', { name: /Зарегистрироваться/i })

    await user.type(usernameInput, 'testuser')
    await user.type(passwordInput, 'password123')
    await user.type(confirmInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Сохраните recovery codes')).toBeDefined()
      expect(screen.getByText('CODE1')).toBeDefined()
    })
  })

  it('should copy codes to clipboard when "Копировать все" clicked', async () => {
    const user = userEvent.setup()
    vi.stubGlobal('fetch', vi.fn().mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        recovery_codes: ['CODE1', 'CODE2'],
      }),
    }))

    const writeTextMock = vi.fn().mockResolvedValueOnce(undefined)
    vi.stubGlobal('navigator', {
      clipboard: {
        writeText: writeTextMock,
      },
    } as any)

    render(<RegisterPage />)

    const usernameInput = screen.getByLabelText('Username')
    const passwordInput = screen.getByLabelText('Password')
    const confirmInput = screen.getByLabelText('Confirm password')
    const submitButton = screen.getByRole('button', { name: /Зарегистрироваться/i })

    await user.type(usernameInput, 'testuser')
    await user.type(passwordInput, 'password123')
    await user.type(confirmInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Сохраните recovery codes')).toBeDefined()
    })

    const copyButton = screen.getByRole('button', { name: /Копировать все/i })
    await user.click(copyButton)

    await waitFor(() => {
      expect(writeTextMock).toHaveBeenCalledWith('CODE1\nCODE2')
      expect(screen.getByText('Скопировано!')).toBeDefined()
    })
  })

  it('should disable "Продолжить" button until checkbox is checked', async () => {
    const user = userEvent.setup()
    vi.stubGlobal('fetch', vi.fn().mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        recovery_codes: ['CODE1', 'CODE2'],
      }),
    }))

    render(<RegisterPage />)

    const usernameInput = screen.getByLabelText('Username')
    const passwordInput = screen.getByLabelText('Password')
    const confirmInput = screen.getByLabelText('Confirm password')
    const submitButton = screen.getByRole('button', { name: /Зарегистрироваться/i })

    await user.type(usernameInput, 'testuser')
    await user.type(passwordInput, 'password123')
    await user.type(confirmInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Сохраните recovery codes')).toBeDefined()
    })

    const continueButton = screen.getByRole('button', { name: /Продолжить/i })
    expect(continueButton).toBeDisabled()

    const checkbox = screen.getByRole('checkbox')
    await user.click(checkbox)

    await waitFor(() => {
      expect(continueButton).not.toBeDisabled()
    })
  })

  it('should show error message on registration failure', async () => {
    const user = userEvent.setup()
    vi.stubGlobal('fetch', vi.fn().mockResolvedValueOnce({
      ok: false,
      json: async () => ({
        message: 'Username already exists',
      }),
    }))

    render(<RegisterPage />)

    const usernameInput = screen.getByLabelText('Username')
    const passwordInput = screen.getByLabelText('Password')
    const confirmInput = screen.getByLabelText('Confirm password')
    const submitButton = screen.getByRole('button', { name: /Зарегистрироваться/i })

    await user.type(usernameInput, 'testuser')
    await user.type(passwordInput, 'password123')
    await user.type(confirmInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Username already exists')).toBeDefined()
    })
  })

  it('should handle network errors gracefully', async () => {
    const user = userEvent.setup()
    vi.stubGlobal('fetch', vi.fn().mockRejectedValueOnce(new Error('Network error')))

    render(<RegisterPage />)

    const usernameInput = screen.getByLabelText('Username')
    const passwordInput = screen.getByLabelText('Password')
    const confirmInput = screen.getByLabelText('Confirm password')
    const submitButton = screen.getByRole('button', { name: /Зарегистрироваться/i })

    await user.type(usernameInput, 'testuser')
    await user.type(passwordInput, 'password123')
    await user.type(confirmInput, 'password123')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText('Network error')).toBeDefined()
    })
  })
})
