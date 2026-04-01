import { Link } from 'react-router-dom';
import { IssueSummary } from 'types/domain';
import { Badge } from 'components/common/Badge';
import { Avatar } from 'components/common/Avatar';

export function IssueCard({ issue, compact = false }: { issue: IssueSummary; compact?: boolean }) {
  return (
    <Link to={`/issues/${issue.id}`} className="panel-soft" style={{ padding: 14, display: 'grid', gap: 10 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', gap: 8 }}>
        <div className="label" style={{ fontSize: 11, color: 'var(--text-secondary)' }}>
          {issue.identifier}
        </div>
        <Badge tone="info">{issue.priority}</Badge>
      </div>
      <div style={{ fontWeight: 700 }}>{issue.title}</div>
      {!compact ? (
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
          <Badge>{issue.status}</Badge>
          {issue.labels.slice(0, 2).map((label) => (
            <span
              key={label.id}
              style={{
                padding: '4px 8px',
                borderRadius: 999,
                border: '1px solid var(--border-strong)',
                background: label.color,
                color: '#fff',
                fontSize: 12,
              }}
            >
              {label.name}
            </span>
          ))}
        </div>
      ) : null}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <span style={{ color: 'var(--text-secondary)', fontSize: 13 }}>{issue.project.key}</span>
        <Avatar user={issue.assignee} size={28} />
      </div>
    </Link>
  );
}
