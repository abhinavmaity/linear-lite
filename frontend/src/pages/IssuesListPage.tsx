import { useMemo } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { Badge } from 'components/common/Badge';
import { Button } from 'components/common/Button';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { PageHeader } from 'components/common/PageHeader';
import { Select } from 'components/common/Select';
import { Spinner } from 'components/common/Spinner';
import { useDebouncedValue } from 'hooks/useDebouncedValue';
import { useIssuesList } from 'features/issues/issuesQueries';
import { useLabelsSelector, useProjectsSelector, useSprintsSelector, useUsersSelector } from 'features/issues/selectorsQueries';
import { titleCase } from 'utils/format';

export function IssuesListPage() {
  const [params, setParams] = useSearchParams();
  const debouncedSearch = useDebouncedValue(params.get('search') ?? '', 500);
  const projects = useProjectsSelector();
  const users = useUsersSelector();
  const labels = useLabelsSelector();
  const sprints = useSprintsSelector(params.get('project_id') ?? undefined);

  const queryParams = useMemo(
    () => ({
      page: Number(params.get('page') ?? 1),
      limit: 20,
      search: debouncedSearch || undefined,
      status: params.get('status') ? [params.get('status') as never] : undefined,
      priority: params.get('priority') ? [params.get('priority') as never] : undefined,
      assignee_id: params.get('assignee_id') || undefined,
      project_id: params.get('project_id') || undefined,
      sprint_id: params.get('sprint_id') || undefined,
      label_id: params.get('label_id') ? [params.get('label_id')!] : undefined,
      label_mode: 'any' as const,
      sort_by: (params.get('sort_by') as never) || 'updated_at',
      sort_order: (params.get('sort_order') as never) || 'desc',
    }),
    [debouncedSearch, params],
  );

  const issues = useIssuesList(queryParams);

  function updateParam(key: string, value: string) {
    const next = new URLSearchParams(params);
    if (value) {
      next.set(key, value);
    } else {
      next.delete(key);
    }
    if (key !== 'page') {
      next.set('page', '1');
    }
    setParams(next);
  }

  return (
    <div>
      <PageHeader
        title="Issues"
        subtitle="Search, filter, sort, and paginate against the canonical issues endpoint."
        actions={
          <Link to="/board">
            <Button variant="ghost">Board View</Button>
          </Link>
        }
      />
      <div className="panel" style={{ padding: 18, marginBottom: 20 }}>
        <div style={{ display: 'grid', gridTemplateColumns: '2fr repeat(5, 1fr)', gap: 12 }}>
          <Input
            placeholder="Search issues"
            value={params.get('search') ?? ''}
            onChange={(event) => updateParam('search', event.target.value)}
          />
          <Select value={params.get('status') ?? ''} onChange={(event) => updateParam('status', event.target.value)}>
            <option value="">All statuses</option>
            <option value="backlog">Backlog</option>
            <option value="todo">Todo</option>
            <option value="in_progress">In Progress</option>
            <option value="in_review">In Review</option>
            <option value="done">Done</option>
            <option value="cancelled">Cancelled</option>
          </Select>
          <Select value={params.get('priority') ?? ''} onChange={(event) => updateParam('priority', event.target.value)}>
            <option value="">All priorities</option>
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="urgent">Urgent</option>
          </Select>
          <Select value={params.get('assignee_id') ?? ''} onChange={(event) => updateParam('assignee_id', event.target.value)}>
            <option value="">All assignees</option>
            {users.data?.map((user) => (
              <option key={user.id} value={user.id}>
                {user.name}
              </option>
            ))}
          </Select>
          <Select value={params.get('project_id') ?? ''} onChange={(event) => updateParam('project_id', event.target.value)}>
            <option value="">All projects</option>
            {projects.data?.map((project) => (
              <option key={project.id} value={project.id}>
                {project.name}
              </option>
            ))}
          </Select>
          <Select value={params.get('sprint_id') ?? ''} onChange={(event) => updateParam('sprint_id', event.target.value)}>
            <option value="">All sprints</option>
            {sprints.data?.map((sprint) => (
              <option key={sprint.id} value={sprint.id}>
                {sprint.name}
              </option>
            ))}
          </Select>
          <Select value={params.get('label_id') ?? ''} onChange={(event) => updateParam('label_id', event.target.value)}>
            <option value="">All labels</option>
            {labels.data?.map((label) => (
              <option key={label.id} value={label.id}>
                {label.name}
              </option>
            ))}
          </Select>
        </div>
      </div>
      {issues.isLoading ? <Spinner label="Loading issues" /> : null}
      {issues.isError ? <ErrorBanner message={(issues.error as Error).message} /> : null}
      {issues.data && issues.data.items.length === 0 ? (
        <EmptyState title="No issues found" description="No issues match the current filters." />
      ) : null}
      {issues.data && issues.data.items.length > 0 ? (
        <div className="panel" style={{ overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ textAlign: 'left', background: 'var(--bg-muted)' }}>
                {['Identifier', 'Title', 'Status', 'Priority', 'Assignee', 'Project'].map((header) => (
                  <th key={header} className="label" style={{ padding: 14 }}>
                    {header}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {issues.data.items.map((issue) => (
                <tr key={issue.id} style={{ borderTop: '1px solid var(--border-soft)' }}>
                  <td style={{ padding: 14 }}>{issue.identifier}</td>
                  <td style={{ padding: 14 }}>
                    <Link to={`/issues/${issue.id}`} style={{ fontWeight: 700 }}>
                      {issue.title}
                    </Link>
                  </td>
                  <td style={{ padding: 14 }}>
                    <Badge>{titleCase(issue.status)}</Badge>
                  </td>
                  <td style={{ padding: 14 }}>
                    <Badge tone="accent">{issue.priority}</Badge>
                  </td>
                  <td style={{ padding: 14 }}>{issue.assignee?.name ?? 'Unassigned'}</td>
                  <td style={{ padding: 14 }}>{issue.project.name}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <div style={{ display: 'flex', justifyContent: 'space-between', padding: 14, borderTop: '1px solid var(--border-soft)' }}>
            <span style={{ color: 'var(--text-secondary)' }}>
              Page {issues.data.pagination.page} of {issues.data.pagination.total_pages || 1}
            </span>
            <div style={{ display: 'flex', gap: 12 }}>
              <Button
                variant="ghost"
                disabled={issues.data.pagination.page <= 1}
                onClick={() => updateParam('page', String(issues.data!.pagination.page - 1))}
              >
                Previous
              </Button>
              <Button
                variant="ghost"
                disabled={issues.data.pagination.page >= issues.data.pagination.total_pages}
                onClick={() => updateParam('page', String(issues.data!.pagination.page + 1))}
              >
                Next
              </Button>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
