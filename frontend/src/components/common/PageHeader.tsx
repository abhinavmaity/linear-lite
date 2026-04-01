import { ReactNode } from 'react';

export function PageHeader({
  title,
  subtitle,
  actions,
}: {
  title: string;
  subtitle?: string;
  actions?: ReactNode;
}) {
  return (
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'end', gap: 16, marginBottom: 24 }}>
      <div>
        <h1
          className="label"
          style={{
            margin: 0,
            fontSize: 'clamp(2rem, 3vw, 3.6rem)',
            lineHeight: 1,
          }}
        >
          {title}
        </h1>
        {subtitle ? <p style={{ color: 'var(--text-secondary)', marginTop: 8 }}>{subtitle}</p> : null}
      </div>
      {actions}
    </div>
  );
}
