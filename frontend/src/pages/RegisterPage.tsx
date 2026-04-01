import { useMutation } from '@tanstack/react-query';
import { FormEvent, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Button } from 'components/common/Button';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { authApi } from 'services/authApi';
import { useAuthStore } from 'store/authStore';
import { ApiError } from 'types/api';

export function RegisterPage() {
  const navigate = useNavigate();
  const setSession = useAuthStore((state) => state.setSession);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [localError, setLocalError] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: () => authApi.register({ name, email, password }),
    onSuccess: (response) => {
      setSession(response.data.token, response.data.user);
      navigate('/dashboard', { replace: true });
    },
  });

  const error = localError || (mutation.error instanceof ApiError ? mutation.error.message : null);

  function onSubmit(event: FormEvent) {
    event.preventDefault();
    if (password !== confirmPassword) {
      setLocalError('Passwords must match.');
      return;
    }
    setLocalError(null);
    mutation.mutate();
  }

  return (
    <main
      className="grid-bg"
      style={{ minHeight: '100vh', display: 'grid', placeItems: 'center', padding: 24, backgroundColor: 'var(--bg-canvas)' }}
    >
      <div className="panel" style={{ width: 'min(560px, 100%)', padding: 28 }}>
        <div className="label" style={{ fontSize: 40, marginBottom: 8 }}>
          Register
        </div>
        <p style={{ color: 'var(--text-secondary)', marginBottom: 24 }}>
          Create an account with the MVP register contract.
        </p>
        <form onSubmit={onSubmit} style={{ display: 'grid', gap: 16 }}>
          {error ? <ErrorBanner message={error} /> : null}
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Name
            </div>
            <Input value={name} onChange={(event) => setName(event.target.value)} />
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Email
            </div>
            <Input type="email" value={email} onChange={(event) => setEmail(event.target.value)} />
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
            <div>
              <div className="label" style={{ marginBottom: 8 }}>
                Password
              </div>
              <Input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
            </div>
            <div>
              <div className="label" style={{ marginBottom: 8 }}>
                Confirm Password
              </div>
              <Input
                type="password"
                value={confirmPassword}
                onChange={(event) => setConfirmPassword(event.target.value)}
              />
            </div>
          </div>
          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? 'Creating' : 'Create Account'}
          </Button>
        </form>
        <div style={{ marginTop: 18, color: 'var(--text-secondary)' }}>
          Already have an account? <Link to="/login">Return to login</Link>
        </div>
      </div>
    </main>
  );
}
