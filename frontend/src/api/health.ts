import { useQuery } from '@tanstack/react-query';
import { gatewayApi } from './client';
import type { HealthResponse } from '@/types/api';

export function useGatewayHealth() {
  return useQuery({
    queryKey: ['health', 'gateway'],
    queryFn: async () => {
      const response = await gatewayApi.get<HealthResponse>('/health');
      return response.data;
    },
    refetchInterval: 30000,
    retry: false,
  });
}

// Aliases for backward compatibility
export const useAuthHealth = useGatewayHealth;
export const useBooksHealth = useGatewayHealth;
