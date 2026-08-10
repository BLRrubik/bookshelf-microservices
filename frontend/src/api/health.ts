import { useQuery } from '@tanstack/react-query';
import { authApi, booksApi } from './client';
import type { ReadyResponse } from '@/types/api';

export function useAuthHealth() {
  return useQuery({
    queryKey: ['health', 'auth'],
    queryFn: async () => {
      const response = await authApi.get<ReadyResponse>('/ready');
      return response.data;
    },
    refetchInterval: 30000, // Poll every 30 seconds
    retry: false,
  });
}

export function useBooksHealth() {
  return useQuery({
    queryKey: ['health', 'books'],
    queryFn: async () => {
      const response = await booksApi.get<ReadyResponse>('/ready');
      return response.data;
    },
    refetchInterval: 30000, // Poll every 30 seconds
    retry: false,
  });
}

