import { useUIStore } from 'store/uiStore';

export function Toaster() {
  const toasts = useUIStore((state) => state.toasts);

  return (
    <div style={{ position: 'fixed', right: 16, bottom: 16, display: 'grid', gap: 12, zIndex: 1100 }}>
      {toasts.map((toast) => (
        <div
          key={toast.id}
          className="panel-soft"
          style={{
            padding: '12px 14px',
            minWidth: 240,
            background:
              toast.tone === 'error'
                ? 'color-mix(in srgb, var(--danger) 18%, var(--bg-elevated))'
                : toast.tone === 'success'
                  ? 'color-mix(in srgb, var(--success) 18%, var(--bg-elevated))'
                  : 'var(--bg-elevated)',
          }}
        >
          {toast.message}
        </div>
      ))}
    </div>
  );
}
