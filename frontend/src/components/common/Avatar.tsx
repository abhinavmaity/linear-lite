import { UserSummary } from 'types/domain';

export function Avatar({ user, size = 36 }: { user: Pick<UserSummary, 'name' | 'avatar_url'> | null; size?: number }) {
  const initials =
    user?.name
      ?.split(' ')
      .slice(0, 2)
      .map((part) => part[0])
      .join('')
      .toUpperCase() ?? '?';

  if (user?.avatar_url) {
    return (
      <img
        src={user.avatar_url}
        alt={user.name}
        style={{
          width: size,
          height: size,
          borderRadius: 999,
          border: '2px solid var(--border-strong)',
          objectFit: 'cover',
        }}
      />
    );
  }

  return (
    <div
      className="label"
      style={{
        width: size,
        height: size,
        borderRadius: 999,
        border: '2px solid var(--border-strong)',
        display: 'grid',
        placeItems: 'center',
        background: 'var(--bg-accent-soft)',
        fontWeight: 700,
      }}
    >
      {initials}
    </div>
  );
}
