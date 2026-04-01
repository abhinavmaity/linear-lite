import { InputHTMLAttributes, forwardRef } from 'react';

export const Input = forwardRef<HTMLInputElement, InputHTMLAttributes<HTMLInputElement>>(function Input(props, ref) {
  return (
    <input
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
