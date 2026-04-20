import { useMutation } from '@tanstack/react-query';
import { FormEvent, useEffect, useRef, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Button } from 'components/common/Button';
import { ErrorBanner } from 'components/common/ErrorBanner';
import { Input } from 'components/common/Input';
import { authApi } from 'services/authApi';
import { getGoogleClientId } from 'services/env';
import { useAuthStore } from 'store/authStore';
import { parseUiError } from 'utils/errorPresentation';

type LoginLocationState = {
  from?: {
    pathname?: string;
  };
  registrationSuccess?: boolean;
  email?: string;
};

export function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const locationState = (location.state ?? {}) as LoginLocationState;
  const setSession = useAuthStore((state) => state.setSession);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const googleButtonRef = useRef<HTMLDivElement | null>(null);
  const googleClientID = getGoogleClientId();
  const [googleReady, setGoogleReady] = useState(false);

  const mutation = useMutation({
    mutationFn: () => authApi.login({ email, password }),
    onSuccess: (response) => {
      setSession(response.data.token, response.data.user);
      navigate(locationState.from?.pathname ?? '/dashboard', { replace: true });
    },
  });

  const googleMutation = useMutation({
    mutationFn: (idToken: string) => authApi.loginWithGoogle({ id_token: idToken }),
    onSuccess: (response) => {
      setSession(response.data.token, response.data.user);
      navigate(locationState.from?.pathname ?? '/dashboard', { replace: true });
    },
  });

  const parsedError = mutation.error ? parseUiError(mutation.error, 'Unable to sign in right now. Please try again.') : null;
  const parsedGoogleError = googleMutation.error
    ? parseUiError(googleMutation.error, 'Unable to sign in with Google right now. Please try again.')
    : null;
  const error = parsedGoogleError?.summary ?? parsedError?.summary ?? null;
  const fieldErrors = parsedError?.fields ?? null;
  const registrationMessage = locationState.registrationSuccess
    ? `Account created${locationState.email ? ` for ${locationState.email}` : ''}. Please log in.`
    : null;

  useEffect(() => {
    if (!googleClientID || !googleButtonRef.current) {
      return;
    }

    let cancelled = false;
    const setupGoogleButton = () => {
      if (cancelled || !window.google?.accounts?.id || !googleButtonRef.current) return;
      window.google.accounts.id.initialize({
        client_id: googleClientID,
        callback: (response: { credential?: string }) => {
          const credential = response.credential?.trim();
          if (!credential) {
            return;
          }
          googleMutation.mutate(credential);
        },
      });
      googleButtonRef.current.innerHTML = '';
      window.google.accounts.id.renderButton(googleButtonRef.current, {
        theme: 'outline',
        size: 'large',
        width: 320,
      });
      setGoogleReady(true);
    };

    if (window.google?.accounts?.id) {
      setupGoogleButton();
      return () => {
        cancelled = true;
      };
    }

    const existing = document.querySelector<HTMLScriptElement>('script[data-google-identity="true"]');
    if (existing) {
      existing.addEventListener('load', setupGoogleButton, { once: true });
      return () => {
        cancelled = true;
      };
    }

    const script = document.createElement('script');
    script.src = 'https://accounts.google.com/gsi/client';
    script.async = true;
    script.defer = true;
    script.dataset.googleIdentity = 'true';
    script.addEventListener('load', setupGoogleButton, { once: true });
    document.head.appendChild(script);

    return () => {
      cancelled = true;
    };
  }, [googleClientID, googleMutation]);

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
        {registrationMessage ? (
          <div className="panel-soft" style={{ marginBottom: 16, padding: 12, color: 'var(--text-primary)' }}>
            {registrationMessage}
          </div>
        ) : null}
        <form onSubmit={onSubmit} style={{ display: 'grid', gap: 16 }}>
          {error ? <ErrorBanner message={error} /> : null}
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Email
            </div>
            <Input type="email" value={email} onChange={(event) => setEmail(event.target.value)} placeholder="alex@example.com" />
            {fieldErrors?.email ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.email}</div> : null}
          </div>
          <div>
            <div className="label" style={{ marginBottom: 8 }}>
              Password
            </div>
            <Input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
            {fieldErrors?.password ? <div style={{ color: 'var(--danger)', marginTop: 6 }}>{fieldErrors.password}</div> : null}
          </div>
          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? 'Signing In' : 'Enter Dashboard'}
          </Button>
        </form>
        {googleClientID ? (
          <div style={{ marginTop: 16 }}>
            <div className="label" style={{ marginBottom: 8 }}>
              Or continue with Google
            </div>
            <div ref={googleButtonRef} />
            {!googleReady ? <div style={{ marginTop: 8, color: 'var(--text-secondary)' }}>Loading Google sign-in...</div> : null}
          </div>
        ) : null}
        <div style={{ marginTop: 18, color: 'var(--text-secondary)' }}>
          Need an account? <Link to="/register">Register</Link>
        </div>
      </div>
    </main>
  );
}
