import { ReactNode } from 'react';

export function EmptyState({ title, description, action }: { title: string; description: string; action?: ReactNode }) {
  return (
    <div className="panel" style={{ padding: 24, textAlign: 'center' }}>
      <h3 className="label" style={{ fontSize: 24, margin: '0 0 8px' }}>
        {title}
      </h3>
      <p style={{ color: 'var(--text-secondary)', margin: '0 0 16px' }}>{description}</p>
      {action}
    </div>
  );
}
