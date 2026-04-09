import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { FormEvent, useMemo, useState } from 'react';
import { Button } from 'components/common/Button';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { projectsApi, ProjectCreateInput, ProjectUpdateInput } from 'services/projectsApi';
import { useUIStore } from 'store/uiStore';
import { ApiError } from 'types/api';
import { formatDate } from 'utils/format';

export function ProjectsPage() {
  const queryClient = useQueryClient();
  const pushToast = useUIStore((state) => state.pushToast);
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'created_at' | 'updated_at'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [createName, setCreateName] = useState('');
  const [createKey, setCreateKey] = useState('');
  const [createDescription, setCreateDescription] = useState('');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState('');
  const [editKey, setEditKey] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [actionError, setActionError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string> | null>(null);

  const projects = useQuery({
    queryKey: ['projects', 'page', search, sortBy, sortOrder],
    queryFn: () =>
      projectsApi
        .list({
          limit: 50,
          search: search || undefined,
          sort_by: sortBy,
          sort_order: sortOrder,
        })
        .then((response) => response.items),
  });

  const sortedCount = useMemo(() => projects.data?.length ?? 0, [projects.data]);

  const createProject = useMutation({
    mutationFn: (payload: ProjectCreateInput) => projectsApi.create(payload),
    onSuccess: () => {
      setCreateName('');
      setCreateKey('');
      setCreateDescription('');
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Project created.' });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: ['issues'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to create project.');
      setActionError(parsed.message);
      setFieldErrors(parsed.fields);
    },
  });

  const updateProject = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: ProjectUpdateInput }) => projectsApi.update(id, payload),
    onSuccess: () => {
      setEditingId(null);
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Project updated.' });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: ['issues'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to update project.');
      setActionError(parsed.message);
      setFieldErrors(parsed.fields);
    },
  });

  const deleteProject = useMutation({
    mutationFn: (id: string) => projectsApi.delete(id),
    onSuccess: () => {
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Project deleted.' });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: ['issues'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to delete project.');
      setActionError(parsed.message);
      setFieldErrors(parsed.fields);
    },
  });

  function startEdit(id: string, name: string, key: string, description: string | null) {
    setEditingId(id);
    setEditName(name);
    setEditKey(key);
    setEditDescription(description ?? '');
    setActionError(null);
    setFieldErrors(null);
  }

  function handleCreate(event: FormEvent) {
    event.preventDefault();
    setActionError(null);
    setFieldErrors(null);
    createProject.mutate({
      name: createName,
      key: createKey.toUpperCase(),
      description: createDescription || null,
    });
  }

  function handleUpdate(event: FormEvent, id: string) {
    event.preventDefault();
    setActionError(null);
    setFieldErrors(null);
    updateProject.mutate({
      id,
      payload: {
        name: editName,
        key: editKey.toUpperCase(),
        description: editDescription || null,
      },
    });
  }

  return (
    <div>
      <PageHeader title="Projects" subtitle="Track project ownership, keys, and progress at a glance." />
      <div className="panel" style={{ padding: 18, marginBottom: 18 }}>
        <form onSubmit={handleCreate} style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 2fr auto', gap: 12, alignItems: 'end' }}>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Name
            </div>
            <Input value={createName} onChange={(event) => setCreateName(event.target.value)} placeholder="Platform" required />
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Key
            </div>
            <Input
              value={createKey}
              onChange={(event) => setCreateKey(event.target.value.toUpperCase())}
              placeholder="PLAT"
              required
            />
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Description
            </div>
            <Input value={createDescription} onChange={(event) => setCreateDescription(event.target.value)} placeholder="Optional description" />
          </div>
          <Button type="submit" disabled={createProject.isPending}>
            {createProject.isPending ? 'Creating' : 'Create'}
          </Button>
        </form>
        <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 1fr', gap: 12, marginTop: 14 }}>
          <Input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search projects" />
          <select value={sortBy} onChange={(event) => setSortBy(event.target.value as 'name' | 'created_at' | 'updated_at')} className="panel-soft" style={{ padding: 10, border: '2px solid var(--border-strong)', borderRadius: 10 }}>
            <option value="name">Sort: Name</option>
            <option value="created_at">Sort: Created</option>
            <option value="updated_at">Sort: Updated</option>
          </select>
          <select value={sortOrder} onChange={(event) => setSortOrder(event.target.value as 'asc' | 'desc')} className="panel-soft" style={{ padding: 10, border: '2px solid var(--border-strong)', borderRadius: 10 }}>
            <option value="asc">Order: Asc</option>
            <option value="desc">Order: Desc</option>
          </select>
        </div>
      </div>
      {actionError ? <ErrorBanner message={actionError} /> : null}
      {projects.isLoading ? <Spinner label="Loading projects" /> : null}
      {projects.isFetching && !projects.isLoading ? <div style={{ color: 'var(--text-secondary)', marginBottom: 12 }}>Refreshing projects...</div> : null}
      {projects.isError ? <ErrorBanner message={(projects.error as Error).message} /> : null}
      {projects.data?.length ? (
        <div style={{ display: 'grid', gap: 12 }}>
          <div style={{ color: 'var(--text-secondary)' }}>{sortedCount} projects</div>
        </div>
      ) : null}
      {projects.data?.length ? (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: 16 }}>
          {projects.data.map((project) => (
            <article key={project.id} className="panel" style={{ padding: 18 }}>
              {editingId === project.id ? (
                <form onSubmit={(event) => handleUpdate(event, project.id)} style={{ display: 'grid', gap: 10 }}>
                  <Input value={editName} onChange={(event) => setEditName(event.target.value)} required />
                  <Input value={editKey} onChange={(event) => setEditKey(event.target.value.toUpperCase())} required />
                  <Input value={editDescription} onChange={(event) => setEditDescription(event.target.value)} placeholder="Optional description" />
                  {fieldErrors?.name ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.name}</div> : null}
                  {fieldErrors?.key ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.key}</div> : null}
                  {fieldErrors?.description ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.description}</div> : null}
                  <div style={{ display: 'flex', gap: 10 }}>
                    <Button type="submit" disabled={updateProject.isPending}>
                      {updateProject.isPending ? 'Saving' : 'Save'}
                    </Button>
                    <Button type="button" variant="ghost" onClick={() => setEditingId(null)}>
                      Cancel
                    </Button>
                  </div>
                </form>
              ) : (
                <>
                  <div className="label" style={{ color: 'var(--text-secondary)', marginBottom: 8 }}>
                    {project.key}
                  </div>
                  <h3 style={{ margin: 0 }}>{project.name}</h3>
                  <p style={{ color: 'var(--text-secondary)' }}>{project.description ?? 'No description'}</p>
                  <div style={{ color: 'var(--text-secondary)', fontSize: 13 }}>
                    Updated {formatDate(project.updated_at)} · Total {project.issue_counts.total} issues
                  </div>
                  <div style={{ color: 'var(--text-secondary)', fontSize: 13, marginTop: 6 }}>
                    Status mix: {project.issue_counts.backlog} backlog, {project.issue_counts.in_progress} in progress, {project.issue_counts.done} done
                  </div>
                  <div style={{ color: 'var(--text-secondary)', fontSize: 13, marginTop: 6 }}>
                    Active sprint: {project.active_sprint?.name ?? 'None'}
                  </div>
                  <div style={{ display: 'flex', gap: 10, marginTop: 12 }}>
                    <Button type="button" variant="ghost" onClick={() => startEdit(project.id, project.name, project.key, project.description)}>
                      Edit
                    </Button>
                    <Button
                      type="button"
                      variant="danger"
                      disabled={deleteProject.isPending}
                      onClick={() => {
                        if (!window.confirm(`Delete project "${project.name}"?`)) return;
                        setActionError(null);
                        setFieldErrors(null);
                        deleteProject.mutate(project.id);
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
      {projects.data && projects.data.length === 0 ? (
        <EmptyState title="No projects" description="Projects will appear here once the backend contains project data." />
      ) : null}
    </div>
  );
}

function parseApiError(error: unknown, fallback: string) {
  if (error instanceof ApiError) {
    return {
      message: error.message || fallback,
      fields: error.fields ?? null,
    };
  }
  return { message: fallback, fields: null };
}
