import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { booksApi } from './client';
import type { 
  Book,
  BookListResponse,
  BooksParams,
  CreateBookRequest,
  UpdateBookRequest
} from '@/types/api';

export function useBooks(params: BooksParams = {}) {
  return useQuery({
    queryKey: ['books', params],
    queryFn: async () => {
      const response = await booksApi.get<BookListResponse>('/api/v1/books', { params });
      return response.data;
    },
    retry: false,
  });
}

export function useBook(id: string) {
  return useQuery({
    queryKey: ['books', id],
    queryFn: async () => {
      const response = await booksApi.get<Book>(`/api/v1/books/${id}`);
      return response.data;
    },
    enabled: !!id,
    retry: false,
  });
}

export function useCreateBook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateBookRequest) => {
      const response = await booksApi.post<Book>('/api/v1/books', data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
      toast.success('Книга добавлена!');
    },
    onError: () => {
      toast.error('Ошибка добавления книги');
    },
  });
}

export function useUpdateBook(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: UpdateBookRequest) => {
      const response = await booksApi.put<Book>(`/api/v1/books/${id}`, data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
      queryClient.invalidateQueries({ queryKey: ['books', id] });
      toast.success('Книга обновлена');
    },
    onError: () => {
      toast.error('Ошибка обновления книги');
    },
  });
}

export function useDeleteBook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await booksApi.delete(`/api/v1/books/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
      toast.success('Книга удалена');
    },
    onError: () => {
      toast.error('Ошибка удаления книги');
    },
  });
}
