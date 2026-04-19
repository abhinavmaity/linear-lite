import { ReactNode } from 'react';
import { NavLink } from 'react-router-dom';
import { Button } from './Button';
import { Avatar } from './Avatar';
import { useThemeStore } from 'store/themeStore';
import { useAuthStore } from 'store/authStore';
import { useUIStore } from 'store/uiStore';

const links = [
  ['Dashboard', '/dashboard'],
  ['Issues', '/issues'],
  ['Board', '/board'],
  ['Projects', '/projects'],
  ['Sprints', '/sprints'],
  ['Labels', '/labels'],
  ['Team', '/team'],
];

export function AppShell({ children }: { children: ReactNode }) {
  const toggleTheme = useThemeStore((state) => state.toggleTheme);
  const user = useAuthStore((state) => state.user);
  const clearSession = useAuthStore((state) => state.clearSession);
  const openCreateIssue = useUIStore((state) => state.openCreateIssue);

  return (
    <div className="app-shell grid-bg" style={{ display: 'grid', gridTemplateColumns: '250px 1fr' }}>
      <aside
        style={{
          padding: 20,
          borderRight: '3px solid var(--border-strong)',
          background: 'color-mix(in srgb, var(--bg-canvas) 88%, var(--bg-muted))',
          position: 'sticky',
          top: 0,
          height: '100vh',
        }}
      >
        <div className="panel" style={{ padding: 16, marginBottom: 20 }}>
          <div className="label" style={{ fontSize: 28 }}>
            Linear Lite
          </div>
        </div>
        <nav style={{ display: 'grid', gap: 10 }}>
          {links.map(([label, to]) => (
            <NavLink
              key={to}
              to={to}
              style={({ isActive }) => ({
                padding: '14px 16px',
                border: '2px solid var(--border-strong)',
                borderRadius: 10,
                boxShadow: 'var(--shadow-soft)',
                background: isActive ? 'var(--bg-accent)' : 'var(--bg-elevated)',
                color: isActive ? 'var(--text-on-accent)' : 'var(--text-primary)',
                fontFamily: 'Space Grotesk, sans-serif',
                textTransform: 'uppercase',
                letterSpacing: '0.08em',
                fontWeight: 700,
              })}
            >
              {label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <div>
        <header
          style={{
            position: 'sticky',
            top: 0,
            zIndex: 20,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            gap: 16,
            padding: '18px 24px',
            background: 'color-mix(in srgb, var(--bg-canvas) 88%, transparent)',
            backdropFilter: 'blur(10px)',
            borderBottom: '2px solid var(--border-soft)',
          }}
        >
          <div className="label" style={{ fontSize: 13 }}>
              Issue tracking for focused team workflows
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <Button variant="secondary" onClick={toggleTheme}>
              Toggle Theme
            </Button>
            <Button onClick={openCreateIssue}>Create Issue</Button>
            <div className="panel-soft" style={{ padding: '6px 10px', display: 'flex', alignItems: 'center', gap: 10 }}>
              <Avatar user={user} />
              <div>
                <div style={{ fontWeight: 700 }}>{user?.name ?? 'Guest'}</div>
                <button onClick={clearSession} style={{ background: 'transparent', border: 'none', padding: 0, color: 'var(--text-secondary)' }}>
                  Logout
                </button>
              </div>
            </div>
          </div>
        </header>
        <main className="page">{children}</main>
      </div>
    </div>
  );
}
