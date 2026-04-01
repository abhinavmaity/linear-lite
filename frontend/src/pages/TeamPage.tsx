import { useQuery } from '@tanstack/react-query';
import { Avatar } from 'components/common/Avatar';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { usersApi } from 'services/usersApi';

export function TeamPage() {
  const users = useQuery({
    queryKey: ['users', 'team-page'],
    queryFn: () => usersApi.list({ limit: 50 }).then((response) => response.items),
  });

  return (
    <div>
      <PageHeader title="Team" subtitle="Read-only MVP route backed by the user list endpoint." />
      {users.isLoading ? <Spinner label="Loading team" /> : null}
      {users.isError ? <ErrorBanner message={(users.error as Error).message} /> : null}
      {users.data?.length ? (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: 16 }}>
          {users.data.map((user) => (
            <article key={user.id} className="panel" style={{ padding: 18, display: 'flex', gap: 14, alignItems: 'center' }}>
              <Avatar user={user} size={48} />
              <div>
                <div style={{ fontWeight: 700 }}>{user.name}</div>
                <div style={{ color: 'var(--text-secondary)' }}>{user.email}</div>
              </div>
            </article>
          ))}
        </div>
      ) : null}
      {users.data && users.data.length === 0 ? (
        <EmptyState title="No team members" description="Users will appear here when the backend has seeded user data." />
      ) : null}
    </div>
  );
}
