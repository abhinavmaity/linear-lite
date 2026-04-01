import { QueryClientProvider } from '@tanstack/react-query';
import { ReactNode, useEffect } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { queryClient } from 'services/queryClient';
import { ThemeProvider } from './ThemeProvider';
import { Toaster } from 'components/common/Toaster';
import { useAuthBootstrap } from 'features/auth/useAuthBootstrap';
import { useUIStore } from 'store/uiStore';

function Bootstrap() {
  useAuthBootstrap();
  const { bindShortcuts } = useUIStore();

  useEffect(() => bindShortcuts(), [bindShortcuts]);

  return null;
}

export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <BrowserRouter>
          <Bootstrap />
          {children}
          <Toaster />
        </BrowserRouter>
      </ThemeProvider>
    </QueryClientProvider>
  );
}
