# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: auth-flow.spec.ts >> registers a new user and logs in with the same credentials
- Location: tests/playwright/auth-flow.spec.ts:7:1

# Error details

```
Test timeout of 60000ms exceeded.
```

```
Error: locator.click: Test timeout of 60000ms exceeded.
Call log:
  - waiting for getByRole('link', { name: /зарегистрироваться/i })

```

# Page snapshot

```yaml
- generic [ref=e2]:
  - generic [ref=e3]:
    - banner [ref=e4]:
      - generic [ref=e5]:
        - generic [ref=e6]:
          - paragraph [ref=e7]: HabitFlow
          - heading "Modern frontend shell" [level=1] [ref=e8]
        - navigation [ref=e9]:
          - link "Home" [ref=e10] [cursor=pointer]:
            - /url: /
          - link "Login" [ref=e11] [cursor=pointer]:
            - /url: /login
          - link "Register" [ref=e12] [cursor=pointer]:
            - /url: /register
          - button "Hide side panel" [ref=e13]
    - main [ref=e14]:
      - complementary [ref=e15]:
        - paragraph [ref=e16]: Session state
        - paragraph [ref=e17]: The UI store persists in session storage while this tab stays open.
      - generic [ref=e20]:
        - generic [ref=e21]:
          - paragraph [ref=e22]: Access
          - heading "Sign in" [level=2] [ref=e23]
          - paragraph [ref=e24]: Sign in to continue working with HabitFlow.
        - generic [ref=e25]:
          - generic [ref=e26]:
            - generic [ref=e27]: Username
            - textbox "Username" [ref=e28]
          - generic [ref=e29]:
            - generic [ref=e30]: Password
            - textbox "Password" [ref=e31]
          - generic [ref=e32]:
            - button "Sign in" [ref=e33]
            - generic [ref=e34]:
              - link "Create an account" [ref=e35] [cursor=pointer]:
                - /url: /register
              - link "Sign in with a recovery code" [ref=e36] [cursor=pointer]:
                - /url: /register
  - generic [ref=e37]:
    - img [ref=e39]
    - button "Open Tanstack query devtools" [ref=e87] [cursor=pointer]:
      - img [ref=e88]
  - generic:
    - contentinfo:
      - button "Open TanStack Router Devtools" [ref=e136] [cursor=pointer]:
        - generic [ref=e137]:
          - img [ref=e139]
          - img [ref=e174]
        - generic [ref=e208]: "-"
        - generic [ref=e209]: TanStack Router
```

# Test source

```ts
  1  | import { test, expect } from '@playwright/test'
  2  | 
  3  | const unique = `qa_${Date.now()}`
  4  | const username = `${unique}_user`
  5  | const password = `StrongPass1${Date.now()}`
  6  | 
  7  | test('registers a new user and logs in with the same credentials', async ({ page }) => {
  8  |   await page.goto('/login')
  9  | 
> 10 |   await page.getByRole('link', { name: /зарегистрироваться/i }).click()
     |                                                                 ^ Error: locator.click: Test timeout of 60000ms exceeded.
  11 |   await page.getByLabel(/логин/i).fill(username)
  12 |   await page.getByLabel(/пароль/i).fill(password)
  13 |   await page.getByLabel(/подтвердите пароль/i).fill(password)
  14 |   await page.getByRole('button', { name: /зарегистрироваться/i }).click()
  15 | 
  16 |   await expect(page.getByText(/коды восстановления/i)).toBeVisible({ timeout: 10000 })
  17 |   await page.getByLabel(/я сохранил/i).check()
  18 |   await page.getByRole('button', { name: /продолжить/i }).click()
  19 | 
  20 |   await expect(page).toHaveURL(/\/board\//)
  21 | 
  22 |   await page.goto('/login')
  23 |   await page.getByLabel(/логин/i).fill(username)
  24 |   await page.getByLabel(/пароль/i).fill(password)
  25 |   await page.getByRole('button', { name: /войти/i }).click()
  26 | 
  27 |   await expect(page).toHaveURL(/\/board\//)
  28 | })
  29 | 
```