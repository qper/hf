import { test, expect } from '@playwright/test'

const unique = `qa_${Date.now()}`
const username = `${unique}_user`
const password = `StrongPass1${Date.now()}`

test('registers a new user and logs in with the same credentials', async ({ page }) => {
  await page.goto('/login')

  await page.getByRole('link', { name: /зарегистрироваться/i }).click()
  await page.getByLabel(/логин/i).fill(username)
  await page.getByLabel(/пароль/i).fill(password)
  await page.getByLabel(/подтвердите пароль/i).fill(password)
  await page.getByRole('button', { name: /зарегистрироваться/i }).click()

  await expect(page.getByText(/коды восстановления/i)).toBeVisible({ timeout: 10000 })
  await page.getByLabel(/я сохранил/i).check()
  await page.getByRole('button', { name: /продолжить/i }).click()

  await expect(page).toHaveURL(/\/board\//)

  await page.goto('/login')
  await page.getByLabel(/логин/i).fill(username)
  await page.getByLabel(/пароль/i).fill(password)
  await page.getByRole('button', { name: /войти/i }).click()

  await expect(page).toHaveURL(/\/board\//)
})
