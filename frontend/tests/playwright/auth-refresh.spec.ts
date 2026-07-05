import { test, expect } from '@playwright/test'

test('auto-refresh after token expiry (mocked)', async ({ page, context }) => {
  let refreshCalled = false

  // Stub login to return initial access token and set refresh cookie
  await context.route('**/auth/login', (route) =>
    route.fulfill({
      status: 200,
      headers: {
        'Content-Type': 'application/json',
        // set a dummy refresh cookie so auth.tryRefresh can use it
        'Set-Cookie': 'refresh=rt; HttpOnly; Path=/; SameSite=Strict',
      },
      body: JSON.stringify({ access_token: 'token1', expires_in: 900 }),
    }),
  )

  // Protected endpoint: first attempt with token1 -> 401, after refresh -> 200
  await context.route('**/api/v1/protected', async (route) => {
    const req = route.request()
    const auth = req.headers()['authorization'] || ''
    if (!refreshCalled && auth === 'Bearer token1') {
      await route.fulfill({ status: 401, body: 'unauthorized' })
      return
    }
    await route.fulfill({ status: 200, body: 'ok' })
  })

  // Stub refresh endpoint to return a new token
  await context.route('**/auth/refresh', async (route) => {
    refreshCalled = true
    await route.fulfill({
      status: 200,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ access_token: 'token2', expires_in: 900 }),
    })
  })

  await page.goto('/test-auth')

  // perform login via the test UI
  await page.click('#login-btn')
  await expect(page.locator('#status')).toHaveText('logged-in')

  // simulate waiting 15 minutes by relying on the server to respond 401 for the first protected call
  await page.click('#call-btn')

  // final status should be 200 after automatic refresh
  await expect(page.locator('#status')).toHaveText('200')
  expect(refreshCalled).toBeTruthy()
})
