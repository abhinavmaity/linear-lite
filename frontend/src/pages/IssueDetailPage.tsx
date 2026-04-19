import { Fragment, ReactNode, useEffect, useMemo, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { Badge } from 'components/common/Badge';
import { Button } from 'components/common/Button';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { PageHeader } from 'components/common/PageHeader';
import { Select } from 'components/common/Select';
import { Spinner } from 'components/common/Spinner';
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
  const [contentDirty, setContentDirty] = useState(false);
  const [contentTab, setContentTab] = useState<'write' | 'preview'>('write');

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
    setContentDirty(false);
  }, [issue.data]);

  const saving = updateIssue.isPending;
  const issueError = issue.error instanceof ApiError ? issue.error : null;

  function saveContent() {
    if (!issue.data) return;
    const trimmedTitle = title.trim();
    if (!trimmedTitle) {
      pushToast({ tone: 'error', message: 'Title is required.' });
      return;
    }

    const payload: Record<string, unknown> = {};
    if (trimmedTitle !== issue.data.title) {
      payload.title = trimmedTitle;
    }
    const normalizedDescription = description.trim() ? description : null;
    const currentDescription = issue.data.description ?? null;
    if (normalizedDescription !== currentDescription) {
      payload.description = normalizedDescription;
    }

    if (Object.keys(payload).length === 0) {
      setContentDirty(false);
      return;
    }

    updateIssue.mutate(payload, {
      onSuccess: () => {
        setContentDirty(false);
        pushToast({ tone: 'success', message: 'Issue content saved.' });
      },
      onError: (error) => {
        setTitle(issue.data?.title ?? '');
        setDescription(issue.data?.description ?? '');
        setContentDirty(false);
        pushToast({ tone: 'error', message: parseUiError(error, 'Failed to save issue content.').message });
      },
    });
  }

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
          {!sprints.isLoading && projectId && (sprints.data?.length ?? 0) === 0 ? (
            <div style={{ marginTop: 6, color: 'var(--text-secondary)' }}>
              No sprints in this project yet. <Link to="/sprints">Create one from Sprints.</Link>
            </div>
          ) : null}
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
      sprints.isLoading,
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
              <div style={{ display: 'flex', gap: 12, alignItems: 'center', marginBottom: 14, flexWrap: 'wrap' }}>
                <Badge tone="info">{issue.data.identifier}</Badge>
                <Badge>{titleCase(issue.data.status)}</Badge>
                <Badge tone="accent">{issue.data.priority}</Badge>
                {saving ? <span style={{ color: 'var(--text-secondary)' }}>Saving...</span> : null}
                {contentDirty ? <span style={{ color: 'var(--text-secondary)' }}>Unsaved changes</span> : null}
                <Button type="button" variant="ghost" onClick={() => setContentTab((prev) => (prev === 'write' ? 'preview' : 'write'))}>
                  {contentTab === 'write' ? 'Preview Markdown' : 'Edit Markdown'}
                </Button>
                <Button type="button" onClick={saveContent} disabled={!contentDirty || saving}>
                  Save Content
                </Button>
              </div>
              <Input
                value={title}
                onChange={(event) => {
                  setTitle(event.target.value);
                  setContentDirty(true);
                }}
                style={{ fontSize: 28, fontWeight: 700 }}
              />
              {contentTab === 'write' ? (
                <textarea
                  value={description}
                  onChange={(event) => {
                    setDescription(event.target.value);
                    setContentDirty(true);
                  }}
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
              ) : (
                <div className="panel-soft" style={{ marginTop: 16, padding: 16 }}>
                  <MarkdownPreview value={description} />
                </div>
              )}
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
                          ? `${activity.old_value ?? 'empty'} -> ${activity.new_value ?? 'empty'}`
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

function MarkdownPreview({ value }: { value: string }) {
  const content = value.trim();
  if (!content) {
    return <div style={{ color: 'var(--text-secondary)' }}>No description yet.</div>;
  }

  const lines = content.split('\n');
  const nodes: ReactNode[] = [];
  let bulletItems: string[] = [];

  const flushBullets = () => {
    if (bulletItems.length === 0) return;
    nodes.push(
      <ul key={`list-${nodes.length}`} style={{ margin: '8px 0 12px 20px' }}>
        {bulletItems.map((item, idx) => (
          <li key={idx}>{renderInlineMarkdown(item)}</li>
        ))}
      </ul>,
    );
    bulletItems = [];
  };

  for (let i = 0; i < lines.length; i += 1) {
    const line = lines[i].trim();
    if (!line) {
      flushBullets();
      nodes.push(<div key={`spacer-${i}`} style={{ height: 8 }} />);
      continue;
    }

    if (line.startsWith('- ')) {
      bulletItems.push(line.slice(2));
      continue;
    }

    flushBullets();
    if (line.startsWith('### ')) {
      nodes.push(
        <h3 key={`h3-${i}`} style={{ margin: '8px 0 6px' }}>
          {renderInlineMarkdown(line.slice(4))}
        </h3>,
      );
      continue;
    }
    if (line.startsWith('## ')) {
      nodes.push(
        <h2 key={`h2-${i}`} style={{ margin: '10px 0 8px' }}>
          {renderInlineMarkdown(line.slice(3))}
        </h2>,
      );
      continue;
    }
    if (line.startsWith('# ')) {
      nodes.push(
        <h1 key={`h1-${i}`} style={{ margin: '12px 0 10px', fontSize: 24 }}>
          {renderInlineMarkdown(line.slice(2))}
        </h1>,
      );
      continue;
    }

    nodes.push(
      <p key={`p-${i}`} style={{ margin: '6px 0' }}>
        {renderInlineMarkdown(line)}
      </p>,
    );
  }

  flushBullets();
  return <div>{nodes}</div>;
}

function renderInlineMarkdown(text: string): ReactNode {
  const parts: ReactNode[] = [];
  const regex = /(`[^`]+`|\*\*[^*]+\*\*|\*[^*]+\*)/g;
  let cursor = 0;
  let match: RegExpExecArray | null;

  while ((match = regex.exec(text)) !== null) {
    if (match.index > cursor) {
      parts.push(<Fragment key={`txt-${cursor}`}>{text.slice(cursor, match.index)}</Fragment>);
    }
    const token = match[0];
    if (token.startsWith('`') && token.endsWith('`')) {
      parts.push(
        <code key={`code-${match.index}`} style={{ background: 'var(--bg-elevated)', padding: '1px 4px', borderRadius: 4 }}>
          {token.slice(1, -1)}
        </code>,
      );
    } else if (token.startsWith('**') && token.endsWith('**')) {
      parts.push(<strong key={`strong-${match.index}`}>{token.slice(2, -2)}</strong>);
    } else if (token.startsWith('*') && token.endsWith('*')) {
      parts.push(<em key={`em-${match.index}`}>{token.slice(1, -1)}</em>);
    }
    cursor = match.index + token.length;
  }

  if (cursor < text.length) {
    parts.push(<Fragment key={`tail-${cursor}`}>{text.slice(cursor)}</Fragment>);
  }
  return <>{parts}</>;
}
