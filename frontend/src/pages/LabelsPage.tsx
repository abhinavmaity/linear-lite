import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { FormEvent, useState } from 'react';
import { Skeleton } from 'boneyard-js/react';
import { Button } from 'components/common/Button';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { PageHeader } from 'components/common/PageHeader';
import { Select } from 'components/common/Select';
import { Spinner } from 'components/common/Spinner';
import { DEFAULT_LABEL_COLOR, LABEL_COLOR_OPTIONS } from 'constants/labelColors';
import { labelsApi, LabelCreateInput, LabelUpdateInput } from 'services/labelsApi';
import { useUIStore } from 'store/uiStore';
import { getBannerErrorMessage, parseUiError } from 'utils/errorPresentation';
import { formatDate } from 'utils/format';

export function LabelsPage() {
  const queryClient = useQueryClient();
  const pushToast = useUIStore((state) => state.pushToast);
  const [search, setSearch] = useState('');
  const [name, setName] = useState('');
  const [color, setColor] = useState<string>(DEFAULT_LABEL_COLOR);
  const [description, setDescription] = useState('');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState('');
  const [editColor, setEditColor] = useState<string>(DEFAULT_LABEL_COLOR);
  const [editDescription, setEditDescription] = useState('');
  const [actionError, setActionError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string> | null>(null);

  const labels = useQuery({
    queryKey: ['labels', 'page', search],
    queryFn: () => labelsApi.list({ limit: 100, search: search || undefined }).then((response) => response.items),
  });

  const createLabel = useMutation({
    mutationFn: (payload: LabelCreateInput) => labelsApi.create(payload),
    onSuccess: () => {
      setName('');
      setColor(DEFAULT_LABEL_COLOR);
      setDescription('');
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Label created.' });
      queryClient.invalidateQueries({ queryKey: ['labels'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to create label.');
      setActionError(parsed.summary);
      setFieldErrors(parsed.fields);
    },
  });

  const updateLabel = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: LabelUpdateInput }) => labelsApi.update(id, payload),
    onSuccess: () => {
      setEditingId(null);
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Label updated.' });
      queryClient.invalidateQueries({ queryKey: ['labels'] });
      queryClient.invalidateQueries({ queryKey: ['issues'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to update label.');
      setActionError(parsed.summary);
      setFieldErrors(parsed.fields);
    },
  });

  const deleteLabel = useMutation({
    mutationFn: (id: string) => labelsApi.delete(id),
    onSuccess: () => {
      setActionError(null);
      setFieldErrors(null);
      pushToast({ tone: 'success', message: 'Label deleted.' });
      queryClient.invalidateQueries({ queryKey: ['labels'] });
      queryClient.invalidateQueries({ queryKey: ['issues'] });
    },
    onError: (error) => {
      const parsed = parseApiError(error, 'Failed to delete label.');
      setActionError(parsed.summary);
      setFieldErrors(parsed.fields);
    },
  });

  function handleCreate(event: FormEvent) {
    event.preventDefault();
    setActionError(null);
    setFieldErrors(null);
    createLabel.mutate({
      name,
      color,
      description: description || null,
    });
  }

  function startEdit(labelId: string, labelName: string, labelColor: string, labelDescription: string | null) {
    setEditingId(labelId);
    setEditName(labelName);
    setEditColor(labelColor);
    setEditDescription(labelDescription ?? '');
    setActionError(null);
    setFieldErrors(null);
  }

  function handleUpdate(event: FormEvent, labelId: string) {
    event.preventDefault();
    setActionError(null);
    setFieldErrors(null);
    updateLabel.mutate({
      id: labelId,
      payload: {
        name: editName,
        color: editColor,
        description: editDescription || null,
      },
    });
  }

  return (
    <div>
      <PageHeader title="Labels" subtitle="Organize issue categories with consistent label definitions." />
      <div className="panel" style={{ padding: 18, marginBottom: 18 }}>
        <form onSubmit={handleCreate} style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 2fr auto', gap: 12, alignItems: 'end' }}>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Name
            </div>
            <Input value={name} onChange={(event) => setName(event.target.value)} placeholder="bug" required />
            {fieldErrors?.name ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.name}</div> : null}
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Color
            </div>
            <Select value={color} onChange={(event) => setColor(event.target.value)} required>
              {LABEL_COLOR_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label} ({option.value})
                </option>
              ))}
            </Select>
            {fieldErrors?.color ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.color}</div> : null}
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Description
            </div>
            <Input value={description} onChange={(event) => setDescription(event.target.value)} placeholder="Optional description" />
            {fieldErrors?.description ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.description}</div> : null}
          </div>
          <Button type="submit" disabled={createLabel.isPending}>
            {createLabel.isPending ? 'Creating' : 'Create'}
          </Button>
        </form>
        <div style={{ marginTop: 12 }}>
          <Input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search labels" />
        </div>
      </div>
      {actionError ? <ErrorBanner message={actionError} /> : null}
      <Skeleton
        name="labels-page"
        loading={labels.isLoading}
        fallback={<Spinner label="Loading labels" />}
        fixture={
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: 16 }}>
            <article className="panel" style={{ padding: 18 }}>
              <div style={{ width: 28, height: 28, borderRadius: 999, background: '#3B82F6', border: '2px solid var(--border-strong)', marginBottom: 12 }} />
              <h3 style={{ margin: '0 0 6px' }}>bug</h3>
              <p style={{ color: 'var(--text-secondary)', margin: 0 }}>Something not working</p>
            </article>
          </div>
        }
      >
        {labels.isFetching && !labels.isLoading ? <div style={{ color: 'var(--text-secondary)', marginBottom: 12 }}>Refreshing labels...</div> : null}
        {labels.isError ? <ErrorBanner message={getBannerErrorMessage(labels.error, 'Unable to load labels right now.')} /> : null}
        {labels.data?.length ? (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: 16 }}>
            {labels.data.map((label) => (
              <article key={label.id} className="panel" style={{ padding: 18 }}>
                {editingId === label.id ? (
                  <form onSubmit={(event) => handleUpdate(event, label.id)} style={{ display: 'grid', gap: 10 }}>
                    <Input value={editName} onChange={(event) => setEditName(event.target.value)} required />
                    <Select value={editColor} onChange={(event) => setEditColor(event.target.value)} required>
                      {LABEL_COLOR_OPTIONS.map((option) => (
                        <option key={option.value} value={option.value}>
                          {option.label} ({option.value})
                        </option>
                      ))}
                    </Select>
                    <Input value={editDescription} onChange={(event) => setEditDescription(event.target.value)} placeholder="Optional description" />
                    {fieldErrors?.name ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.name}</div> : null}
                    {fieldErrors?.color ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.color}</div> : null}
                    {fieldErrors?.description ? <div style={{ color: 'var(--danger)' }}>{fieldErrors.description}</div> : null}
                    <div style={{ display: 'flex', gap: 10 }}>
                      <Button type="submit" disabled={updateLabel.isPending}>
                        {updateLabel.isPending ? 'Saving' : 'Save'}
                      </Button>
                      <Button type="button" variant="ghost" onClick={() => setEditingId(null)}>
                        Cancel
                      </Button>
                    </div>
                  </form>
                ) : (
                  <>
                    <div
                      style={{
                        width: 28,
                        height: 28,
                        borderRadius: 999,
                        background: label.color,
                        border: '2px solid var(--border-strong)',
                        marginBottom: 12,
                      }}
                    />
                    <h3 style={{ margin: '0 0 6px' }}>{label.name}</h3>
                    <p style={{ color: 'var(--text-secondary)', margin: 0 }}>{label.description ?? 'No description'}</p>
                    <p style={{ color: 'var(--text-secondary)', marginTop: 8, fontSize: 13 }}>Updated {formatDate(label.updated_at)}</p>
                    <div style={{ display: 'flex', gap: 10, marginTop: 10 }}>
                      <Button type="button" variant="ghost" onClick={() => startEdit(label.id, label.name, label.color, label.description)}>
                        Edit
                      </Button>
                      <Button
                        type="button"
                        variant="danger"
                        disabled={deleteLabel.isPending}
                        onClick={() => {
                          if (!window.confirm(`Delete label "${label.name}"?`)) return;
                          setActionError(null);
                          setFieldErrors(null);
                          deleteLabel.mutate(label.id);
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
        {labels.data && labels.data.length === 0 ? (
          <EmptyState title="No labels" description="Labels will appear here when available." />
        ) : null}
      </Skeleton>
    </div>
  );
}

function parseApiError(error: unknown, fallback: string) {
  return parseUiError(error, fallback);
}
