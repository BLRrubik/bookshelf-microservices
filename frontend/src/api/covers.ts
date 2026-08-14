import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { gatewayApi } from './client';
import type { CoverUploadResponse, CoverStatusResponse } from '@/types/api';

export function useUploadCover(bookId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('file', file);
      
      const response = await gatewayApi.post<CoverUploadResponse>(
        `/api/v1/books/${bookId}/cover`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        }
      );
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books', bookId] });
      queryClient.invalidateQueries({ queryKey: ['cover-status', bookId] });
      toast.success('Обложка загружена и обрабатывается');
    },
    onError: (error: any) => {
      const message = error.response?.data?.message || 'Ошибка загрузки обложки';
      toast.error(message);
    },
  });
}

export function useCoverStatus(bookId: string, initialStatus?: string) {
  return useQuery({
    queryKey: ['cover-status', bookId],
    queryFn: async () => {
      const response = await gatewayApi.get<CoverStatusResponse>(
        `/api/v1/books/${bookId}/cover/status`
      );
      return response.data;
    },
    // Conditional polling: poll while processing
    refetchInterval: (query) => {
      const data = query.state.data;
      return data?.status === 'processing' ? 2000 : false;
    },
    // Only query if book exists and has a cover (not 'none')
    enabled: !!bookId && initialStatus !== 'none',
    // Don't retry on 404 (cover not found)
    retry: (failureCount, error: any) => {
      if (error?.response?.status === 404) return false;
      return failureCount < 3;
    },
  });
}

export function useDeleteCover(bookId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await gatewayApi.delete(`/api/v1/books/${bookId}/cover`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books', bookId] });
      queryClient.invalidateQueries({ queryKey: ['cover-status', bookId] });
      toast.success('Обложка удалена');
    },
    onError: () => {
      toast.error('Ошибка удаления обложки');
    },
  });
}
