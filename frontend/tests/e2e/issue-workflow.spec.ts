import { expect, test } from '@playwright/test';
import {
  createIssue,
  createLabel,
  createProject,
  createSprint,
  registerViaApi,
  restoreIssue,
  signInWithToken,
  uniqueSeed,
} from './helpers';

test.describe('M6-09 Core Issue Workflow Journey', () => {
  test('create, update, board-verify, archive, and restore issue', async ({ page, request }) => {
    const seed = uniqueSeed('issue-flow');
    const session = await registerViaApi(request, seed);
    const project = await createProject(request, session.token, seed);
    const sprint = await createSprint(request, session.token, project.id, seed, 'planned');
    const label = await createLabel(request, session.token, seed);
    const title = `Issue ${seed}`;
    const createdIssue = await createIssue(request, session.token, {
      title,
      project_id: project.id,
      sprint_id: sprint.id,
      assignee_id: session.user.id,
      label_ids: [label.id],
    });

    await signInWithToken(page, session.token);
    await page.goto(`/issues/${createdIssue.id}`);
    await expect(page).toHaveURL(new RegExp(`/issues/${createdIssue.id}$`));
    const issueId = createdIssue.id;
    await expect(page.locator('input').first()).toHaveValue(title);

    const statusSelect = page.locator('select').first();
    await statusSelect.selectOption('in_progress');
    await expect(statusSelect).toHaveValue('in_progress');

    await page.getByRole('link', { name: 'Board', exact: true }).click();
    await expect(page).toHaveURL(/\/board$/);
    await page.getByPlaceholder('Search issues').fill(title);
    await expect(page.getByText(title)).toBeVisible();

    page.once('dialog', (dialog) => dialog.accept());
    await page.getByText(title).first().click();
    await expect(page).toHaveURL(new RegExp(`/issues/${issueId}$`));
    await page.getByRole('button', { name: 'Archive Issue' }).click();
    await expect(page).toHaveURL(/\/issues$/);

    await restoreIssue(request, session.token, issueId);
    await page.goto(`/issues/${issueId}`);
    await expect(page.locator('input').first()).toHaveValue(title);
  });
});
