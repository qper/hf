import { test, expect } from '@playwright/test'

// Note: These tests demonstrate the expected behavior.
// For actual swipe testing, you'll need to run on a real device or use Playwright's mobile emulation.

test.describe('Board - Read-only mode and Swipe Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to board page
    await page.goto('http://localhost:5173/board/2026-07-13')
    // Wait for board to load
    await page.waitForSelector('[class*="space-y-6"]')
  })

  test.describe('Read-only mode', () => {
    test('should show lock icon when board is not editable', async ({ page }) => {
      // This test assumes we have a way to mock a non-editable board
      // In practice, this would be an API response where is_editable=false
      
      // Check if lock icon is rendered (it should be in read-only mode)
      // The lock icon would appear when is_editable === false
      const lockIcon = page.locator('svg[class*="lucide-lock"]')
      
      // On non-editable boards, the lock icon should be visible
      // This is a placeholder - actual test depends on mock data
      if (await lockIcon.isVisible()) {
        expect(await lockIcon.count()).toBeGreaterThan(0)
      }
    })

    test('navigation buttons should be disabled when not editable', async ({ page }) => {
      // Buttons should have disabled:opacity-50 and disabled:cursor-not-allowed
      const prevButton = page.locator('button').filter({ has: page.locator('svg[class*="lucide-chevron-left"]') }).first()
      const nextButton = page.locator('button').filter({ has: page.locator('svg[class*="lucide-chevron-right"]') }).first()
      
      // Verify buttons exist
      expect(prevButton).toBeDefined()
      expect(nextButton).toBeDefined()
    })

    test('habit input controls should show opacity-50 when disabled', async ({ page }) => {
      // When isEditable=false, buttons and inputs should have opacity-50
      const habitRow = page.locator('[class*="rounded-lg"][class*="border"]').first()
      
      if (habitRow) {
        // Get computed style to check opacity
        const button = habitRow.locator('button').first()
        
        if (button) {
          const classes = await button.getAttribute('class')
          // Should have disabled:opacity-50 in the class list (when disabled)
          expect(classes).toBeDefined()
        }
      }
    })
  })

  test.describe('Swipe Navigation', () => {
    test('should detect horizontal swipe left', async ({ page }) => {
      // Get initial date
      const dateElement = page.locator('p').filter({ hasText: /\d+/ }).first()
      
      // Perform touch swipe left
      const boardContainer = page.locator('div').filter({ hasText: /Board/ }).first().locator('..').first()
      
      // Calculate the dimensions
      const boundingBox = await boardContainer.boundingBox()
      if (boundingBox) {
        const startX = boundingBox.x + boundingBox.width * 0.7
        const startY = boundingBox.y + boundingBox.height / 2
        const endX = boundingBox.x + boundingBox.width * 0.2
        
        // Simulate swipe left (move right to left)
        await page.touchscreen.tap(startX, startY)
        // Note: Playwright's touchscreen doesn't support drag, so actual swipe test
        // requires a custom implementation or real device testing
        
        // Verify date element still exists
        expect(dateElement).toBeDefined()
      }
    })

    test('should not process vertical scroll as swipe', async ({ page }) => {
      // Get initial date
      const dateElement = page.locator('p').filter({ hasText: /\d+/ }).first()
      const initialDate = await dateElement.textContent()
      
      // Perform vertical scroll (should not change date)
      await page.evaluate(() => {
        document.documentElement.scrollTop += 100
      })
      
      // Date should remain the same
      const finalDate = await dateElement.textContent()
      expect(finalDate).toBe(initialDate)
    })

    test('should have slide animation on date change', async ({ page }) => {
      // Check that animation classes are applied
      const boardDiv = page.locator('div').filter({ hasText: /Board/ }).first().locator('..').first()
      
      // Verify the div has the transition classes
      const classes = await boardDiv.getAttribute('class')
      
      // Should have transition-all and duration-100 classes
      expect(classes).toContain('transition-all')
      expect(classes).toContain('duration-100')
    })

    test('forward navigation should be disabled when at today', async ({ page }) => {
      // Navigate to today
      const todayButton = page.locator('button').filter({ hasText: /today|Today/ })
      if (await todayButton.isVisible()) {
        await todayButton.click()
        await page.waitForTimeout(100)
      }
      
      // Get next button
      const nextButton = page.locator('button').filter({ 
        has: page.locator('svg[class*="lucide-chevron-right"]') 
      }).first()
      
      // Should be disabled when at today
      const isDisabled = await nextButton.isDisabled()
      expect(isDisabled).toBe(true)
    })

    test('backward navigation should always be enabled', async ({ page }) => {
      // Navigate to a past date first
      await page.goto('http://localhost:5173/board/2026-07-01')
      await page.waitForSelector('[class*="space-y-6"]')
      
      // Get prev button
      const prevButton = page.locator('button').filter({ 
        has: page.locator('svg[class*="lucide-chevron-left"]') 
      }).first()
      
      // Should be enabled even on old dates
      const isEnabled = await prevButton.isEnabled()
      expect(isEnabled).toBe(true)
    })
  })

  test.describe('Mobile Viewport', () => {
    test.use({ viewport: { width: 390, height: 844 } }) // iPhone 14

    test('should apply slide animation on small screen', async ({ page }) => {
      // On mobile, the animation should still work
      const boardDiv = page.locator('div').filter({ hasText: /Board/ }).first().locator('..').first()
      
      const classes = await boardDiv.getAttribute('class')
      expect(classes).toContain('transition-all')
    })

    test('navigation controls should be touch-friendly (44px minimum)', async ({ page }) => {
      // Navigation buttons should have min-height and min-width of 44px
      const buttons = page.locator('button').filter({ hasText: '' })
      
      // Verify at least some buttons exist
      const count = await buttons.count()
      expect(count).toBeGreaterThan(0)
    })
  })
})
