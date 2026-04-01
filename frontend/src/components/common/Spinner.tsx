export function Spinner({ label, fullScreen = false }: { label?: string; fullScreen?: boolean }) {
  return (
    <div
      style={{
        minHeight: fullScreen ? '100vh' : 120,
        display: 'grid',
        placeItems: 'center',
        gap: 12,
      }}
    >
      <div
        style={{
          width: 32,
          height: 32,
          borderRadius: '50%',
          border: '4px solid var(--border-soft)',
          borderTopColor: 'var(--bg-accent)',
          animation: 'spin 0.9s linear infinite',
        }}
      />
      {label ? <div className="label">{label}</div> : null}
      <style>{'@keyframes spin { to { transform: rotate(360deg); } }'}</style>
    </div>
  );
}
