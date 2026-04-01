import { create } from 'zustand';
import { UserSummary } from 'types/domain';

const TOKEN_KEY = 'linear-lite-token';

interface AuthState {
  token: string | null;
  user: UserSummary | null;
  bootstrapped: boolean;
  setSession: (token: string, user: UserSummary) => void;
  clearSession: () => void;
  setUser: (user: UserSummary | null) => void;
  setBootstrapped: (bootstrapped: boolean) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem(TOKEN_KEY),
  user: null,
  bootstrapped: false,
  setSession: (token, user) => {
    localStorage.setItem(TOKEN_KEY, token);
    set({ token, user, bootstrapped: true });
  },
  clearSession: () => {
    localStorage.removeItem(TOKEN_KEY);
    set({ token: null, user: null, bootstrapped: true });
  },
  setUser: (user) => set({ user }),
  setBootstrapped: (bootstrapped) => set({ bootstrapped }),
}));

export function getStoredToken() {
  return localStorage.getItem(TOKEN_KEY);
}
