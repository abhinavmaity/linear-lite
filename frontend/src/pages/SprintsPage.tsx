import { useQuery } from '@tanstack/react-query';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { sprintsApi } from 'services/sprintsApi';

export function SprintsPage() {
  const sprints = useQuery({
    queryKey: ['sprints', 'page'],
    queryFn: () => sprintsApi.list({ limit: 50 }).then((response) => response.items),
  });

  return (
    <div>
      <PageHeader title="Sprints" subtitle="Scaffolded MVP route with live sprint list data." />
      {sprints.isLoading ? <Spinner label="Loading sprints" /> : null}
      {sprints.isError ? <ErrorBanner message={(sprints.error as Error).message} /> : null}
      {sprints.data?.length ? (
        <div style={{ display: 'grid', gap: 16 }}>
          {sprints.data.map((sprint) => (
            <article key={sprint.id} className="panel" style={{ padding: 18 }}>
              <h3 style={{ marginTop: 0 }}>{sprint.name}</h3>
              <p style={{ color: 'var(--text-secondary)' }}>
                {sprint.status} · {sprint.start_date} to {sprint.end_date}
              </p>
            </article>
          ))}
        </div>
      ) : null}
      {sprints.data && sprints.data.length === 0 ? (
        <EmptyState title="No sprints" description="Sprint data will appear here when available." />
      ) : null}
    </div>
  );
}
