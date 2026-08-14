import { create } from 'zustand';

interface RateLimitState {
  limit: number;
  remaining: number;
  reset: number | null;
  isWarning: boolean;
  isBlocked: boolean;
  retryAfter: number;
  
  updateRateLimit: (limit: number, remaining: number, reset?: number) => void;
  setWarning: (isWarning: boolean) => void;
  setBlocked: (retryAfter: number) => void;
  clearBlocked: () => void;
}

export const useRateLimitStore = create<RateLimitState>((set, get) => ({
  limit: 100,
  remaining: 100,
  reset: null,
  isWarning: false,
  isBlocked: false,
  retryAfter: 0,
  
  updateRateLimit: (limit, remaining, reset) => 
    set({ limit, remaining, reset: reset || null }),
  
  setWarning: (isWarning) => 
    set({ isWarning }),
  
  setBlocked: (retryAfter) => {
    set({ isBlocked: true, retryAfter });
    
    // Auto-clear after retry period
    setTimeout(() => {
      get().clearBlocked();
    }, retryAfter * 1000);
  },
  
  clearBlocked: () => 
    set({ isBlocked: false, retryAfter: 0, remaining: 100, isWarning: false }),
}));





