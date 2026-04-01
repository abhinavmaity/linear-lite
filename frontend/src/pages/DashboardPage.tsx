import { EmptyState } from 'components/common/EmptyState';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { PageHeader } from 'components/common/PageHeader';
import { Spinner } from 'components/common/Spinner';
import { StatCard } from 'components/dashboard/StatCard';
import { useDashboardStats } from 'features/dashboard/dashboardQueries';
import { relativeTime, titleCase } from 'utils/format';

export function DashboardPage() {
  const stats = useDashboardStats();

  return (
    <div>
      <PageHeader title="Dashboard" subtitle="Only architecture-supported metrics are rendered here." />
      {stats.isLoading ? <Spinner label="Loading dashboard" /> : null}
      {stats.isError ? <ErrorBanner message={(stats.error as Error).message} /> : null}
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
                      <div style={{ fontWeight: 700 }}>{activity.user.name}</div>
                      <div style={{ color: 'var(--text-secondary)' }}>
                        {titleCase(activity.action)} on issue {activity.issue_id.slice(0, 8)} · {relativeTime(activity.created_at)}
                      </div>
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
    </div>
  );
}
