/*
 * Lightweight auth client: keeps access token in closure,
 * provides request helper that auto-refreshes on 401 and
 * exposes login/logout helpers.
 */

export type LoginResponse = {
  access_token?: string
  refresh_token?: string
  expires_in?: number
  token_type?: string
}

let accessToken: string | null = null
let refreshPromise: Promise<void> | null = null

export function getAccessToken(): string | null {
  return accessToken
}

export function setAccessToken(token: string | null) {
  accessToken = token
}

async function doRefresh(): Promise<void> {
  if (refreshPromise) return refreshPromise

  refreshPromise = (async () => {
    const res = await fetch('/auth/refresh', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
    })

    if (!res.ok) {
      accessToken = null
      throw new Error('refresh_failed')
    }

    const data = await res.json().catch(() => ({})) as LoginResponse
    accessToken = data.access_token || null
    refreshPromise = null
  })()

  return refreshPromise
}

export async function tryRefresh(): Promise<boolean> {
  try {
    await doRefresh()
    return accessToken !== null
  } catch {
    refreshPromise = null
    return false
  }
}

export async function authFetch(input: unknown, init: Record<string, unknown> = {}): Promise<Response> {
  const mergedInit = { credentials: 'include', ...(init as Record<string, unknown>) } as Record<string, unknown>

  const existingHeaders = ((mergedInit.headers as unknown) as Record<string, string>) || {}
  const headers: Record<string, string> = { ...existingHeaders }

  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`
  }

  mergedInit.headers = headers as unknown

  const nativeFetch = (fetch as unknown) as (input: unknown, init?: unknown) => Promise<Response>
  let res = await nativeFetch(input, mergedInit)

  if (res.status !== 401) return res

  const refreshed = await tryRefresh()
  if (!refreshed) {
    // failed refresh: force client to login
    try {
      window.location.assign('/login')
    } catch (e) {
      void e
    }
    throw new Error('session_expired')
  }

  // retry once with new token
  ;(mergedInit.headers as unknown as Record<string, string>)['Authorization'] = `Bearer ${accessToken}`
  res = await nativeFetch(input, mergedInit)
  return res
}

export function getAuthErrorMessage(
  status: number,
  data: unknown,
  t: (key: string) => string,
): string {
  const payload = (data as Record<string, unknown> | null) || {}
  const message = typeof payload.error === 'string' ? payload.error : ''
  const detail = typeof payload.message === 'string' ? payload.message : ''

  if (status === 429 || message === 'rate limit exceeded' || detail === 'rate limit exceeded') {
    return t('errors.rateLimitExceeded')
  }

  if (status === 401 || message === 'invalid credentials' || detail === 'invalid credentials') {
    return t('errors.invalidCredentials')
  }

  if (status === 422 || message === 'invalid registration payload' || detail === 'invalid registration payload') {
    return `${t('errors.registrationFailed')}: ${t('errors.passwordTooShort')} / ${t('errors.passwordDigit')}`
  }

  return typeof message === 'string' && message.length > 0
    ? message
    : t('errors.loginFailed')
}

export async function login(username: string, password: string): Promise<LoginResponse> {
  const res = await fetch('/auth/login', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })

  if (!res.ok) {
    throw res
  }

  const data = (await res.json()) as LoginResponse
  accessToken = data.access_token || null
  return data
}

export async function logout(): Promise<void> {
  try {
    await fetch('/auth/logout', { method: 'POST', credentials: 'include' })
  } catch (e) {
    void e
  }
  accessToken = null
  try {
    window.location.assign('/login')
  } catch (e) {
    void e
  }
}

export default {
  getAccessToken,
  setAccessToken,
  tryRefresh,
  authFetch,
  login,
  logout,
}
