import { Link } from 'react-router-dom';
import { Skeleton } from 'boneyard-js/react';
import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { StatCard } from 'components/dashboard/StatCard';
import { useDashboardStats } from 'features/dashboard/dashboardQueries';
import { getBannerErrorMessage } from 'utils/errorPresentation';
import { relativeTime, titleCase } from 'utils/format';

export function DashboardPage() {
  const stats = useDashboardStats();

  return (
    <div>
      <PageHeader title="Dashboard" subtitle="Only architecture-supported metrics are rendered here." />
      {stats.isError ? <ErrorBanner message={getBannerErrorMessage(stats.error, 'Unable to load dashboard right now.')} /> : null}
      <Skeleton
        name="dashboard-page"
        loading={stats.isLoading}
        fallback={<Spinner label="Loading dashboard" />}
        fixture={
          <div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, minmax(0, 1fr))', gap: 16, marginBottom: 24 }}>
              <StatCard label="Total Issues" value={42} />
              <StatCard label="My Issues" value={12} />
              <StatCard label="In Progress" value={5} accent />
              <StatCard label="Done This Week" value={7} />
            </div>
            <div className="two-col">
              <section className="panel" style={{ padding: 20 }}>
                <div className="label" style={{ fontSize: 24, marginBottom: 14 }}>
                  Recent Activity
                </div>
                <div style={{ display: 'grid', gap: 14 }}>
                  <div className="panel-soft" style={{ padding: 14 }}>
                    <div style={{ fontWeight: 700 }}>Alex</div>
                    <div style={{ color: 'var(--text-secondary)' }}>Updated · Status · 5m ago</div>
                  </div>
                </div>
              </section>
              <aside className="panel" style={{ padding: 20 }}>
                <div className="label" style={{ fontSize: 24, marginBottom: 14 }}>
                  Active Sprint
                </div>
                <div className="panel-soft" style={{ padding: 16 }}>
                  <div style={{ fontWeight: 700, fontSize: 22 }}>Sprint 1</div>
                  <div style={{ color: 'var(--text-secondary)', marginTop: 8 }}>2026-04-01 to 2026-04-14</div>
                </div>
              </aside>
            </div>
          </div>
        }
      >
        {stats.data ? (
          <>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, minmax(0, 1fr))', gap: 16, marginBottom: 24 }}>
            <StatCard label="Total Issues" value={stats.data.total_issues} />
            <StatCard label="My Issues" value={stats.data.my_issues} />
            <StatCard label="In Progress" value={stats.data.in_progress} accent />
            <StatCard label="Done This Week" value={stats.data.done_this_week} />
          </div>
          <div className="two-col">
            <section className="panel" style={{ padding: 20 }}>
              <div className="label" style={{ fontSize: 24, marginBottom: 14 }}>
                Recent Activity
              </div>
              {stats.data.recent_activity.length === 0 ? (
                <EmptyState title="No activity" description="No recent issue activity is available yet." />
              ) : (
                <div style={{ display: 'grid', gap: 14 }}>
                  {stats.data.recent_activity.map((activity) => (
                    <div key={activity.id} className="panel-soft" style={{ padding: 14 }}>
                      <div style={{ fontWeight: 700 }}>{activity.user?.name ?? 'Unknown user'}</div>
                      <div style={{ color: 'var(--text-secondary)' }}>
                        {titleCase(activity.action)}
                        {activity.field_name ? ` · ${titleCase(activity.field_name)}` : ''} · {relativeTime(activity.created_at)}
                      </div>
                      <div style={{ marginTop: 8, fontSize: 14 }}>
                        {activity.old_value || activity.new_value
                          ? `${activity.old_value ?? 'empty'} → ${activity.new_value ?? 'empty'}`
                          : 'No field change payload'}
                      </div>
                      <Link to={`/issues/${activity.issue_id}`} style={{ display: 'inline-block', marginTop: 8, fontWeight: 700 }}>
                        Open issue
                      </Link>
                    </div>
                  ))}
                </div>
              )}
            </section>
            <aside className="panel" style={{ padding: 20 }}>
              <div className="label" style={{ fontSize: 24, marginBottom: 14 }}>
                Active Sprint
              </div>
              {stats.data.active_sprint ? (
                <div className="panel-soft" style={{ padding: 16 }}>
                  <div style={{ fontWeight: 700, fontSize: 22 }}>{stats.data.active_sprint.name}</div>
                  <div style={{ color: 'var(--text-secondary)', marginTop: 8 }}>
                    {stats.data.active_sprint.start_date} to {stats.data.active_sprint.end_date}
                  </div>
                </div>
              ) : (
                <EmptyState title="No active sprint" description="Sprint data will appear here when a project has one active sprint." />
              )}
            </aside>
          </div>
          </>
        ) : null}
      </Skeleton>
    </div>
  );
}
