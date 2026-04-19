/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string;
  readonly VITE_USE_MOCK_DATA?: 'true' | 'false';
  readonly VITE_GOOGLE_CLIENT_ID?: string;
}

interface Window {
  __APP_CONFIG__?: {
    API_BASE_URL?: string;
  };
  google?: {
    accounts?: {
      id?: {
        initialize: (config: { client_id: string; callback: (response: { credential?: string }) => void }) => void;
        renderButton: (
          parent: HTMLElement,
          options: { theme?: string; size?: string; width?: number | string },
        ) => void;
      };
    };
  };
}
