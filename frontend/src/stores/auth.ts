import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '@/types/api.ts';

let sessionExpiryTimer: ReturnType<typeof setTimeout> | null = null;

function clearSessionExpiryTimer() {
  if (sessionExpiryTimer !== null) {
    clearTimeout(sessionExpiryTimer);
    sessionExpiryTimer = null;
  }
}

function scheduleSessionExpiry(expiresAt: number, logout: () => void) {
  clearSessionExpiryTimer();
  const remaining = expiresAt - Date.now();

  if (remaining <= 0) {
    logout();
    return;
  }

  sessionExpiryTimer = setTimeout(logout, remaining);
}

interface AuthState {
  token: string | null;
  user: User | null;
  expiresAt: number | null;
  isAuthenticated: boolean;
  setAuth: (token: string, user: User, expiresIn: number) => void;
  setUser: (user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      expiresAt: null,
      isAuthenticated: false,
      
      setAuth: (token, user, expiresIn) => {
        const expiresAt = Date.now() + expiresIn * 1000;
        set({ token, user, expiresAt, isAuthenticated: true });
        scheduleSessionExpiry(expiresAt, () => useAuthStore.getState().logout());
      },
      
      setUser: (user) => 
        set({ user }),
      
      logout: () => {
        clearSessionExpiryTimer();
        set({ token: null, user: null, expiresAt: null, isAuthenticated: false });
      },
    }),
    {
      name: 'bookshelf-auth-p2',
      partialize: (state) => ({ 
        token: state.token,
        user: state.user,
        expiresAt: state.expiresAt,
        isAuthenticated: state.isAuthenticated 
      }),
      onRehydrateStorage: () => (state) => {
        if (!state?.isAuthenticated) {
          return;
        }

        if (!state.token || !state.user || !state.expiresAt || state.expiresAt <= Date.now()) {
          state.logout();
          return;
        }

        scheduleSessionExpiry(state.expiresAt, () => useAuthStore.getState().logout());
      },
    }
  )
);
