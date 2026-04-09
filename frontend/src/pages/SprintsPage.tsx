import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { FormEvent, useMemo, useState } from 'react';
import { Skeleton } from 'boneyard-js/react';
import { Button } from 'components/common/Button';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { PageHeader } from 'components/common/PageHeader';
import { Select } from 'components/common/Select';
import { Spinner } from 'components/common/Spinner';
import { projectsApi } from 'services/projectsApi';
import { sprintsApi, SprintCreateInput, SprintUpdateInput } from 'services/sprintsApi';
import { useUIStore } from 'store/uiStore';
import { ApiError } from 'types/api';
import { formatDate, titleCase } from 'utils/format';

export function SprintsPage() {
  const queryClient = useQueryClient();
  const pushToast = useUIStore((state) => state.pushToast);
  const [projectFilterId, setProjectFilterId] = useState('');
  const [createProjectId, setCreateProjectId] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [sortBy, setSortBy] = useState<'start_date' | 'name' | 'end_date' | 'created_at'>('start_date');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [status, setStatus] = useState<'planned' | 'active' | 'completed'>('planned');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [editStartDate, setEditStartDate] = useState('');
  const [editEndDate, setEditEndDate] = useState('');
  const [editStatus, setEditStatus] = useState<'planned' | 'active' | 'completed'>('planned');
  const [actionError, setActionError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string> | null>(null);

  const projects = useQuery({
    queryKey: ['projects', 'sprint-selector'],
    queryFn: () => projectsApi.list({ limit: 100, sort_by: 'name', sort_order: 'asc' }).then((response) => response.items),
  });

  const sprints = useQuery({
    queryKey: ['sprints', 'page', projectFilterId, statusFilter, sortBy, sortOrder],
    queryFn: () =>
      sprintsApi
        .list({
          limit: 50,
          project_id: projectFilterId || undefined,
          status: (statusFilter || undefined) as 'planned' | 'active' | 'completed' | undefined,
          sort_by: sortBy,
          sort_order: sortOrder,
        })
        .then((response) => response.items),
  });

  const createSprint = useMutation({
    mutationFn: (payload: SprintCreateInput) => sprintsApi.create(payload),
    onSuccess: () => {
      setName('');
      setDescription('');
      setStartDate('');
      setEndDate('');
      setStatus('planned');
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Sprint created.' });
      queryClient.invalidateQueries({ queryKey: ['sprints'] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to create sprint.');
      setActionError(parsed.message);
      setFieldErrors(parsed.fields);
    },
  });

  const updateSprint = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: SprintUpdateInput }) => sprintsApi.update(id, payload),
    onSuccess: () => {
      setEditingId(null);
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Sprint updated.' });
      queryClient.invalidateQueries({ queryKey: ['sprints'] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to update sprint.');
      setActionError(parsed.message);
      setFieldErrors(parsed.fields);
    },
  });

  const deleteSprint = useMutation({
    mutationFn: (id: string) => sprintsApi.delete(id),
    onSuccess: () => {
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Sprint deleted.' });
      queryClient.invalidateQueries({ queryKey: ['sprints'] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to delete sprint.');
      setActionError(parsed.message);
      setFieldErrors(parsed.fields);
    },
  });

  const projectNameById = useMemo(() => new Map((projects.data ?? []).map((p) => [p.id, p.name])), [projects.data]);

  function handleCreate(event: FormEvent) {
    event.preventDefault();
    setActionError(null);
    setFieldErrors(null);
    createSprint.mutate({
      name,
      description: description || null,
      project_id: createProjectId,
      start_date: startDate,
      end_date: endDate,
      status,
    });
  }

  function startEdit(sprintId: string, sprintName: string, sprintDescription: string | null, sprintStartDate: string, sprintEndDate: string, sprintStatus: 'planned' | 'active' | 'completed') {
    setEditingId(sprintId);
    setEditName(sprintName);
    setEditDescription(sprintDescription ?? '');
    setEditStartDate(sprintStartDate);
    setEditEndDate(sprintEndDate);
    setEditStatus(sprintStatus);
    setActionError(null);
    setFieldErrors(null);
  }

  function handleUpdate(event: FormEvent, sprintId: string) {
    event.preventDefault();
    setActionError(null);
    setFieldErrors(null);
    updateSprint.mutate({
      id: sprintId,
      payload: {
        name: editName,
        description: editDescription || null,
        start_date: editStartDate,
        end_date: editEndDate,
        status: editStatus,
      },
    });
  }

  return (
    <div>
      <PageHeader title="Sprints" subtitle="Plan and track sprint windows across active projects." />
      <div className="panel" style={{ padding: 18, marginBottom: 18 }}>
        <form onSubmit={handleCreate} style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 1fr 1fr 1fr auto', gap: 12, alignItems: 'end' }}>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Sprint name
            </div>
            <Input value={name} onChange={(event) => setName(event.target.value)} required />
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Project
            </div>
            <Select value={createProjectId} onChange={(event) => setCreateProjectId(event.target.value)} required>
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
              Start
            </div>
            <Input type="date" value={startDate} onChange={(event) => setStartDate(event.target.value)} required />
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              End
            </div>
            <Input type="date" value={endDate} onChange={(event) => setEndDate(event.target.value)} required />
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Status
            </div>
            <Select value={status} onChange={(event) => setStatus(event.target.value as 'planned' | 'active' | 'completed')}>
              <option value="planned">Planned</option>
              <option value="active">Active</option>
              <option value="completed">Completed</option>
            </Select>
          </div>
          <Button type="submit" disabled={createSprint.isPending}>
            {createSprint.isPending ? 'Creating' : 'Create'}
          </Button>
        </form>
        <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 1fr 1fr', gap: 12, marginTop: 12 }}>
          <Input value={description} onChange={(event) => setDescription(event.target.value)} placeholder="Description for create flow (optional)" />
          <Select value={projectFilterId} onChange={(event) => setProjectFilterId(event.target.value)}>
            <option value="">Filter: all projects</option>
            {projects.data?.map((project) => (
              <option key={project.id} value={project.id}>
                {project.name}
              </option>
            ))}
          </Select>
          <Select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value)}>
            <option value="">Filter: all statuses</option>
            <option value="planned">Planned</option>
            <option value="active">Active</option>
            <option value="completed">Completed</option>
          </Select>
          <Select value={sortBy} onChange={(event) => setSortBy(event.target.value as 'start_date' | 'name' | 'end_date' | 'created_at')}>
            <option value="start_date">Sort: start date</option>
            <option value="name">Sort: name</option>
            <option value="end_date">Sort: end date</option>
            <option value="created_at">Sort: created</option>
          </Select>
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: 12, marginTop: 12 }}>
          <Select value={sortOrder} onChange={(event) => setSortOrder(event.target.value as 'asc' | 'desc')}>
            <option value="desc">Order: desc</option>
            <option value="asc">Order: asc</option>
          </Select>
        </div>
      </div>
      {actionError ? <ErrorBanner message={actionError} /> : null}
      <Skeleton
        name="sprints-page"
        loading={sprints.isLoading}
        fallback={<Spinner label="Loading sprints" />}
        fixture={
          <div style={{ display: 'grid', gap: 16 }}>
            <article className="panel" style={{ padding: 18 }}>
              <h3 style={{ marginTop: 0 }}>Sprint Alpha</h3>
              <p style={{ color: 'var(--text-secondary)' }}>Active · Apr 01 to Apr 14</p>
            </article>
          </div>
        }
      >
        {sprints.isFetching && !sprints.isLoading ? <div style={{ color: 'var(--text-secondary)', marginBottom: 12 }}>Refreshing sprints...</div> : null}
        {sprints.isError ? <ErrorBanner message={(sprints.error as Error).message} /> : null}
        {sprints.data?.length ? (
          <div style={{ display: 'grid', gap: 16 }}>
            {sprints.data.map((sprint) => (
              <article key={sprint.id} className="panel" style={{ padding: 18 }}>
                {editingId === sprint.id ? (
                  <form onSubmit={(event) => handleUpdate(event, sprint.id)} style={{ display: 'grid', gap: 10 }}>
                    <Input value={editName} onChange={(event) => setEditName(event.target.value)} required />
                    <Input value={editDescription} onChange={(event) => setEditDescription(event.target.value)} placeholder="Optional description" />
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 8 }}>
                      <Input type="date" value={editStartDate} onChange={(event) => setEditStartDate(event.target.value)} required />
                      <Input type="date" value={editEndDate} onChange={(event) => setEditEndDate(event.target.value)} required />
                      <Select value={editStatus} onChange={(event) => setEditStatus(event.target.value as 'planned' | 'active' | 'completed')}>
                        <option value="planned">Planned</option>
                        <option value="active">Active</option>
                        <option value="completed">Completed</option>
                      </Select>
                    </div>
                    {fieldErrors?.name ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.name}</div> : null}
                    {fieldErrors?.start_date ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.start_date}</div> : null}
                    {fieldErrors?.end_date ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.end_date}</div> : null}
                    {fieldErrors?.status ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.status}</div> : null}
                    <div style={{ display: 'flex', gap: 10 }}>
                      <Button type="submit" disabled={updateSprint.isPending}>
                        {updateSprint.isPending ? 'Saving' : 'Save'}
                      </Button>
                      <Button type="button" variant="ghost" onClick={() => setEditingId(null)}>
                        Cancel
                      </Button>
                    </div>
                  </form>
                ) : (
                  <>
                    <h3 style={{ marginTop: 0 }}>{sprint.name}</h3>
                    <p style={{ color: 'var(--text-secondary)' }}>
                      {titleCase(sprint.status)} · {formatDate(sprint.start_date)} to {formatDate(sprint.end_date)}
                    </p>
                    <p style={{ color: 'var(--text-secondary)' }}>Project: {projectNameById.get(sprint.project_id) ?? sprint.project_id}</p>
                    <p style={{ color: 'var(--text-secondary)' }}>
                      Issues: total {sprint.issue_counts.total}, in progress {sprint.issue_counts.in_progress}, done {sprint.issue_counts.done}
                    </p>
                    <div style={{ display: 'flex', gap: 10 }}>
                      <Button
                        type="button"
                        variant="ghost"
                        onClick={() => startEdit(sprint.id, sprint.name, sprint.description, sprint.start_date, sprint.end_date, sprint.status)}
                      >
                        Edit
                      </Button>
                      <Button
                        type="button"
                        variant="danger"
                        disabled={deleteSprint.isPending}
                        onClick={() => {
                          if (!window.confirm(`Delete sprint "${sprint.name}"?`)) return;
                          setActionError(null);
                          setFieldErrors(null);
                          deleteSprint.mutate(sprint.id);
                        }}
                      >
                        Delete
                      </Button>
                    </div>
                  </>
                )}
              </article>
            ))}
          </div>
        ) : null}
        {sprints.data && sprints.data.length === 0 ? (
          <EmptyState title="No sprints" description="Sprint data will appear here when available." />
        ) : null}
      </Skeleton>
    </div>
  );
}

function parseApiError(error: unknown, fallback: string) {
  if (error instanceof ApiError) {
    return { message: error.message || fallback, fields: error.fields ?? null };
  }
  return { message: fallback, fields: null };
}
