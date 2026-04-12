import { APIRequestContext, Page, expect } from '@playwright/test';

const API_BASE_URL = process.env.E2E_API_BASE_URL ?? 'http://127.0.0.1:8080/api/v1';
const TOKEN_KEY = 'linear-lite-token';

type AuthSession = {
  token: string;
  user: {
    id: string;
    email: string;
    name: string;
  };
};

export function uniqueSeed(prefix: string) {
  return `${prefix}-${Date.now()}-${Math.floor(Math.random() * 10_000)}`;
}

export async function registerViaApi(request: APIRequestContext, seed: string): Promise<AuthSession> {
  const email = `${seed}@example.com`;
  const response = await request.post(`${API_BASE_URL}/auth/register`, {
    data: {
      name: `E2E ${seed}`,
      email,
      password: 'Password123',
    },
  });
  expect(response.ok()).toBeTruthy();
  const body = await response.json();
  return body.data as AuthSession;
}

export async function authHeaders(token: string) {
  return {
    Authorization: `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
}

export async function createProject(request: APIRequestContext, token: string, seed: string) {
  const timestampPart = Date.now().toString().slice(-5);
  const randomPart = Math.floor(Math.random() * 100).toString().padStart(2, '0');
  const key = `E${timestampPart}${randomPart}`;
  const response = await request.post(`${API_BASE_URL}/projects`, {
    headers: await authHeaders(token),
    data: {
      name: `Project ${seed}`,
      key,
      description: `E2E project ${seed}`,
    },
  });
  if (!response.ok()) {
    const body = await response.text();
    throw new Error(`createProject failed: status=${response.status()} body=${body}`);
  }
  const body = await response.json();
  return body.data;
}

export async function createSprint(request: APIRequestContext, token: string, projectId: string, seed: string, status = 'planned') {
  const response = await request.post(`${API_BASE_URL}/sprints`, {
    headers: await authHeaders(token),
    data: {
      name: `Sprint ${seed}`,
      description: `E2E sprint ${seed}`,
      project_id: projectId,
      start_date: '2026-04-12',
      end_date: '2026-04-19',
      status,
    },
  });
  expect(response.ok()).toBeTruthy();
  const body = await response.json();
  return body.data;
}

export async function createLabel(request: APIRequestContext, token: string, seed: string) {
  const response = await request.post(`${API_BASE_URL}/labels`, {
    headers: await authHeaders(token),
    data: {
      name: `label-${seed}`,
      color: '#EF4444',
      description: `E2E label ${seed}`,
    },
  });
  expect(response.ok()).toBeTruthy();
  const body = await response.json();
  return body.data;
}

export async function createIssue(
  request: APIRequestContext,
  token: string,
  input: {
    title: string;
    project_id: string;
    sprint_id?: string | null;
    assignee_id?: string | null;
    label_ids?: string[];
    status?: string;
    priority?: string;
  },
) {
  const response = await request.post(`${API_BASE_URL}/issues`, {
    headers: await authHeaders(token),
    data: {
      ...input,
      description: `E2E issue ${input.title}`,
      status: input.status ?? 'backlog',
      priority: input.priority ?? 'medium',
    },
  });
  expect(response.ok()).toBeTruthy();
  const body = await response.json();
  return body.data;
}

export async function restoreIssue(request: APIRequestContext, token: string, issueId: string) {
  const response = await request.put(`${API_BASE_URL}/issues/${issueId}`, {
    headers: await authHeaders(token),
    data: { archived: false },
  });
  expect(response.ok()).toBeTruthy();
}

export async function signInWithToken(page: Page, token: string) {
  await page.addInitScript(
    ([key, value]) => {
      localStorage.setItem(key, value);
    },
    [TOKEN_KEY, token],
  );
  await page.goto('/dashboard');
  await expect(page).toHaveURL(/\/dashboard$/);
}
