import { useMutation } from '@tanstack/react-query';
import { FormEvent, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Button } from 'components/common/Button';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { authApi } from 'services/authApi';
import { useAuthStore } from 'store/authStore';
import { ApiError } from 'types/api';

export function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const setSession = useAuthStore((state) => state.setSession);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const mutation = useMutation({
    mutationFn: () => authApi.login({ email, password }),
    onSuccess: (response) => {
      setSession(response.data.token, response.data.user);
      navigate(location.state?.from?.pathname ?? '/dashboard', { replace: true });
    },
  });

  const error = mutation.error instanceof ApiError ? mutation.error.message : null;

  function onSubmit(event: FormEvent) {
    event.preventDefault();
    mutation.mutate();
  }

  return (
    <main
      className="grid-bg"
      style={{ minHeight: '100vh', display: 'grid', placeItems: 'center', padding: 24, backgroundColor: 'var(--bg-canvas)' }}
    >
      <div className="panel" style={{ width: 'min(460px, 100%)', padding: 28 }}>
        <div className="label" style={{ fontSize: 40, marginBottom: 8 }}>
          Login
        </div>
        <p style={{ color: 'var(--text-secondary)', marginBottom: 24 }}>
          Sign in with the architecture-defined email and password flow.
        </p>
        <form onSubmit={onSubmit} style={{ display: 'grid', gap: 16 }}>
          {error ? <ErrorBanner message={error} /> : null}
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Email
            </div>
            <Input type="email" value={email} onChange={(event) => setEmail(event.target.value)} placeholder="alex@example.com" />
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Password
            </div>
            <Input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
          </div>
          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? 'Signing In' : 'Enter Dashboard'}
          </Button>
        </form>
        <div style={{ marginTop: 18, color: 'var(--text-secondary)' }}>
          Need an account? <Link to="/register">Register</Link>
        </div>
      </div>
    </main>
  );
}
