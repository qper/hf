import { describe, expect, it } from 'vitest'
import { getAuthErrorMessage } from './auth'

describe('getAuthErrorMessage', () => {
  it('returns a helpful password policy message for invalid registration payloads', () => {
    const message = getAuthErrorMessage(
      422,
      { error: 'invalid registration payload' },
      (key: string) => {
        const map: Record<string, string> = {
          'errors.registrationFailed': 'Ошибка при регистрации',
          'errors.loginFailed': 'Не удалось выполнить вход',
          'errors.passwordTooShort': 'Минимум 8 символов',
          'errors.passwordDigit': 'Нужна хотя бы одна цифра',
          'errors.rateLimitExceeded': 'Слишком много попыток. Подождите немного.',
          'errors.invalidCredentials': 'Неверный логин или пароль',
          'errors.network': 'Ошибка сети',
        }
        return map[key] || key
      },
    )

    expect(message).toContain('8')
    expect(message).toContain('цифра')
  })

  it('maps rate limit responses to a friendly message', () => {
    const message = getAuthErrorMessage(
      429,
      { error: 'rate limit exceeded' },
      (key: string) => {
        const map: Record<string, string> = {
          'errors.rateLimitExceeded': 'Слишком много попыток. Подождите немного.',
        }
        return map[key] || key
      },
    )

    expect(message).toContain('Слишком много попыток')
  })
})
