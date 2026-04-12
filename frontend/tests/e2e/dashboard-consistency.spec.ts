import { expect, test, type Page } from '@playwright/test';
import { createIssue, createProject, createSprint, registerViaApi, signInWithToken, uniqueSeed } from './helpers';

async function dashboardStatValue(page: Page, label: string) {
  const card = page
    .locator('.panel')
    .filter({ has: page.locator('.label', { hasText: label }) })
    .first();
  const valueText = (await card.locator('div').nth(1).textContent()) ?? '';
  return Number(valueText.trim());
}

test.describe('M6-11 Dashboard Consistency Journey', () => {
  test('dashboard total issues reflects issue mutations', async ({ page, request }) => {
    const seed = uniqueSeed('dashboard');
    const session = await registerViaApi(request, seed);
    const project = await createProject(request, session.token, seed);
    const sprint = await createSprint(request, session.token, project.id, seed, 'planned');

    await createIssue(request, session.token, {
      title: `Dashboard Issue A ${seed}`,
      project_id: project.id,
      sprint_id: sprint.id,
      assignee_id: session.user.id,
    });

    await signInWithToken(page, session.token);
    await page.goto('/dashboard');
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
    const before = await dashboardStatValue(page, 'Total Issues');

    await createIssue(request, session.token, {
      title: `Dashboard Issue B ${seed}`,
      project_id: project.id,
      sprint_id: sprint.id,
      assignee_id: session.user.id,
    });

    await page.reload();
    const after = await dashboardStatValue(page, 'Total Issues');
    expect(after).toBe(before + 1);
  });
});
