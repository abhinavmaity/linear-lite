import { useQuery } from '@tanstack/react-query';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { labelsApi } from 'services/labelsApi';

export function LabelsPage() {
  const labels = useQuery({
    queryKey: ['labels', 'page'],
    queryFn: () => labelsApi.list({ limit: 100 }).then((response) => response.items),
  });

  return (
    <div>
      <PageHeader title="Labels" subtitle="Scaffolded route with label management data source." />
      {labels.isLoading ? <Spinner label="Loading labels" /> : null}
      {labels.isError ? <ErrorBanner message={(labels.error as Error).message} /> : null}
      {labels.data?.length ? (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: 16 }}>
          {labels.data.map((label) => (
            <article key={label.id} className="panel" style={{ padding: 18 }}>
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
            </article>
          ))}
        </div>
      ) : null}
      {labels.data && labels.data.length === 0 ? (
        <EmptyState title="No labels" description="Labels will appear here when available." />
      ) : null}
    </div>
  );
}
