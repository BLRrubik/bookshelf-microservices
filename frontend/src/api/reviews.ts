import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { booksApi } from './client';
import type { 
  Review,
  ReviewListResponse,
  ReviewsParams,
  CreateReviewRequest,
  UpdateReviewRequest
} from '@/types/api';

export function useBookReviews(bookId: string, params: ReviewsParams = {}) {
  return useQuery({
    queryKey: ['reviews', bookId, params],
    queryFn: async () => {
      const response = await booksApi.get<ReviewListResponse>(
        `/api/v1/books/${bookId}/reviews`, 
        { params }
      );
      return response.data;
    },
    enabled: !!bookId,
    retry: false,
  });
}

export function useReview(id: string) {
  return useQuery({
    queryKey: ['reviews', 'single', id],
    queryFn: async () => {
      const response = await booksApi.get<Review>(`/api/v1/reviews/${id}`);
      return response.data;
    },
    enabled: !!id,
    retry: false,
  });
}

export function useCreateReview(bookId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateReviewRequest) => {
      const response = await booksApi.post<Review>(
        `/api/v1/books/${bookId}/reviews`, 
        data
      );
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews', bookId] });
      queryClient.invalidateQueries({ queryKey: ['books', bookId] });
      queryClient.invalidateQueries({ queryKey: ['books'] });
      toast.success('Рецензия добавлена!');
    },
    onError: (error: any) => {
      const message = error.response?.status === 409 
        ? 'Вы уже оставили рецензию на эту книгу'
        : 'Ошибка добавления рецензии';
      toast.error(message);
    },
  });
}

export function useUpdateReview(reviewId: string, bookId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: UpdateReviewRequest) => {
      const response = await booksApi.put<Review>(
        `/api/v1/reviews/${reviewId}`, 
        data
      );
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews', bookId] });
      queryClient.invalidateQueries({ queryKey: ['reviews', 'single', reviewId] });
      queryClient.invalidateQueries({ queryKey: ['books', bookId] });
      queryClient.invalidateQueries({ queryKey: ['books'] });
      toast.success('Рецензия обновлена');
    },
    onError: () => {
      toast.error('Ошибка обновления рецензии');
    },
  });
}

export function useDeleteReview(bookId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (reviewId: string) => {
      await booksApi.delete(`/api/v1/reviews/${reviewId}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews', bookId] });
      queryClient.invalidateQueries({ queryKey: ['books', bookId] });
      queryClient.invalidateQueries({ queryKey: ['books'] });
      toast.success('Рецензия удалена');
    },
    onError: () => {
      toast.error('Ошибка удаления рецензии');
    },
  });
}

export function useUserReviews(userId: string, params: ReviewsParams = {}) {
  return useQuery({
    queryKey: ['reviews', 'user', userId, params],
    queryFn: async () => {
      const response = await booksApi.get<ReviewListResponse>(
        `/api/v1/users/${userId}/reviews`,
        { params }
      );
      return response.data;
    },
    enabled: !!userId,
  });
}
