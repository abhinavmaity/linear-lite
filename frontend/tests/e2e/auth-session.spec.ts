import { test, expect } from '@playwright/test';

test.describe('M6-08 Auth + Session Journey', () => {
  test('redirects unauthenticated users from protected routes to login', async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page).toHaveURL(/\/login$/);
    await expect(page.getByText('Login')).toBeVisible();
  });

  test('registers via UI and preserves session after refresh', async ({ page }) => {
    const seed = `auth-${Date.now()}`;
    const email = `${seed}@example.com`;

    await page.goto('/register');
    await page.getByText('Register').first().waitFor();

    await page.locator('input').nth(0).fill(`User ${seed}`);
    await page.locator('input').nth(1).fill(email);
    await page.locator('input').nth(2).fill('Password123');
    await page.locator('input').nth(3).fill('Password123');
    await page.getByRole('button', { name: 'Create Account' }).click();

    await expect(page).toHaveURL(/\/login$/);
    await expect(page.getByText(/Account created/i)).toBeVisible();

    await page.locator('input[type="email"]').fill(email);
    await page.locator('input[type="password"]').fill('Password123');
    await page.getByRole('button', { name: 'Enter Dashboard' }).click();

    await expect(page).toHaveURL(/\/dashboard$/);
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();

    await page.reload();
    await expect(page).toHaveURL(/\/dashboard$/);
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
  });
});
