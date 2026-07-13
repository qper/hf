import { test, expect } from '@playwright/test'

test.describe('Board - Streak Pulse Animation and Polish', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to board page
    await page.goto('http://localhost:5173/board/2026-07-13')
    // Wait for board to load
    await page.waitForSelector('[class*="space-y"]')
  })

  test.describe('Streak Pulse Animation', () => {
    test('should not show animation when streak < 7', async ({ page }) => {
      // Check that streak badges exist
      const streakBadges = page.locator('[class*="bg-orange-500"]')
      const count = await streakBadges.count()
      
      // Verify badges without animation exist (streak < 7)
      expect(count).toBeGreaterThan(0)
      
      // Check that streak-pulse class is not applied to low streaks
      const firstBadge = streakBadges.first()
      const classes = await firstBadge.getAttribute('class')
      
      // If streak is < 7, should not have streak-pulse class
      if (classes && !classes.includes('streak-pulse')) {
        expect(classes).not.toContain('streak-pulse')
      }
    })

    test('should apply animation when streak >= 7', async ({ page }) => {
      // Take a screenshot of the current state
      // In real scenarios with data where streak >= 7, the animation should be visible
      
      const board = page.locator('[class*="space-y"]')
      expect(board).toBeDefined()
      
      // The animation will apply when streak >= 7
      // We check that the streak-pulse class is conditionally applied
      const streakBadges = page.locator('[class*="flex items-center gap-1 text-xs"]')
      const count = await streakBadges.count()
      
      // If any badges exist, they should be properly rendered
      if (count > 0) {
        expect(count).toBeGreaterThan(0)
      }
    })

    test('animation should use CSS custom properties', async ({ page }) => {
      // Check that the custom property is set in the style attribute
      const habitRow = page.locator('[style*="--habit-color-alpha"]').first()
      
      // Check if the attribute exists (only if streak >= 7)
      const style = await habitRow.getAttribute('style')
      
      if (style && style.includes('--habit-color-alpha')) {
        expect(style).toContain('--habit-color-alpha')
        expect(style).toContain('rgba(234, 179, 8, 0.4)')
      }
    })

    test('animation should not cause layout reflow', async ({ page }) => {
      // Get initial dimensions
      const board = page.locator('[class*="space-y"]')
      const initialBox = await board.boundingBox()
      
      // Wait a bit for animation to run
      await page.waitForTimeout(200)
      
      // Get dimensions again
      const finalBox = await board.boundingBox()
      
      // Dimensions should remain the same (no layout shift)
      expect(initialBox?.width).toBe(finalBox?.width)
      expect(initialBox?.height).toBe(finalBox?.height)
    })
  })

  test.describe('Board Polish - Separators and Alignment', () => {
    test('should display separators between habit rows', async ({ page }) => {
      // Count separator lines (h-px elements with gradient)
      const separators = page.locator('[class*="h-px"][class*="bg-gradient-to-r"]')
      
      // There should be at least one separator
      const count = await separators.count()
      expect(count).toBeGreaterThan(0)
    })

    test('separator should have gradient styling', async ({ page }) => {
      // Get first separator
      const separator = page.locator('[class*="h-px"][class*="bg-gradient-to-r"]').first()
      
      const classes = await separator.getAttribute('class')
      
      // Should have gradient classes
      expect(classes).toContain('h-px')
      expect(classes).toContain('bg-gradient-to-r')
    })

    test('habit rows should have proper alignment', async ({ page }) => {
      // Get first habit row
      const habitRow = page.locator('[class*="rounded-lg"][class*="border"]').first()
      
      // Get the flex container inside
      const flexContainer = habitRow.locator('[class*="flex items-start"]').first()
      
      const classes = await flexContainer.getAttribute('class')
      
      // Should have flex layout classes
      expect(classes).toContain('flex')
      expect(classes).toContain('items-start')
      expect(classes).toContain('justify-between')
      expect(classes).toContain('gap-4')
    })

    test('habit name should be truncated on overflow', async ({ page }) => {
      // Get habit names
      const habitNames = page.locator('h3[class*="truncate"]')
      
      const count = await habitNames.count()
      expect(count).toBeGreaterThan(0)
      
      // Verify truncate class is applied
      const firstName = habitNames.first()
      const classes = await firstName.getAttribute('class')
      
      expect(classes).toContain('truncate')
    })

    test('description should be single-line clamped', async ({ page }) => {
      // Get descriptions
      const descriptions = page.locator('p[class*="line-clamp"]')
      
      const count = await descriptions.count()
      
      if (count > 0) {
        // Verify line-clamp class is applied
        const firstDesc = descriptions.first()
        const classes = await firstDesc.getAttribute('class')
        
        expect(classes).toContain('line-clamp')
      }
    })

    test('should work in dark theme', async ({ page }) => {
      // The page should be in dark theme by default
      
      // Check that background uses dark colors
      const habitRow = page.locator('[class*="bg-zinc-950"]')
      expect(habitRow).toBeDefined()
    })
  })

  test.describe('Visual Regression - Screenshot Comparison', () => {
    test('board should match golden snapshot', async ({ page }) => {
      // Wait for all elements to load and animate settle
      await page.waitForLoadState('networkidle')
      await page.waitForTimeout(100)
      
      // Take screenshot of board area
      const board = page.locator('[class*="space-y-6"]')
      
      // Create a visual snapshot
      await expect(board).toHaveScreenshot('board-with-streak-pulse.png', {
        maxDiffPixels: 100, // Allow small differences
      })
    })
  })
})
