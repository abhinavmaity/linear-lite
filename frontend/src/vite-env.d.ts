/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string;
  readonly VITE_USE_MOCK_DATA?: 'true' | 'false';
}

interface Window {
  __APP_CONFIG__?: {
    API_BASE_URL?: string;
  };
}
