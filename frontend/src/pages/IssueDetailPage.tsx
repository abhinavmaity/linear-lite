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
import { formatDate, relativeTime, titleCase } from 'utils/format';

export function IssueDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const pushToast = useUIStore((state) => state.pushToast);
  const issue = useIssueDetail(id);
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
      updateIssue.mutate({ title: debouncedTitle });
    }
  }, [debouncedTitle, issue.data, updateIssue]);

  useEffect(() => {
    if (!issue.data) return;
    if (debouncedDescription !== (issue.data.description ?? '')) {
      updateIssue.mutate({ description: debouncedDescription || null });
    }
  }, [debouncedDescription, issue.data, updateIssue]);

  const saving = updateIssue.isPending;

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
              const next = event.target.value;
              setStatus(next);
              updateIssue.mutate({ status: next as never });
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
              const next = event.target.value;
              setPriority(next);
              updateIssue.mutate({ priority: next as never });
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
              const next = event.target.value;
              setProjectId(next);
              updateIssue.mutate({ project_id: next, sprint_id: null });
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
              const next = event.target.value;
              setAssigneeId(next);
              updateIssue.mutate({ assignee_id: next || null });
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
              const next = event.target.value;
              setSprintId(next);
              updateIssue.mutate({ sprint_id: next || null });
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
              const next = Array.from(event.target.selectedOptions).map((option) => option.value);
              setLabelIds(next);
              updateIssue.mutate({ label_ids: next });
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
          onClick={() =>
            archiveIssue.mutate(undefined, {
              onSuccess: () => {
                pushToast({ tone: 'success', message: 'Issue archived.' });
                navigate('/issues');
              },
              onError: (error) => {
                pushToast({ tone: 'error', message: error instanceof Error ? error.message : 'Failed to archive issue.' });
              },
            })
          }
        >
          Archive Issue
        </Button>
      </div>
    ),
    [archiveIssue, assigneeId, labelIds, labels.data, navigate, priority, projectId, projects.data, pushToast, sprintId, sprints.data, status, updateIssue, users.data],
  );

  return (
    <div>
      <PageHeader title="Issue Detail" subtitle="Comments are intentionally omitted; the activity feed is read-only." />
      {issue.isLoading ? <Spinner label="Loading issue" /> : null}
      {issue.isError ? <ErrorBanner message={(issue.error as Error).message} /> : null}
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
                      <div style={{ fontWeight: 700 }}>{activity.user.name}</div>
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
