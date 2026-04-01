import { useEffect } from 'react';
import { authApi } from 'services/authApi';
import { useAuthStore } from 'store/authStore';

export function useAuthBootstrap() {
  const token = useAuthStore((state) => state.token);
  const setUser = useAuthStore((state) => state.setUser);
  const clearSession = useAuthStore((state) => state.clearSession);
  const setBootstrapped = useAuthStore((state) => state.setBootstrapped);

  useEffect(() => {
    let cancelled = false;

    async function bootstrap() {
      if (!token) {
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
  }, [clearSession, setBootstrapped, setUser, token]);
}
