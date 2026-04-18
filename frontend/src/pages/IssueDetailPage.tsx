import { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Badge } from 'components/common/Badge';
import { Button } from 'components/common/Button';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { PageHeader } from 'components/common/PageHeader';
import { Select } from 'components/common/Select';
import { Spinner } from 'components/common/Spinner';
import { useDebouncedValue } from 'hooks/useDebouncedValue';
import { useArchiveIssue, useIssueDetail, useUpdateIssue } from 'features/issues/issuesQueries';
import { useLabelsSelector, useProjectsSelector, useSprintsSelector, useUsersSelector } from 'features/issues/selectorsQueries';
import { useUIStore } from 'store/uiStore';
import { ApiError } from 'types/api';
import { getBannerErrorMessage, parseUiError } from 'utils/errorPresentation';
import { formatDate, relativeTime, titleCase } from 'utils/format';

export function IssueDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const pushToast = useUIStore((state) => state.pushToast);
  const [includeArchived, setIncludeArchived] = useState(false);
  const issue = useIssueDetail(id, includeArchived);
  const projects = useProjectsSelector();
  const users = useUsersSelector();
  const labels = useLabelsSelector();
  const [projectId, setProjectId] = useState('');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState('');
  const [priority, setPriority] = useState('');
  const [assigneeId, setAssigneeId] = useState('');
  const [sprintId, setSprintId] = useState('');
  const [labelIds, setLabelIds] = useState<string[]>([]);

  const sprints = useSprintsSelector(projectId || undefined);
  const updateIssue = useUpdateIssue(id ?? '');
  const archiveIssue = useArchiveIssue(id ?? '');

  useEffect(() => {
    if (!issue.data) return;
    setProjectId(issue.data.project_id);
    setTitle(issue.data.title);
    setDescription(issue.data.description ?? '');
    setStatus(issue.data.status);
    setPriority(issue.data.priority);
    setAssigneeId(issue.data.assignee_id ?? '');
    setSprintId(issue.data.sprint_id ?? '');
    setLabelIds(issue.data.labels.map((label) => label.id));
  }, [issue.data]);

  const debouncedTitle = useDebouncedValue(title, 1000);
  const debouncedDescription = useDebouncedValue(description, 1000);

  useEffect(() => {
    if (!issue.data) return;
    if (debouncedTitle !== issue.data.title) {
      updateIssue.mutate(
        { title: debouncedTitle },
        {
          onError: (error) => {
            setTitle(issue.data?.title ?? '');
            pushToast({ tone: 'error', message: parseUiError(error, 'Failed to update title.').message });
          },
        },
      );
    }
  }, [debouncedTitle, issue.data, pushToast, updateIssue]);

  useEffect(() => {
    if (!issue.data) return;
    if (debouncedDescription !== (issue.data.description ?? '')) {
      updateIssue.mutate(
        { description: debouncedDescription || null },
        {
          onError: (error) => {
            setDescription(issue.data?.description ?? '');
            pushToast({ tone: 'error', message: parseUiError(error, 'Failed to update description.').message });
          },
        },
      );
    }
  }, [debouncedDescription, issue.data, pushToast, updateIssue]);

  const saving = updateIssue.isPending;
  const issueError = issue.error instanceof ApiError ? issue.error : null;

  function updateWithFeedback<T>(payload: T, rollback: () => void, fallbackMessage: string) {
    updateIssue.mutate(payload as never, {
      onError: (error) => {
        rollback();
        pushToast({ tone: 'error', message: parseUiError(error, fallbackMessage).message });
      },
    });
  }

  const sidebar = useMemo(
    () => (
      <div className="panel" style={{ padding: 20, display: 'grid', gap: 16 }}>
        <div>
          <div className="label" style={{ marginBottom: 8 }}>
            Status
          </div>
          <Select
            value={status}
            onChange={(event) => {
              const previous = status;
              const next = event.target.value;
              setStatus(next);
              updateWithFeedback({ status: next as never }, () => setStatus(previous), 'Failed to update status.');
            }}
          >
            <option value="backlog">Backlog</option>
            <option value="todo">Todo</option>
            <option value="in_progress">In Progress</option>
            <option value="in_review">In Review</option>
            <option value="done">Done</option>
            <option value="cancelled">Cancelled</option>
          </Select>
        </div>
        <div>
          <div className="label" style={{ marginBottom: 8 }}>
            Priority
          </div>
          <Select
            value={priority}
            onChange={(event) => {
              const previous = priority;
              const next = event.target.value;
              setPriority(next);
              updateWithFeedback({ priority: next as never }, () => setPriority(previous), 'Failed to update priority.');
            }}
          >
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="urgent">Urgent</option>
          </Select>
        </div>
        <div>
          <div className="label" style={{ marginBottom: 8 }}>
            Project
          </div>
          <Select
            value={projectId}
            onChange={(event) => {
              const previousProject = projectId;
              const previousSprint = sprintId;
              const next = event.target.value;
              setProjectId(next);
              setSprintId('');
              updateWithFeedback(
                { project_id: next, sprint_id: null },
                () => {
                  setProjectId(previousProject);
                  setSprintId(previousSprint);
                },
                'Failed to update project.',
              );
            }}
          >
            {projects.data?.map((project) => (
              <option key={project.id} value={project.id}>
                {project.name}
              </option>
            ))}
          </Select>
        </div>
        <div>
          <div className="label" style={{ marginBottom: 8 }}>
            Assignee
          </div>
          <Select
            value={assigneeId}
            onChange={(event) => {
              const previous = assigneeId;
              const next = event.target.value;
              setAssigneeId(next);
              updateWithFeedback({ assignee_id: next || null }, () => setAssigneeId(previous), 'Failed to update assignee.');
            }}
          >
            <option value="">Unassigned</option>
            {users.data?.map((user) => (
              <option key={user.id} value={user.id}>
                {user.name}
              </option>
            ))}
          </Select>
        </div>
        <div>
          <div className="label" style={{ marginBottom: 8 }}>
            Sprint
          </div>
          <Select
            value={sprintId}
            onChange={(event) => {
              const previous = sprintId;
              const next = event.target.value;
              setSprintId(next);
              updateWithFeedback({ sprint_id: next || null }, () => setSprintId(previous), 'Failed to update sprint.');
            }}
          >
            <option value="">No sprint</option>
            {sprints.data?.map((sprint) => (
              <option key={sprint.id} value={sprint.id}>
                {sprint.name}
              </option>
            ))}
          </Select>
        </div>
        <div>
          <div className="label" style={{ marginBottom: 8 }}>
            Labels
          </div>
          <Select
            multiple
            value={labelIds}
            onChange={(event) => {
              const previous = labelIds;
              const next = Array.from(event.target.selectedOptions).map((option) => option.value);
              setLabelIds(next);
              updateWithFeedback({ label_ids: next }, () => setLabelIds(previous), 'Failed to update labels.');
            }}
            style={{ minHeight: 120 }}
          >
            {labels.data?.map((label) => (
              <option key={label.id} value={label.id}>
                {label.name}
              </option>
            ))}
          </Select>
        </div>
        <Button
          variant="danger"
          disabled={archiveIssue.isPending}
          onClick={() => {
            if (!window.confirm('Archive this issue? You can restore it later from archived views.')) {
              return;
            }
            archiveIssue.mutate(undefined, {
              onSuccess: () => {
                pushToast({ tone: 'success', message: 'Issue archived.' });
                navigate('/issues');
              },
              onError: (error) => {
                pushToast({ tone: 'error', message: parseUiError(error, 'Failed to archive issue.').message });
              },
            });
          }}
        >
          Archive Issue
        </Button>
      </div>
    ),
    [
      archiveIssue,
      assigneeId,
      labelIds,
      labels.data,
      navigate,
      priority,
      projectId,
      projects.data,
      pushToast,
      sprintId,
      sprints.data,
      status,
      updateIssue,
      users.data,
    ],
  );

  return (
    <div>
      <PageHeader title="Issue Detail" subtitle="Comments are intentionally omitted; the activity feed is read-only." />
      {issue.isLoading ? <Spinner label="Loading issue" /> : null}
      {issue.isError ? <ErrorBanner message={getBannerErrorMessage(issue.error, 'Unable to load issue details right now.')} /> : null}
      {issueError?.status === 404 && !includeArchived ? (
        <div style={{ marginBottom: 16 }}>
          <Button variant="ghost" onClick={() => setIncludeArchived(true)}>
            Try Loading Archived Issue
          </Button>
        </div>
      ) : null}
      {!issue.isLoading && !issue.data && !issue.isError ? (
        <EmptyState title="Issue not found" description="The requested issue could not be loaded." />
      ) : null}
      {issue.data ? (
        <div className="two-col">
          <section style={{ display: 'grid', gap: 20 }}>
            <div className="panel" style={{ padding: 20 }}>
              <div style={{ display: 'flex', gap: 12, alignItems: 'center', marginBottom: 14 }}>
                <Badge tone="info">{issue.data.identifier}</Badge>
                <Badge>{titleCase(issue.data.status)}</Badge>
                <Badge tone="accent">{issue.data.priority}</Badge>
                {saving ? <span style={{ color: 'var(--text-secondary)' }}>Saving...</span> : null}
              </div>
              <Input value={title} onChange={(event) => setTitle(event.target.value)} style={{ fontSize: 28, fontWeight: 700 }} />
              <textarea
                value={description}
                onChange={(event) => setDescription(event.target.value)}
                rows={12}
                style={{
                  width: '100%',
                  marginTop: 16,
                  padding: 16,
                  borderRadius: 12,
                  border: '2px solid var(--border-strong)',
                  background: 'var(--bg-elevated)',
                  color: 'var(--text-primary)',
                }}
              />
            </div>
            <div className="panel" style={{ padding: 20 }}>
              <div className="label" style={{ fontSize: 24, marginBottom: 14 }}>
                Activity History
              </div>
              <div style={{ display: 'grid', gap: 12 }}>
                {issue.data.activities.length === 0 ? (
                  <EmptyState title="No activity" description="This issue has no recorded changes yet." />
                ) : (
                  issue.data.activities.map((activity) => (
                    <div key={activity.id} className="panel-soft" style={{ padding: 14 }}>
                      <div style={{ fontWeight: 700 }}>{activity.user?.name ?? 'Unknown user'}</div>
                      <div style={{ color: 'var(--text-secondary)', marginTop: 4 }}>
                        {titleCase(activity.action)}
                        {activity.field_name ? ` · ${titleCase(activity.field_name)}` : ''} · {relativeTime(activity.created_at)}
                      </div>
                      <div style={{ marginTop: 8, fontSize: 14 }}>
                        {activity.old_value || activity.new_value
                          ? `${activity.old_value ?? 'empty'} → ${activity.new_value ?? 'empty'}`
                          : `Recorded at ${formatDate(activity.created_at)}`}
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </section>
          <aside>{sidebar}</aside>
        </div>
      ) : null}
    </div>
  );
}
