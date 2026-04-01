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
import { ApiError } from 'types/api';

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

  const users = useUsersSelector();
  const projects = useProjectsSelector();
  const sprints = useSprintsSelector(projectId || undefined);
  const labels = useLabelsSelector();

  const mutation = useMutation({
    mutationFn: () =>
      issuesApi.create({
        title,
        description: description || null,
        project_id: projectId,
        status: status as 'backlog',
        priority: priority as 'medium',
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
      navigate(`/issues/${response.data.id}`);
    },
    onError: (error) => {
      pushToast({ tone: 'error', message: error instanceof Error ? error.message : 'Failed to create issue.' });
    },
  });

  const error = mutation.error instanceof ApiError ? mutation.error.message : null;

  const selectedProjectSprints = useMemo(() => sprints.data ?? [], [sprints.data]);

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    mutation.mutate();
  }

  return (
    <Modal open={open} title="Create Issue" onClose={close}>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: 16 }}>
        {error ? <ErrorBanner message={error} /> : null}
        <div>
          <div className="label" style={{ marginBottom: 8 }}>
            Title
          </div>
          <Input value={title} onChange={(event) => setTitle(event.target.value)} placeholder="Build authentication flow" />
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
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, minmax(0, 1fr))', gap: 16 }}>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Project
            </div>
            <Select value={projectId} onChange={(event) => setProjectId(event.target.value)} required>
              <option value="">Select project</option>
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
            <Select value={assigneeId} onChange={(event) => setAssigneeId(event.target.value)}>
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
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Sprint
            </div>
            <Select value={sprintId} onChange={(event) => setSprintId(event.target.value)}>
              <option value="">No sprint</option>
              {selectedProjectSprints.map((sprint) => (
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
