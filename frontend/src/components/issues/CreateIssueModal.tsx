import { useMutation, useQueryClient } from '@tanstack/react-query';
import { FormEvent, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from 'components/common/Button';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { Modal } from 'components/common/Modal';
import { Select } from 'components/common/Select';
import { useLabelsSelector, useProjectsSelector, useSprintsSelector, useUsersSelector } from 'features/issues/selectorsQueries';
import { issuesApi } from 'services/issuesApi';
import { useUIStore } from 'store/uiStore';
import { parseUiError } from 'utils/errorPresentation';
import { IssuePriority, IssueStatus } from 'types/domain';

export function CreateIssueModal() {
  const open = useUIStore((state) => state.createIssueOpen);
  const close = useUIStore((state) => state.closeCreateIssue);
  const pushToast = useUIStore((state) => state.pushToast);
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [projectId, setProjectId] = useState('');
  const [status, setStatus] = useState('backlog');
  const [priority, setPriority] = useState('medium');
  const [assigneeId, setAssigneeId] = useState('');
  const [sprintId, setSprintId] = useState('');
  const [labelIds, setLabelIds] = useState<string[]>([]);
  const [clientError, setClientError] = useState<string | null>(null);

  const users = useUsersSelector();
  const projects = useProjectsSelector();
  const sprints = useSprintsSelector(projectId || undefined);
  const labels = useLabelsSelector();

  const mutation = useMutation({
    mutationFn: () =>
      issuesApi.create({
        title: title.trim(),
        description: description || null,
        project_id: projectId,
        status: status as IssueStatus,
        priority: priority as IssuePriority,
        assignee_id: assigneeId || null,
        sprint_id: sprintId || null,
        label_ids: labelIds,
      }),
    onSuccess: (response) => {
      queryClient.invalidateQueries({ queryKey: ['issues'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
      pushToast({ tone: 'success', message: 'Issue created.' });
      close();
      setTitle('');
      setDescription('');
      setProjectId('');
      setSprintId('');
      setLabelIds([]);
      setClientError(null);
      navigate(`/issues/${response.data.id}`);
    },
    onError: (error) => {
      const parsed = parseUiError(error, 'Failed to create issue.');
      pushToast({ tone: 'error', message: parsed.message });
    },
  });

  const parsedError = mutation.error ? parseUiError(mutation.error, 'Failed to create issue.') : null;
  const error = parsedError?.summary ?? null;
  const fieldErrors = parsedError?.fields ?? undefined;

  const selectedProjectSprints = useMemo(() => sprints.data ?? [], [sprints.data]);

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (!title.trim()) {
      setClientError('Title is required.');
      return;
    }
    if (!projectId) {
      setClientError('Project is required.');
      return;
    }
    setClientError(null);
    mutation.mutate();
  }

  return (
    <Modal open={open} title="Create Issue" onClose={close}>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: 16 }}>
        {clientError ? <ErrorBanner message={clientError} /> : null}
        {error ? <ErrorBanner message={error} /> : null}
        <div>
          <div className="label" style={{ marginBottom: 8 }}>
            Title
          </div>
          <Input value={title} onChange={(event) => setTitle(event.target.value)} placeholder="Build authentication flow" />
          {fieldErrors?.title ? <div style={{ color: 'var(--text-secondary)', marginTop: 6 }}>{fieldErrors.title}</div> : null}
        </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Description
            </div>
          <textarea
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            rows={5}
            style={{
              width: '100%',
              padding: 14,
              borderRadius: 10,
              border: '2px solid var(--border-strong)',
              background: 'var(--bg-elevated)',
              color: 'var(--text-primary)',
            }}
          />
          {fieldErrors?.description ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.description}</div> : null}
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, minmax(0, 1fr))', gap: 16 }}>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Project
            </div>
            <Select
              value={projectId}
              onChange={(event) => {
                const nextProject = event.target.value;
                setProjectId(nextProject);
                setSprintId('');
              }}
              required
            >
              <option value="">Select project</option>
              {projects.data?.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.name}
                </option>
              ))}
            </Select>
            {fieldErrors?.project_id ? <div style={{ color: 'var(--text-secondary)', marginTop: 6 }}>{fieldErrors.project_id}</div> : null}
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Assignee
            </div>
            <Select value={assigneeId} onChange={(event) => setAssigneeId(event.target.value)}>
              <option value="">Unassigned</option>
              {users.data?.map((user) => (
                <option key={user.id} value={user.id}>
                  {user.name}
                </option>
              ))}
            </Select>
            {fieldErrors?.assignee_id ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.assignee_id}</div> : null}
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Status
            </div>
            <Select value={status} onChange={(event) => setStatus(event.target.value)}>
              <option value="backlog">Backlog</option>
              <option value="todo">Todo</option>
              <option value="in_progress">In Progress</option>
              <option value="in_review">In Review</option>
              <option value="done">Done</option>
              <option value="cancelled">Cancelled</option>
            </Select>
            {fieldErrors?.status ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.status}</div> : null}
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Priority
            </div>
            <Select value={priority} onChange={(event) => setPriority(event.target.value)}>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
              <option value="urgent">Urgent</option>
            </Select>
            {fieldErrors?.priority ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.priority}</div> : null}
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Sprint
            </div>
            <Select value={sprintId} onChange={(event) => setSprintId(event.target.value)} disabled={!projectId}>
              <option value="">No sprint</option>
              {selectedProjectSprints.map((sprint) => (
                <option key={sprint.id} value={sprint.id}>
                  {sprint.name}
                </option>
              ))}
            </Select>
            {!projectId ? <div style={{ color: 'var(--text-secondary)', marginTop: 6 }}>Select a project first.</div> : null}
            {fieldErrors?.sprint_id ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.sprint_id}</div> : null}
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Labels
            </div>
            <Select
              multiple
              value={labelIds}
              onChange={(event) =>
                setLabelIds(Array.from(event.target.selectedOptions).map((option) => option.value))
              }
              style={{ minHeight: 110 }}
            >
              {labels.data?.map((label) => (
                <option key={label.id} value={label.id}>
                  {label.name}
                </option>
              ))}
            </Select>
            {fieldErrors?.label_ids ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.label_ids}</div> : null}
          </div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12 }}>
          <Button type="button" variant="ghost" onClick={close}>
            Cancel
          </Button>
          <Button type="submit" disabled={!title.trim() || !projectId || mutation.isPending}>
            {mutation.isPending ? 'Creating' : 'Create Issue'}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
