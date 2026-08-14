import { create } from 'zustand';

interface RequestInfo {
  requestId?: string;
  cacheStatus?: string;
  responseTime: number;
  url: string;
  method: string;
  status: number;
  timestamp?: number;
}

interface RequestState {
  lastRequest: RequestInfo | null;
  history: RequestInfo[];
  
  setLastRequest: (info: RequestInfo) => void;
  clearHistory: () => void;
}

export const useRequestStore = create<RequestState>((set) => ({
  lastRequest: null,
  history: [],
  
  setLastRequest: (info) => 
    set((state) => ({
      lastRequest: { ...info, timestamp: Date.now() },
      history: [
        { ...info, timestamp: Date.now() },
        ...state.history.slice(0, 49), // Keep last 50
      ],
    })),
  
  clearHistory: () => 
    set({ history: [] }),
}));





