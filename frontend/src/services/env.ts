export const isMockDataEnabled = import.meta.env.DEV && import.meta.env.VITE_USE_MOCK_DATA === 'true';

export function shouldUseMockForPath(path: string) {
  if (!isMockDataEnabled) {
    return false;
  }

  // Milestone 2 auth must use the real backend even when other domains are still mocked.
  return !path.startsWith('/auth/');
}
