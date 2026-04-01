import { create } from 'zustand';

type Theme = 'light' | 'dark';

interface ThemeState {
  theme: Theme;
  toggleTheme: () => void;
}

const getInitialTheme = (): Theme => {
  const stored = localStorage.getItem('linear-lite-theme');
  return stored === 'dark' ? 'dark' : 'light';
};

export const useThemeStore = create<ThemeState>((set) => ({
  theme: getInitialTheme(),
  toggleTheme: () =>
    set((state) => {
      const next = state.theme === 'light' ? 'dark' : 'light';
      localStorage.setItem('linear-lite-theme', next);
      return { theme: next };
    }),
}));
