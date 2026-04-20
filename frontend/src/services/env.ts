export const isMockDataEnabled = import.meta.env.DEV && import.meta.env.VITE_USE_MOCK_DATA === 'true';

export function getApiBaseUrl() {
  const runtimeBase = window.__APP_CONFIG__?.API_BASE_URL?.trim();
  if (runtimeBase) {
    return runtimeBase;
  }

  const buildBase = import.meta.env.VITE_API_BASE_URL?.trim();
  if (buildBase) {
    return buildBase;
  }

  return '/api/v1';
}

export function getGoogleClientId() {
  const runtimeClientId = window.__APP_CONFIG__?.GOOGLE_CLIENT_ID?.trim();
  if (runtimeClientId) {
    return runtimeClientId;
  }

  return import.meta.env.VITE_GOOGLE_CLIENT_ID?.trim() ?? '';
}

export function shouldUseMockForPath(path: string) {
  if (!isMockDataEnabled) {
    return false;
  }

  // Milestone 2 auth must use the real backend even when other domains are still mocked.
  return !path.startsWith('/auth/');
}
