import { useEffect } from 'react';
import { isMockDataEnabled } from 'services/env';
import { getMockSession } from 'services/mockBackend';
import { authApi } from 'services/authApi';
import { useAuthStore } from 'store/authStore';

export function useAuthBootstrap() {
  const token = useAuthStore((state) => state.token);
  const setSession = useAuthStore((state) => state.setSession);
  const setUser = useAuthStore((state) => state.setUser);
  const clearSession = useAuthStore((state) => state.clearSession);
  const setBootstrapped = useAuthStore((state) => state.setBootstrapped);

  useEffect(() => {
    let cancelled = false;

    async function bootstrap() {
      if (!token) {
        if (isMockDataEnabled) {
          const session = getMockSession();
          setSession(session.token, session.user);
          return;
        }
        setBootstrapped(true);
        return;
      }

      try {
        const response = await authApi.me();
        if (!cancelled) {
          setUser(response.data);
          setBootstrapped(true);
        }
      } catch {
        if (!cancelled) {
          clearSession();
        }
      }
    }

    bootstrap();

    return () => {
      cancelled = true;
    };
  }, [clearSession, setBootstrapped, setSession, setUser, token]);
}
