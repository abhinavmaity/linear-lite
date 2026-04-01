import { SelectHTMLAttributes, forwardRef } from 'react';

export const Select = forwardRef<HTMLSelectElement, SelectHTMLAttributes<HTMLSelectElement>>(function Select(
  props,
  ref,
) {
  return (
    <select
      ref={ref}
      {...props}
      style={{
        width: '100%',
        padding: '12px 14px',
        borderRadius: '10px',
        border: '2px solid var(--border-strong)',
        background: 'var(--bg-elevated)',
        color: 'var(--text-primary)',
      }}
    />
  );
});
