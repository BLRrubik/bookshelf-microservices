import { useQuery } from '@tanstack/react-query';
import { gatewayApi } from './client';
import type { DashboardResponse, SearchResponse, SearchParams } from '@/types/api';

export function useDashboard() {
  return useQuery({
    queryKey: ['dashboard'],
    queryFn: async () => {
      const response = await gatewayApi.get<DashboardResponse>('/api/v1/dashboard');
      return {
        data: response.data,
        cacheStatus: response.headers['x-cache'],
      };
    },
    staleTime: 1000 * 60 * 5, // 5 minutes - leverage gateway cache
  });
}

export function useSearch(params: SearchParams) {
  return useQuery({
    queryKey: ['search', params],
    queryFn: async () => {
      const response = await gatewayApi.get<SearchResponse>('/api/v1/search', { 
        params: {
          q: params.q,
          limit: params.limit || 20,
        }
      });
      return response.data;
    },
    enabled: params.q.length >= 2, // Only search with 2+ chars
  });
}





