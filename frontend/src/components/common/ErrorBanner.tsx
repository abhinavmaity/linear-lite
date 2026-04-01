export function ErrorBanner({ message }: { message: string }) {
  return (
    <div
      className="panel-soft"
      style={{
        padding: 16,
        background: 'color-mix(in srgb, var(--danger) 10%, var(--bg-elevated))',
        color: 'var(--danger)',
      }}
    >
      {message}
    </div>
  );
}
