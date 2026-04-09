import { useMemo, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Skeleton } from 'boneyard-js/react';
import { Input } from 'components/common/Input';
import { Select } from 'components/common/Select';
import { Button } from 'components/common/Button';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { IssueCard } from 'components/issues/IssueCard';
import { useIssuesList } from 'features/issues/issuesQueries';
import { useDebouncedValue } from 'hooks/useDebouncedValue';
import { useLabelsSelector, useProjectsSelector, useSprintsSelector, useUsersSelector } from 'features/issues/selectorsQueries';
import { issuesApi } from 'services/issuesApi';
import { useUIStore } from 'store/uiStore';
import { IssueSummary } from 'types/domain';

type BoardColumn = 'backlog' | 'todo' | 'in_progress' | 'in_review' | 'done' | 'cancelled';

const columns: BoardColumn[] = ['backlog', 'todo', 'in_progress', 'in_review', 'done', 'cancelled'];

export function IssuesBoardPage() {
  const queryClient = useQueryClient();
  const pushToast = useUIStore((state) => state.pushToast);
  const [search, setSearch] = useState('');
  const debouncedSearch = useDebouncedValue(search, 500);
  const [status, setStatus] = useState('');
  const [priority, setPriority] = useState('');
  const [assigneeId, setAssigneeId] = useState('');
  const [projectId, setProjectId] = useState('');
  const [sprintId, setSprintId] = useState('');
  const [labelId, setLabelId] = useState('');
  const users = useUsersSelector();
  const projects = useProjectsSelector();
  const labels = useLabelsSelector();
  const sprints = useSprintsSelector(projectId || undefined);
  const [items, setItems] = useState<ReturnType<typeof groupByStatus> | null>(null);
  const [movingIssueId, setMovingIssueId] = useState<string | null>(null);
  const issues = useIssuesList({
    page: 1,
    limit: 100,
    sort_by: 'updated_at',
    sort_order: 'desc',
    search: debouncedSearch || undefined,
    status: status ? [status as BoardColumn] : undefined,
    priority: priority ? [priority as 'low' | 'medium' | 'high' | 'urgent'] : undefined,
    assignee_id: assigneeId || undefined,
    project_id: projectId || undefined,
    sprint_id: sprintId || undefined,
    label_id: labelId ? [labelId] : undefined,
    label_mode: 'any',
  });
  const loading = issues.isLoading || users.isLoading || projects.isLoading || labels.isLoading || sprints.isLoading;

  const grouped = useMemo(() => {
    const base = groupByStatus(issues.data?.items ?? []);
    return items ?? base;
  }, [issues.data?.items, items]);

  async function moveIssue(issueId: string, nextStatus: BoardColumn) {
    if (movingIssueId) return;
    if (!issues.data) return;
    const moved = issues.data.items.find((issue) => issue.id === issueId);
    if (!moved || moved.status === nextStatus) return;
    const previous = groupByStatus(issues.data.items);
    const optimistic = groupByStatus(
      issues.data.items.map((issue) => (issue.id === issueId ? { ...issue, status: nextStatus } : issue)),
    );
    setItems(optimistic);
    setMovingIssueId(issueId);

    try {
      await issuesApi.update(issueId, { status: nextStatus });
      pushToast({ tone: 'success', message: 'Issue updated.' });
      queryClient.invalidateQueries({ queryKey: ['issues'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
      setItems(null);
    } catch (error) {
      setItems(previous);
      pushToast({ tone: 'error', message: error instanceof Error ? error.message : 'Failed to update issue.' });
    } finally {
      setMovingIssueId(null);
    }
  }

  return (
    <div>
      <PageHeader title="Board" subtitle="Shared issue query layer with optimistic drag-and-drop status updates." />
      <Skeleton
        name="issues-board-page"
        loading={loading}
        fallback={<Spinner label="Loading board" />}
        fixture={
          <div>
            <div className="panel" style={{ padding: 18, marginBottom: 20 }}>
              <div style={{ display: 'grid', gridTemplateColumns: '2fr repeat(6, 1fr)', gap: 12 }}>
                <Input placeholder="Search issues" value="loading" readOnly />
                <Input value="All statuses" readOnly />
                <Input value="Priority" readOnly />
                <Input value="Assignee" readOnly />
                <Input value="Project" readOnly />
                <Input value="Sprint" readOnly />
                <Input value="Label" readOnly />
                <Button variant="ghost">Reset</Button>
              </div>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(5, minmax(220px, 1fr))', gap: 16, alignItems: 'start' }}>
              {['Backlog', 'Todo', 'In Progress', 'In Review', 'Done'].map((column) => (
                <section key={column} className="panel" style={{ padding: 14, minHeight: 280 }}>
                  <div className="label" style={{ fontSize: 18, marginBottom: 12 }}>
                    {column} (1)
                  </div>
                  <div className="panel-soft" style={{ padding: 12 }}>
                    Example card
                  </div>
                </section>
              ))}
            </div>
          </div>
        }
      >
        <div className="panel" style={{ padding: 18, marginBottom: 20 }}>
          <div style={{ display: 'grid', gridTemplateColumns: '2fr repeat(6, 1fr)', gap: 12 }}>
            <Input placeholder="Search issues" value={search} onChange={(event) => setSearch(event.target.value)} />
            <Select value={status} onChange={(event) => setStatus(event.target.value)}>
              <option value="">All statuses</option>
              <option value="backlog">Backlog</option>
              <option value="todo">Todo</option>
              <option value="in_progress">In Progress</option>
              <option value="in_review">In Review</option>
              <option value="done">Done</option>
              <option value="cancelled">Cancelled</option>
            </Select>
            <Select value={priority} onChange={(event) => setPriority(event.target.value)}>
              <option value="">All priorities</option>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
              <option value="urgent">Urgent</option>
            </Select>
            <Select value={assigneeId} onChange={(event) => setAssigneeId(event.target.value)}>
              <option value="">All assignees</option>
              {users.data?.map((user) => (
                <option key={user.id} value={user.id}>
                  {user.name}
                </option>
              ))}
            </Select>
            <Select value={projectId} onChange={(event) => setProjectId(event.target.value)}>
              <option value="">All projects</option>
              {projects.data?.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.name}
                </option>
              ))}
            </Select>
            <Select value={sprintId} onChange={(event) => setSprintId(event.target.value)}>
              <option value="">All sprints</option>
              {sprints.data?.map((sprint) => (
                <option key={sprint.id} value={sprint.id}>
                  {sprint.name}
                </option>
              ))}
            </Select>
            <Select value={labelId} onChange={(event) => setLabelId(event.target.value)}>
              <option value="">All labels</option>
              {labels.data?.map((label) => (
                <option key={label.id} value={label.id}>
                  {label.name}
                </option>
              ))}
            </Select>
            <Button
              variant="ghost"
              onClick={() => {
                setSearch('');
                setStatus('');
                setPriority('');
                setAssigneeId('');
                setProjectId('');
                setSprintId('');
                setLabelId('');
              }}
            >
              Reset
            </Button>
          </div>
        </div>
        {issues.isFetching && !issues.isLoading ? <div style={{ color: 'var(--text-secondary)', marginBottom: 12 }}>Updating board...</div> : null}
        {issues.isError ? <ErrorBanner message={(issues.error as Error).message} /> : null}
        {issues.data ? (
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(5, minmax(220px, 1fr))',
              gap: 16,
              alignItems: 'start',
              overflowX: 'auto',
            }}
          >
            {columns.map((column) => (
              <section
                key={column}
                className="panel"
                onDragOver={(event) => event.preventDefault()}
                onDrop={(event) => {
                  if (movingIssueId) return;
                  const issueId = event.dataTransfer.getData('text/plain');
                  if (issueId) {
                    moveIssue(issueId, column);
                  }
                }}
                style={{ padding: 14, minHeight: 420 }}
              >
                <div className="label" style={{ fontSize: 18, marginBottom: 12 }}>
                  {column.replace('_', ' ')} ({grouped[column].length})
                </div>
                <div style={{ display: 'grid', gap: 12 }}>
                  {grouped[column].map((issue) => (
                    <div
                      key={issue.id}
                      draggable={!movingIssueId}
                      onDragStart={(event) => {
                        if (movingIssueId) {
                          event.preventDefault();
                          return;
                        }
                        event.dataTransfer.setData('text/plain', issue.id);
                      }}
                      style={{ opacity: movingIssueId === issue.id ? 0.6 : 1 }}
                    >
                      <IssueCard issue={issue} compact />
                    </div>
                  ))}
                </div>
              </section>
            ))}
          </div>
        ) : null}
        {issues.data && issues.data.items.length === 0 ? (
          <EmptyState title="No issues found" description="No issues match the current board filters." />
        ) : null}
      </Skeleton>
    </div>
  );
}

function groupByStatus(items: IssueSummary[]): Record<BoardColumn, IssueSummary[]> {
  return {
    backlog: items.filter((issue) => issue.status === 'backlog'),
    todo: items.filter((issue) => issue.status === 'todo'),
    in_progress: items.filter((issue) => issue.status === 'in_progress'),
    in_review: items.filter((issue) => issue.status === 'in_review'),
    done: items.filter((issue) => issue.status === 'done'),
    cancelled: items.filter((issue) => issue.status === 'cancelled'),
  };
}
