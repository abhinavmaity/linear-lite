import { ReactNode } from 'react';

export function Badge({ children, tone = 'default' }: { children: ReactNode; tone?: 'default' | 'accent' | 'info' }) {
  const styles =
    tone === 'accent'
      ? { background: 'var(--bg-accent)', color: 'var(--text-on-accent)' }
      : tone === 'info'
        ? { background: 'var(--bg-secondary)', color: '#fff' }
        : { background: 'var(--bg-muted)', color: 'var(--text-primary)' };

  return (
    <span
      className="label"
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: 8,
        padding: '6px 10px',
        border: '2px solid var(--border-strong)',
        borderRadius: 999,
        fontSize: 12,
        fontWeight: 700,
        ...styles,
      }}
    >
      {children}
    </span>
  );
}
