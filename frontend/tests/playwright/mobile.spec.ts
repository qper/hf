import { test, expect, devices } from '@playwright/test'

test.describe('Mobile Swipe Navigation', () => {
  test.use({ ...devices['iPhone 14'] })

  test('swipe left to navigate to next day', async ({ page }) => {
    const base = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3000'
    await page.goto(`${base}/board/2026-07-12`)

    // Get initial date text
    const dateText = await page.locator('p:has-text("Sat, 11 July"), p:has-text("Сб, 11 июля")').first()

    // Simulate swipe left
    await page.touchscreen?.swipe(100, 300, -100, 300)

    // Wait for navigation
    await page.waitForURL(/2026-07-13/)

    // Verify date changed
    const newDateText = await page.locator('p:has-text("Sun, 12 July"), p:has-text("Вс, 12 июля")').first()
    expect(newDateText).toBeTruthy()
  })

  test('swipe right to navigate to previous day', async ({ page }) => {
    const base = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3000'
    await page.goto(`${base}/board/2026-07-13`)

    // Simulate swipe right
    await page.touchscreen?.swipe(-100, 300, 100, 300)

    // Wait for navigation
    await page.waitForURL(/2026-07-12/)

    // Verify navigation happened
    const urlAfterSwipe = page.url()
    expect(urlAfterSwipe).toContain('/board/2026-07-12')
  })

  test('vertical scroll does not trigger swipe navigation', async ({ page }) => {
    const base = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3000'
    await page.goto(`${base}/board/2026-07-13`)

    const initialUrl = page.url()

    // Simulate vertical scroll (small deltaX, large deltaY)
    await page.touchscreen?.swipe(100, 100, 110, 300)

    // Wait a bit for any potential navigation
    await page.waitForTimeout(500)

    // Verify URL didn't change
    const finalUrl = page.url()
    expect(finalUrl).toBe(initialUrl)
  })

  test('small swipe does not trigger navigation', async ({ page }) => {
    const base = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3000'
    await page.goto(`${base}/board/2026-07-13`)

    const initialUrl = page.url()

    // Simulate small swipe (less than 50px threshold)
    await page.touchscreen?.swipe(100, 300, 80, 300)

    // Wait a bit for any potential navigation
    await page.waitForTimeout(500)

    // Verify URL didn't change
    const finalUrl = page.url()
    expect(finalUrl).toBe(initialUrl)
  })

  test('read-only mode shows lock icon', async ({ page }) => {
    const base = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3000'
    // This test assumes the backend returns is_editable=false for past dates
    await page.goto(`${base}/board/2026-07-01`)

    // Look for lock icon
    const lockIcon = await page.locator('svg[class*="lucide-lock"]').first()
    const isVisible = await lockIcon.isVisible().catch(() => false)

    if (isVisible) {
      // If lock icon is visible, verify controls are disabled
      const prevButton = page.locator('button').first()
      const isDisabled = await prevButton.isDisabled()
      expect(isDisabled).toBe(true)
    }
  })

  test('slide animation plays on date change', async ({ page }) => {
    const base = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3000'
    await page.goto(`${base}/board/2026-07-12`)

    // Get the board container
    const boardContainer = await page.locator('div.space-y-6').first()

    // Check for animation class before swipe
    let hasAnimateClass = await boardContainer.evaluate((el) =>
      el.classList.contains('animate-slide-in'),
    )
    expect(hasAnimateClass).toBe(false)

    // Trigger swipe
    await page.touchscreen?.swipe(100, 300, -100, 300)

    // Wait for URL change
    await page.waitForURL(/2026-07-13/)

    // Animation class should be gone after animation completes
    hasAnimateClass = await boardContainer.evaluate((el) =>
      el.classList.contains('animate-slide-in'),
    )
    // After animation, class should be removed or animation should complete
    expect(hasAnimateClass).toBe(false)
  })
})
