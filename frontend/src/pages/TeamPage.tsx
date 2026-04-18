import { useQuery } from '@tanstack/react-query';
import { useMemo, useState } from 'react';
import { Skeleton } from 'boneyard-js/react';
import { Avatar } from 'components/common/Avatar';
import { Badge } from 'components/common/Badge';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { PageHeader } from 'components/common/PageHeader';
import { Select } from 'components/common/Select';
import { Spinner } from 'components/common/Spinner';
import { usersApi } from 'services/usersApi';
import { getBannerErrorMessage } from 'utils/errorPresentation';
import { formatDate } from 'utils/format';

export function TeamPage() {
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'created_at'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  const users = useQuery({
    queryKey: ['users', 'team-page', search, sortBy, sortOrder],
    queryFn: () =>
      usersApi
        .list({
          limit: 50,
          search: search || undefined,
          sort_by: sortBy,
          sort_order: sortOrder,
        })
        .then((response) => response.items),
  });

  const totalUsers = useMemo(() => users.data?.length ?? 0, [users.data]);

  return (
    <div>
      <PageHeader title="Team" subtitle="Read-only MVP route backed by the user list endpoint." />
      <div className="panel" style={{ padding: 18, marginBottom: 18 }}>
        <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 1fr', gap: 12 }}>
          <Input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search team members" />
          <Select value={sortBy} onChange={(event) => setSortBy(event.target.value as 'name' | 'created_at')}>
            <option value="name">Sort: Name</option>
            <option value="created_at">Sort: Joined</option>
          </Select>
          <Select value={sortOrder} onChange={(event) => setSortOrder(event.target.value as 'asc' | 'desc')}>
            <option value="asc">Order: Asc</option>
            <option value="desc">Order: Desc</option>
          </Select>
        </div>
      </div>
      <Skeleton
        name="team-page"
        loading={users.isLoading}
        fallback={<Spinner label="Loading team" />}
        fixture={
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: 16 }}>
            <article className="panel" style={{ padding: 18, display: 'grid', gap: 10 }}>
              <div style={{ display: 'flex', gap: 14, alignItems: 'center' }}>
                <div style={{ width: 48, height: 48, borderRadius: '50%', background: 'var(--bg-muted)' }} />
                <div>
                  <div style={{ fontWeight: 700 }}>Alex Doe</div>
                  <div style={{ color: 'var(--text-secondary)' }}>alex@example.com</div>
                </div>
              </div>
            </article>
          </div>
        }
      >
        {users.isFetching && !users.isLoading ? <div style={{ color: 'var(--text-secondary)', marginBottom: 12 }}>Refreshing team...</div> : null}
        {users.isError ? <ErrorBanner message={getBannerErrorMessage(users.error, 'Unable to load team members right now.')} /> : null}
        {users.data ? (
          <div style={{ marginBottom: 12, color: 'var(--text-secondary)', display: 'flex', gap: 10, alignItems: 'center' }}>
            <span>{totalUsers} members</span>
            <Badge>Read only</Badge>
          </div>
        ) : null}
        {users.data?.length ? (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: 16 }}>
            {users.data.map((user) => (
              <article key={user.id} className="panel" style={{ padding: 18, display: 'grid', gap: 10 }}>
                <div style={{ display: 'flex', gap: 14, alignItems: 'center' }}>
                  <Avatar user={user} size={48} />
                  <div>
                    <div style={{ fontWeight: 700 }}>{user.name}</div>
                    <div style={{ color: 'var(--text-secondary)' }}>{user.email}</div>
                  </div>
                </div>
                <div>
                  <div style={{ color: 'var(--text-secondary)', fontSize: 13 }}>Joined {formatDate(user.created_at)}</div>
                  <div style={{ color: 'var(--text-secondary)', fontSize: 13 }}>Last update {formatDate(user.updated_at)}</div>
                </div>
              </article>
            ))}
          </div>
        ) : null}
        {users.data && users.data.length === 0 ? (
          <EmptyState title="No team members" description="No users match the current search and sort options." />
        ) : null}
      </Skeleton>
    </div>
  );
}
