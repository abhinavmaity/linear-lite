import { ReactNode } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { Spinner } from 'components/common/Spinner';
import { isMockDataEnabled } from 'services/env';
import { useAuthStore } from 'store/authStore';

export function AuthGate({ requireAuth, children }: { requireAuth: boolean; children?: ReactNode }) {
  const location = useLocation();
  const { token, bootstrapped } = useAuthStore();
  const bypassAuth = isMockDataEnabled;

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
