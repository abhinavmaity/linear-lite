import { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
}

export function Button({ children, style, variant = 'primary', ...props }: ButtonProps) {
  const colors =
    variant === 'primary'
      ? { background: 'var(--bg-accent)', color: 'var(--text-on-accent)' }
      : variant === 'secondary'
        ? { background: 'var(--bg-accent-soft)', color: 'var(--text-primary)' }
        : variant === 'danger'
          ? { background: 'var(--danger)', color: '#fff' }
          : { background: 'transparent', color: 'var(--text-primary)' };

  return (
    <button
      {...props}
      style={{
        ...colors,
        border: '2px solid var(--border-strong)',
        borderRadius: '10px',
        padding: '12px 16px',
        fontWeight: 700,
        textTransform: 'uppercase',
        letterSpacing: '0.06em',
        boxShadow: 'var(--shadow-soft)',
        ...style,
      }}
    >
      {children}
    </button>
  );
}
