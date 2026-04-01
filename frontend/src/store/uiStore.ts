import { create } from 'zustand';

export interface ToastItem {
  id: number;
  tone: 'success' | 'error' | 'info';
  message: string;
}

interface UIState {
  createIssueOpen: boolean;
  toasts: ToastItem[];
  openCreateIssue: () => void;
  closeCreateIssue: () => void;
  pushToast: (toast: Omit<ToastItem, 'id'>) => void;
  removeToast: (id: number) => void;
  bindShortcuts: () => () => void;
}

let toastId = 0;

export const useUIStore = create<UIState>((set, get) => ({
  createIssueOpen: false,
  toasts: [],
  openCreateIssue: () => set({ createIssueOpen: true }),
  closeCreateIssue: () => set({ createIssueOpen: false }),
  pushToast: (toast) => {
    const id = ++toastId;
    set((state) => ({ toasts: [...state.toasts, { ...toast, id }] }));
    window.setTimeout(() => get().removeToast(id), 3500);
  },
  removeToast: (id) => set((state) => ({ toasts: state.toasts.filter((toast) => toast.id !== id) })),
  bindShortcuts: () => {
    const handler = (event: KeyboardEvent) => {
      if (event.key.toLowerCase() === 'c' && !event.metaKey && !event.ctrlKey) {
        const target = event.target as HTMLElement | null;
        const tag = target?.tagName;
        if (!tag || !['INPUT', 'TEXTAREA', 'SELECT'].includes(tag)) {
          get().openCreateIssue();
        }
      }
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  },
}));
