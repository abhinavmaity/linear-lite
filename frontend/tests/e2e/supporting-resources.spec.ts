import { expect, test } from '@playwright/test';
import { createProject, registerViaApi, signInWithToken, uniqueSeed } from './helpers';

test.describe('M6-10 Supporting Resources Journey', () => {
  test('projects, sprints, and labels handle happy path and conflict states', async ({ page, request }) => {
    const seed = uniqueSeed('resources');
    const session = await registerViaApi(request, seed);
    await signInWithToken(page, session.token);

    const projectKey = `R${Date.now().toString().slice(-5)}`;

    await page.goto('/projects');
    await expect(page.getByRole('heading', { name: 'Projects' })).toBeVisible();

    const projectCreateForm = page.locator('form').first();
    await projectCreateForm.getByPlaceholder('Platform').fill(`Project ${seed}`);
    await projectCreateForm.getByPlaceholder('PLAT', { exact: true }).fill(projectKey);
    await projectCreateForm.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText('Project created.')).toBeVisible();

    await projectCreateForm.getByPlaceholder('Platform').fill(`Project Duplicate ${seed}`);
    await projectCreateForm.getByPlaceholder('PLAT', { exact: true }).fill(projectKey);
    await projectCreateForm.getByRole('button', { name: 'Create' }).click();
    await expect(projectCreateForm.getByText(/already exists|already in use/i).first()).toBeVisible();

    const createdProject = await createProject(request, session.token, `${seed}-sprint-parent`);

    await page.goto('/sprints');
    await expect(page.getByRole('heading', { name: 'Sprints' })).toBeVisible();
    const sprintCreateForm = page.locator('form').first();
    await sprintCreateForm.locator('input').nth(0).fill(`Sprint A ${seed}`);
    await sprintCreateForm.locator('select').nth(0).selectOption(createdProject.id);
    await sprintCreateForm.locator('input[type="date"]').nth(0).fill('2026-04-12');
    await sprintCreateForm.locator('input[type="date"]').nth(1).fill('2026-04-19');
    await sprintCreateForm.locator('select').nth(1).selectOption('active');
    await sprintCreateForm.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText('Sprint created.')).toBeVisible();

    await sprintCreateForm.locator('input').nth(0).fill(`Sprint B ${seed}`);
    await sprintCreateForm.locator('select').nth(0).selectOption(createdProject.id);
    await sprintCreateForm.locator('input[type="date"]').nth(0).fill('2026-04-20');
    await sprintCreateForm.locator('input[type="date"]').nth(1).fill('2026-04-26');
    await sprintCreateForm.locator('select').nth(1).selectOption('active');
    await sprintCreateForm.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText(/only one active sprint/i)).toBeVisible();

    await page.goto('/labels');
    await expect(page.getByRole('heading', { name: 'Labels' })).toBeVisible();
    const labelsCreateForm = page.locator('form').first();
    const labelName = `label-${seed}`;
    await labelsCreateForm.getByPlaceholder('bug').fill(labelName);
    await labelsCreateForm.locator('select').first().selectOption('#22C55E');
    await labelsCreateForm.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText('Label created.')).toBeVisible();

    await labelsCreateForm.getByPlaceholder('bug').fill(labelName);
    await labelsCreateForm.locator('select').first().selectOption('#EF4444');
    await labelsCreateForm.getByRole('button', { name: 'Create' }).click();
    await expect(labelsCreateForm.getByText(/already exists|already in use/i).first()).toBeVisible();
  });
});
