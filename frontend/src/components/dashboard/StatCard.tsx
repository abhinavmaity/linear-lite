export function StatCard({ label, value, accent }: { label: string; value: string | number; accent?: boolean }) {
  return (
    <div
      className="panel"
      style={{
        padding: 18,
        background: accent ? 'color-mix(in srgb, var(--bg-accent-soft) 28%, var(--bg-elevated))' : undefined,
      }}
    >
      <div className="label" style={{ fontSize: 12, color: 'var(--text-secondary)' }}>
        {label}
      </div>
      <div style={{ fontSize: 36, fontWeight: 700, marginTop: 10 }}>{value}</div>
    </div>
  );
}
