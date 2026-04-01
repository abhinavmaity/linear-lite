import { ReactNode } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { Spinner } from 'components/common/Spinner';
import { useAuthStore } from 'store/authStore';

export function AuthGate({ requireAuth, children }: { requireAuth: boolean; children?: ReactNode }) {
  const location = useLocation();
  const { token, bootstrapped } = useAuthStore();
  const bypassAuth = import.meta.env.DEV && import.meta.env.VITE_BYPASS_AUTH === 'true';

  if (bypassAuth) {
    return <>{children ?? <Outlet />}</>;
  }

  if (!bootstrapped) {
    return <Spinner fullScreen label="Restoring session" />;
  }

  if (requireAuth && !token) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  if (!requireAuth && token) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children ?? <Outlet />}</>;
}
