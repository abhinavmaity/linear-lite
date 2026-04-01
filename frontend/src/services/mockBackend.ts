import { ApiError, CollectionResponse, SingleResponse } from 'types/api';
import {
  AuthResponse,
  DashboardStats,
  IssueActivity,
  IssueDetail,
  IssuePriority,
  IssueStatus,
  IssueSummary,
  Label,
  ProjectSummary,
  SprintSummary,
  UserSummary,
} from 'types/domain';

interface MockIssueRecord {
  id: string;
  identifier: string;
  title: string;
  description: string | null;
  status: IssueStatus;
  priority: IssuePriority;
  project_id: string;
  sprint_id: string | null;
  assignee_id: string | null;
  created_by: string;
  label_ids: string[];
  archived_at: string | null;
  archived_by: string | null;
  created_at: string;
  updated_at: string;
}

interface MockState {
  users: UserSummary[];
  labels: Label[];
  projects: Array<{ id: string; name: string; description: string | null; key: string; created_by: string; created_at: string; updated_at: string }>;
  sprints: SprintSummary[];
  issues: MockIssueRecord[];
  activities: IssueActivity[];
  nextIssueNumber: number;
  currentUserId: string;
}

const ISSUE_STATUSES: IssueStatus[] = ['backlog', 'todo', 'in_progress', 'in_review', 'done', 'cancelled'];
const ISSUE_PRIORITIES: IssuePriority[] = ['low', 'medium', 'high', 'urgent'];

const now = Date.now();

const seedUsers: UserSummary[] = [
  makeUser('user-1', 'alex@linear-lite.dev', 'Alex Rivera', 30),
  makeUser('user-2', 'sam@linear-lite.dev', 'Sam Patel', 45),
  makeUser('user-3', 'morgan@linear-lite.dev', 'Morgan Lee', 60),
  makeUser('user-4', 'nina@linear-lite.dev', 'Nina Park', 15),
  makeUser('user-5', 'drew@linear-lite.dev', 'Drew Kim', 10),
];

const seedLabels: Label[] = [
  makeLabel('label-1', 'bug', '#ef4444', 'Defects and regressions', 28),
  makeLabel('label-2', 'feature', '#16a34a', 'Net-new functionality', 24),
  makeLabel('label-3', 'design', '#eab308', 'UI and UX improvements', 20),
  makeLabel('label-4', 'infra', '#2563eb', 'DevOps and platform work', 18),
  makeLabel('label-5', 'frontend', '#ec4899', 'Front-end ownership', 16),
];

const seedProjects: MockState['projects'] = [
  makeProject('project-1', 'Core Platform', 'API and architecture foundations', 'CORE', 'user-1', 90),
  makeProject('project-2', 'Web App', 'Front-end workflows and polish', 'WEB', 'user-2', 80),
  makeProject('project-3', 'Growth', 'Onboarding and activation experiments', 'GROW', 'user-3', 70),
];

const seedSprints: SprintSummary[] = [
  makeSprint('sprint-1', 'Sprint 24', 'project-2', 'active', 4, 10),
  makeSprint('sprint-2', 'Sprint 23', 'project-2', 'completed', 24, 11),
  makeSprint('sprint-3', 'Hardening', 'project-1', 'active', 2, 14),
  makeSprint('sprint-4', 'Launch Prep', 'project-3', 'planned', -3, 14),
];

const seedIssues: MockIssueRecord[] = [
  makeIssue('issue-1', 101, 'Refine issue detail layout spacing', 'in_progress', 'high', 'project-2', 'sprint-1', 'user-1', 'user-2', ['label-3', 'label-5'], 5),
  makeIssue('issue-2', 102, 'Fix keyboard trap in create issue modal', 'todo', 'urgent', 'project-2', 'sprint-1', 'user-4', 'user-1', ['label-1', 'label-5'], 2),
  makeIssue('issue-3', 103, 'Add API pagination metadata tests', 'in_review', 'medium', 'project-1', 'sprint-3', 'user-2', 'user-1', ['label-4'], 7),
  makeIssue('issue-4', 104, 'Implement team directory skeleton', 'done', 'low', 'project-2', 'sprint-2', 'user-3', 'user-2', ['label-2', 'label-5'], 10),
  makeIssue('issue-5', 105, 'Set up CI cache for pnpm store', 'backlog', 'medium', 'project-1', null, null, 'user-1', ['label-4'], 1),
  makeIssue('issue-6', 106, 'Polish onboarding checklist copy', 'cancelled', 'low', 'project-3', null, 'user-5', 'user-3', ['label-2'], 11),
  makeIssue('issue-7', 107, 'Create issue board empty-state illustration', 'todo', 'medium', 'project-2', 'sprint-1', 'user-5', 'user-4', ['label-3', 'label-5'], 3),
  makeIssue('issue-8', 108, 'Audit auth/session edge cases', 'done', 'high', 'project-1', 'sprint-3', 'user-1', 'user-1', ['label-1', 'label-4'], 8),
];

const initialState: MockState = {
  users: seedUsers,
  labels: seedLabels,
  projects: seedProjects,
  sprints: seedSprints,
  issues: seedIssues,
  activities: buildSeedActivities(seedIssues, seedUsers),
  nextIssueNumber: 109,
  currentUserId: 'user-1',
};

const state: MockState = deepClone(initialState);

export function getMockSession(): { token: string; user: UserSummary } {
  const user = findUser(state.currentUserId);
  return { token: tokenForUser(user.id), user };
}

export async function mockRequest<T>(pathWithQuery: string, init?: RequestInit, token?: string | null): Promise<T> {
  await delay(90);
  const method = (init?.method ?? 'GET').toUpperCase();
  const { pathname, searchParams } = parsePath(pathWithQuery);
  const body = parseJsonBody(init?.body);
  const authedUser = getAuthedUser(token);

  if (pathname === '/auth/login' && method === 'POST') {
    const email = String(body?.email ?? '').toLowerCase();
    const user = state.users.find((candidate) => candidate.email.toLowerCase() === email) ?? state.users[0];
    state.currentUserId = user.id;
    return wrap({
      data: {
        token: tokenForUser(user.id),
        expires_at: new Date(now + 1000 * 60 * 60 * 12).toISOString(),
        user,
      },
    } satisfies SingleResponse<AuthResponse>) as T;
  }

  if (pathname === '/auth/register' && method === 'POST') {
    const email = String(body?.email ?? '').trim().toLowerCase();
    const name = String(body?.name ?? '').trim() || 'New User';
    if (!email) {
      throw apiError(400, 'validation_error', 'Email is required.');
    }
    if (state.users.some((user) => user.email.toLowerCase() === email)) {
      throw apiError(409, 'email_exists', 'This email is already registered.');
    }

    const id = `user-${state.users.length + 1}`;
    const createdAt = new Date().toISOString();
    const user: UserSummary = {
      id,
      email,
      name,
      avatar_url: null,
      created_at: createdAt,
      updated_at: createdAt,
    };
    state.users.unshift(user);
    state.currentUserId = user.id;
    return wrap({
      data: {
        token: tokenForUser(user.id),
        expires_at: new Date(Date.now() + 1000 * 60 * 60 * 12).toISOString(),
        user,
      },
    } satisfies SingleResponse<AuthResponse>) as T;
  }

  if (pathname === '/auth/me' && method === 'GET') {
    if (!authedUser) {
      throw apiError(401, 'unauthorized', 'Authentication required.');
    }
    return wrap({ data: authedUser } satisfies SingleResponse<UserSummary>) as T;
  }

  if (!authedUser) {
    throw apiError(401, 'unauthorized', 'Authentication required.');
  }

  if (pathname === '/dashboard/stats' && method === 'GET') {
    return wrap({ data: computeDashboardStats(authedUser.id) } satisfies SingleResponse<DashboardStats>) as T;
  }

  if (pathname === '/users' && method === 'GET') {
    return wrap(listCollection(state.users, searchParams)) as T;
  }

  if (pathname === '/labels' && method === 'GET') {
    return wrap(listCollection(state.labels, searchParams)) as T;
  }

  if (pathname === '/projects' && method === 'GET') {
    return wrap(listCollection(getProjectSummaries(), searchParams)) as T;
  }

  if (pathname === '/sprints' && method === 'GET') {
    let sprints = [...state.sprints];
    const projectId = searchParams.get('project_id');
    if (projectId) {
      sprints = sprints.filter((sprint) => sprint.project_id === projectId);
    }
    return wrap(listCollection(sprints, searchParams)) as T;
  }

  if (pathname === '/issues' && method === 'GET') {
    return wrap(listIssues(searchParams)) as T;
  }

  if (pathname === '/issues' && method === 'POST') {
    const title = String(body?.title ?? '').trim();
    const projectId = String(body?.project_id ?? '');
    if (!title) {
      throw apiError(400, 'validation_error', 'Title is required.');
    }
    if (!state.projects.some((project) => project.id === projectId)) {
      throw apiError(400, 'validation_error', 'Project is required.');
    }
    const created = createIssue({
      title,
      description: asNullableString(body?.description),
      status: asIssueStatus(body?.status),
      priority: asIssuePriority(body?.priority),
      project_id: projectId,
      sprint_id: asNullableString(body?.sprint_id),
      assignee_id: asNullableString(body?.assignee_id),
      label_ids: Array.isArray(body?.label_ids) ? body.label_ids : [],
      created_by: authedUser.id,
    });
    return wrap({ data: toIssueDetail(created) } satisfies SingleResponse<IssueDetail>) as T;
  }

  if (pathname.startsWith('/issues/') && method === 'GET') {
    const issueId = pathname.replace('/issues/', '');
    const includeArchived = searchParams.get('include_archived') === 'true';
    const issue = state.issues.find((item) => item.id === issueId);
    if (!issue || (!includeArchived && issue.archived_at)) {
      throw apiError(404, 'not_found', 'Issue not found.');
    }
    return wrap({ data: toIssueDetail(issue) } satisfies SingleResponse<IssueDetail>) as T;
  }

  if (pathname.startsWith('/issues/') && method === 'PUT') {
    const issueId = pathname.replace('/issues/', '');
    const issue = findIssue(issueId);
    if (issue.archived_at) {
      throw apiError(404, 'not_found', 'Issue not found.');
    }
    const updates = normalizeIssueUpdates(body ?? {});
    applyIssueUpdates(issue, updates, authedUser.id);
    return wrap({ data: toIssueDetail(issue) } satisfies SingleResponse<IssueDetail>) as T;
  }

  if (pathname.startsWith('/issues/') && method === 'DELETE') {
    const issueId = pathname.replace('/issues/', '');
    const issue = findIssue(issueId);
    if (issue.archived_at) {
      return undefined as T;
    }
    issue.archived_at = new Date().toISOString();
    issue.archived_by = authedUser.id;
    issue.updated_at = issue.archived_at;
    pushActivity(issue.id, authedUser.id, 'archived', null, null, null, issue.updated_at);
    return undefined as T;
  }

  throw apiError(404, 'not_found', `No mock endpoint for ${method} ${pathname}`);
}

function listIssues(searchParams: URLSearchParams): CollectionResponse<IssueSummary> {
  let items = state.issues.filter((issue) => !issue.archived_at);
  const search = searchParams.get('search')?.toLowerCase();
  const statuses = searchParams.getAll('status');
  const priorities = searchParams.getAll('priority');
  const labelIds = searchParams.getAll('label_id');
  const assigneeId = searchParams.get('assignee_id');
  const projectId = searchParams.get('project_id');
  const sprintId = searchParams.get('sprint_id');
  const sortBy = searchParams.get('sort_by') ?? 'updated_at';
  const sortOrder = searchParams.get('sort_order') ?? 'desc';

  if (search) {
    items = items.filter(
      (issue) => issue.title.toLowerCase().includes(search) || issue.identifier.toLowerCase().includes(search),
    );
  }
  if (statuses.length) {
    items = items.filter((issue) => statuses.includes(issue.status));
  }
  if (priorities.length) {
    items = items.filter((issue) => priorities.includes(issue.priority));
  }
  if (labelIds.length) {
    items = items.filter((issue) => labelIds.some((labelId) => issue.label_ids.includes(labelId)));
  }
  if (assigneeId) {
    items = items.filter((issue) => issue.assignee_id === assigneeId);
  }
  if (projectId) {
    items = items.filter((issue) => issue.project_id === projectId);
  }
  if (sprintId) {
    items = items.filter((issue) => issue.sprint_id === sprintId);
  }

  items.sort((a, b) => {
    const left = sortField(a, sortBy);
    const right = sortField(b, sortBy);
    if (left < right) return sortOrder === 'asc' ? -1 : 1;
    if (left > right) return sortOrder === 'asc' ? 1 : -1;
    return 0;
  });

  return listCollection(items.map(toIssueSummary), searchParams);
}

function createIssue(input: {
  title: string;
  description: string | null;
  status: IssueStatus;
  priority: IssuePriority;
  project_id: string;
  sprint_id: string | null;
  assignee_id: string | null;
  label_ids: string[];
  created_by: string;
}) {
  const createdAt = new Date().toISOString();
  const issue: MockIssueRecord = {
    id: `issue-${state.nextIssueNumber}`,
    identifier: `LL-${state.nextIssueNumber}`,
    title: input.title,
    description: input.description,
    status: input.status,
    priority: input.priority,
    project_id: input.project_id,
    sprint_id: input.sprint_id,
    assignee_id: input.assignee_id,
    created_by: input.created_by,
    label_ids: input.label_ids.filter((id) => state.labels.some((label) => label.id === id)),
    archived_at: null,
    archived_by: null,
    created_at: createdAt,
    updated_at: createdAt,
  };
  state.nextIssueNumber += 1;
  state.issues.unshift(issue);
  pushActivity(issue.id, input.created_by, 'created', null, null, null, createdAt);
  return issue;
}

function applyIssueUpdates(issue: MockIssueRecord, updates: Partial<MockIssueRecord>, actorId: string) {
  const updatedAt = new Date().toISOString();
  const fields: Array<keyof MockIssueRecord> = [
    'title',
    'description',
    'status',
    'priority',
    'project_id',
    'sprint_id',
    'assignee_id',
  ];

  fields.forEach((field) => {
    const nextValue = updates[field];
    if (nextValue !== undefined && nextValue !== issue[field]) {
      const previous = issue[field];
      issue[field] = nextValue as never;
      pushActivity(issue.id, actorId, 'updated', String(field), stringify(previous), stringify(nextValue), updatedAt);
    }
  });

  if (updates.label_ids && updates.label_ids.join(',') !== issue.label_ids.join(',')) {
    const oldValue = issue.label_ids.map((labelId) => findLabel(labelId).name).join(', ');
    const newValue = updates.label_ids.map((labelId) => findLabel(labelId).name).join(', ');
    issue.label_ids = updates.label_ids;
    pushActivity(issue.id, actorId, 'updated', 'labels', oldValue || null, newValue || null, updatedAt);
  }

  issue.updated_at = updatedAt;
}

function normalizeIssueUpdates(input: Record<string, unknown>): Partial<MockIssueRecord> {
  const next: Partial<MockIssueRecord> = {};
  if (typeof input.title === 'string') next.title = input.title.trim() || 'Untitled Issue';
  if (typeof input.description === 'string' || input.description === null) next.description = input.description;
  if (typeof input.status === 'string' && ISSUE_STATUSES.includes(input.status as IssueStatus)) next.status = input.status as IssueStatus;
  if (typeof input.priority === 'string' && ISSUE_PRIORITIES.includes(input.priority as IssuePriority)) next.priority = input.priority as IssuePriority;
  if (typeof input.project_id === 'string' && state.projects.some((project) => project.id === input.project_id)) {
    next.project_id = input.project_id;
  }
  if (typeof input.sprint_id === 'string' || input.sprint_id === null) next.sprint_id = input.sprint_id;
  if (typeof input.assignee_id === 'string' || input.assignee_id === null) next.assignee_id = input.assignee_id;
  if (Array.isArray(input.label_ids)) {
    next.label_ids = input.label_ids.filter((id): id is string => typeof id === 'string' && state.labels.some((label) => label.id === id));
  }
  return next;
}

function toIssueSummary(issue: MockIssueRecord): IssueSummary {
  return {
    id: issue.id,
    identifier: issue.identifier,
    title: issue.title,
    description: issue.description,
    status: issue.status,
    priority: issue.priority,
    project_id: issue.project_id,
    sprint_id: issue.sprint_id,
    assignee_id: issue.assignee_id,
    created_by: issue.created_by,
    archived_at: issue.archived_at,
    archived_by: issue.archived_by,
    created_at: issue.created_at,
    updated_at: issue.updated_at,
    project: getProjectSummaries().find((project) => project.id === issue.project_id) ?? getProjectSummaries()[0],
    sprint: issue.sprint_id ? state.sprints.find((sprint) => sprint.id === issue.sprint_id) ?? null : null,
    assignee: issue.assignee_id ? state.users.find((user) => user.id === issue.assignee_id) ?? null : null,
    creator: findUser(issue.created_by),
    labels: issue.label_ids.map(findLabel),
  };
}

function toIssueDetail(issue: MockIssueRecord): IssueDetail {
  return {
    ...toIssueSummary(issue),
    activities: state.activities.filter((activity) => activity.issue_id === issue.id).sort((a, b) => (a.created_at < b.created_at ? 1 : -1)),
  };
}

function computeDashboardStats(userId: string): DashboardStats {
  const activeIssues = state.issues.filter((issue) => !issue.archived_at);
  const weekAgo = Date.now() - 7 * 24 * 60 * 60 * 1000;
  return {
    total_issues: activeIssues.length,
    my_issues: activeIssues.filter((issue) => issue.assignee_id === userId).length,
    in_progress: activeIssues.filter((issue) => issue.status === 'in_progress').length,
    done_this_week: activeIssues.filter((issue) => issue.status === 'done' && new Date(issue.updated_at).getTime() >= weekAgo).length,
    active_sprint: state.sprints.find((sprint) => sprint.status === 'active') ?? null,
    recent_activity: state.activities
      .filter((activity) => {
        const issue = state.issues.find((item) => item.id === activity.issue_id);
        return Boolean(issue && !issue.archived_at);
      })
      .sort((a, b) => (a.created_at < b.created_at ? 1 : -1))
      .slice(0, 10),
  };
}

function getProjectSummaries(): ProjectSummary[] {
  return state.projects.map((project) => {
    const projectIssues = state.issues.filter((issue) => issue.project_id === project.id && !issue.archived_at);
    return {
      ...project,
      issue_counts: {
        total: projectIssues.length,
        backlog: projectIssues.filter((issue) => issue.status === 'backlog').length,
        todo: projectIssues.filter((issue) => issue.status === 'todo').length,
        in_progress: projectIssues.filter((issue) => issue.status === 'in_progress').length,
        in_review: projectIssues.filter((issue) => issue.status === 'in_review').length,
        done: projectIssues.filter((issue) => issue.status === 'done').length,
        cancelled: projectIssues.filter((issue) => issue.status === 'cancelled').length,
      },
      active_sprint: state.sprints.find((sprint) => sprint.project_id === project.id && sprint.status === 'active') ?? null,
    };
  });
}

function listCollection<T extends { name?: string; title?: string }>(
  data: T[],
  searchParams: URLSearchParams,
): CollectionResponse<T> {
  const page = Number(searchParams.get('page') ?? 1);
  const limit = Number(searchParams.get('limit') ?? 50);
  const search = searchParams.get('search')?.toLowerCase();
  const filtered = search
    ? data.filter((item) => (item.name ?? item.title ?? '').toLowerCase().includes(search))
    : data;
  const start = Math.max(0, (page - 1) * limit);
  const end = start + limit;
  return {
    items: filtered.slice(start, end),
    pagination: {
      page,
      limit,
      total: filtered.length,
      total_pages: Math.max(1, Math.ceil(filtered.length / limit)),
    },
  };
}

function pushActivity(
  issueId: string,
  userId: string,
  action: string,
  fieldName: string | null,
  oldValue: string | null,
  newValue: string | null,
  createdAt = new Date().toISOString(),
) {
  const activity: IssueActivity = {
    id: `act-${state.activities.length + 1}`,
    issue_id: issueId,
    user_id: userId,
    action,
    field_name: fieldName,
    old_value: oldValue,
    new_value: newValue,
    created_at: createdAt,
    user: findUser(userId),
  };
  state.activities.unshift(activity);
}

function buildSeedActivities(issues: MockIssueRecord[], users: UserSummary[]) {
  const userById = new Map(users.map((user) => [user.id, user]));
  const activities: IssueActivity[] = [];
  issues.forEach((issue, index) => {
    activities.push({
      id: `act-seed-${index * 2 + 1}`,
      issue_id: issue.id,
      user_id: issue.created_by,
      action: 'created',
      field_name: null,
      old_value: null,
      new_value: null,
      created_at: issue.created_at,
      user: userById.get(issue.created_by)!,
    });
    activities.push({
      id: `act-seed-${index * 2 + 2}`,
      issue_id: issue.id,
      user_id: issue.assignee_id ?? issue.created_by,
      action: 'updated',
      field_name: 'status',
      old_value: 'todo',
      new_value: issue.status,
      created_at: issue.updated_at,
      user: userById.get(issue.assignee_id ?? issue.created_by)!,
    });
  });
  return activities.sort((a, b) => (a.created_at < b.created_at ? 1 : -1));
}

function sortField(issue: MockIssueRecord, sortBy: string) {
  if (sortBy === 'identifier') return issue.identifier;
  if (sortBy === 'title') return issue.title.toLowerCase();
  if (sortBy === 'status') return issue.status;
  if (sortBy === 'priority') return issue.priority;
  if (sortBy === 'created_at') return issue.created_at;
  return issue.updated_at;
}

function getAuthedUser(token?: string | null) {
  const userId = tokenToUserId(token);
  if (!userId) return null;
  return state.users.find((user) => user.id === userId) ?? null;
}

function tokenForUser(userId: string) {
  return `mock-token-${userId}`;
}

function tokenToUserId(token?: string | null) {
  if (!token?.startsWith('mock-token-')) return null;
  return token.slice('mock-token-'.length);
}

function findIssue(id: string) {
  const issue = state.issues.find((item) => item.id === id);
  if (!issue) {
    throw apiError(404, 'not_found', 'Issue not found.');
  }
  return issue;
}

function findUser(id: string) {
  const user = state.users.find((candidate) => candidate.id === id);
  if (!user) {
    throw apiError(404, 'not_found', 'User not found.');
  }
  return user;
}

function findLabel(id: string) {
  const label = state.labels.find((candidate) => candidate.id === id);
  if (!label) {
    throw apiError(404, 'not_found', 'Label not found.');
  }
  return label;
}

function parsePath(pathWithQuery: string) {
  const url = new URL(pathWithQuery, 'http://mock.local');
  return { pathname: url.pathname, searchParams: url.searchParams };
}

function parseJsonBody(body: RequestInit['body']) {
  if (!body || typeof body !== 'string') return null;
  try {
    return JSON.parse(body) as Record<string, unknown>;
  } catch {
    return null;
  }
}

function apiError(status: number, code: string, message: string) {
  return new ApiError(status, { code, message });
}

function wrap<T>(value: T) {
  return deepClone(value);
}

function deepClone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T;
}

function stringify(value: unknown) {
  if (value === null || value === undefined || value === '') return null;
  return String(value);
}

function asNullableString(value: unknown) {
  if (typeof value === 'string') {
    return value;
  }
  return null;
}

function asIssueStatus(value: unknown): IssueStatus {
  if (typeof value === 'string' && ISSUE_STATUSES.includes(value as IssueStatus)) {
    return value as IssueStatus;
  }
  return 'backlog';
}

function asIssuePriority(value: unknown): IssuePriority {
  if (typeof value === 'string' && ISSUE_PRIORITIES.includes(value as IssuePriority)) {
    return value as IssuePriority;
  }
  return 'medium';
}

function delay(ms: number) {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

function makeUser(id: string, email: string, name: string, hoursAgo: number): UserSummary {
  const timestamp = new Date(now - hoursAgo * 60 * 60 * 1000).toISOString();
  return {
    id,
    email,
    name,
    avatar_url: null,
    created_at: timestamp,
    updated_at: timestamp,
  };
}

function makeLabel(id: string, name: string, color: string, description: string, daysAgo: number): Label {
  const timestamp = new Date(now - daysAgo * 24 * 60 * 60 * 1000).toISOString();
  return {
    id,
    name,
    color,
    description,
    created_at: timestamp,
    updated_at: timestamp,
  };
}

function makeProject(
  id: string,
  name: string,
  description: string,
  key: string,
  createdBy: string,
  daysAgo: number,
): MockState['projects'][number] {
  const timestamp = new Date(now - daysAgo * 24 * 60 * 60 * 1000).toISOString();
  return {
    id,
    name,
    description,
    key,
    created_by: createdBy,
    created_at: timestamp,
    updated_at: timestamp,
  };
}

function makeSprint(
  id: string,
  name: string,
  projectId: string,
  status: SprintSummary['status'],
  startOffsetDays: number,
  durationDays: number,
): SprintSummary {
  const start = new Date(now + startOffsetDays * 24 * 60 * 60 * 1000);
  const end = new Date(start.getTime() + durationDays * 24 * 60 * 60 * 1000);
  const updated = new Date(start.getTime() + 2 * 60 * 60 * 1000).toISOString();
  return {
    id,
    name,
    description: `${name} focus sprint`,
    project_id: projectId,
    start_date: start.toISOString().slice(0, 10),
    end_date: end.toISOString().slice(0, 10),
    status,
    created_at: start.toISOString(),
    updated_at: updated,
    issue_counts: {
      total: 0,
      backlog: 0,
      todo: 0,
      in_progress: 0,
      in_review: 0,
      done: 0,
      cancelled: 0,
    },
  };
}

function makeIssue(
  id: string,
  number: number,
  title: string,
  status: IssueStatus,
  priority: IssuePriority,
  projectId: string,
  sprintId: string | null,
  assigneeId: string | null,
  createdBy: string,
  labelIds: string[],
  daysAgo: number,
): MockIssueRecord {
  const created = new Date(now - daysAgo * 24 * 60 * 60 * 1000);
  const updated = new Date(created.getTime() + 3 * 60 * 60 * 1000);
  return {
    id,
    identifier: `LL-${number}`,
    title,
    description: `Mock description for "${title}" to validate spacing, typography, and multiline content rendering.`,
    status,
    priority,
    project_id: projectId,
    sprint_id: sprintId,
    assignee_id: assigneeId,
    created_by: createdBy,
    label_ids: labelIds,
    archived_at: null,
    archived_by: null,
    created_at: created.toISOString(),
    updated_at: updated.toISOString(),
  };
}
