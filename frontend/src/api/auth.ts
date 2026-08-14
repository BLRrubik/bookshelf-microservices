import { useMutation, useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { gatewayApi } from './client';
import { useAuthStore } from '@/stores/auth';
import type { 
  AuthResponse, 
  LoginRequest, 
  RegisterRequest,
  User,
  PublicUser,
  UpdateUserRequest
} from '@/types/api';

export function useLogin() {
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);

  return useMutation({
    mutationFn: async (data: LoginRequest) => {
      const response = await gatewayApi.post<AuthResponse>('/api/v1/auth/login', data);
      return response.data;
    },
    onSuccess: (data) => {
      setAuth(data.access_token, data.refresh_token, data.user);
      toast.success('Добро пожаловать!');
      navigate('/');
    },
    onError: () => {
      toast.error('Неверный email или пароль');
    },
  });
}

export function useRegister() {
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);

  return useMutation({
    mutationFn: async (data: RegisterRequest) => {
      const response = await gatewayApi.post<AuthResponse>('/api/v1/auth/register', data);
      return response.data;
    },
    onSuccess: (data) => {
      setAuth(data.access_token, data.refresh_token, data.user);
      toast.success('Регистрация успешна!');
      navigate('/');
    },
    onError: (error: any) => {
      const message = error.response?.data?.message || 'Ошибка регистрации';
      toast.error(message);
    },
  });
}

export function useLogout() {
  const navigate = useNavigate();
  const logout = useAuthStore((state) => state.logout);

  return useMutation({
    mutationFn: async () => {
      await gatewayApi.post('/api/v1/auth/logout');
    },
    onSuccess: () => {
      logout();
      navigate('/');
      toast.success('Вы вышли из аккаунта');
    },
    onError: () => {
      logout();
      navigate('/');
    },
  });
}

export function useCurrentUser() {
  const token = useAuthStore((state) => state.token);
  
  return useQuery({
    queryKey: ['currentUser'],
    queryFn: async () => {
      const response = await gatewayApi.get<User>('/api/v1/users/me');
      return response.data;
    },
    enabled: !!token,
  });
}

export function useUpdateProfile() {
  const setUser = useAuthStore((state) => state.setUser);

  return useMutation({
    mutationFn: async (data: UpdateUserRequest) => {
      const response = await gatewayApi.put<User>('/api/v1/users/me', data);
      return response.data;
    },
    onSuccess: (data) => {
      setUser(data);
      toast.success('Профиль обновлён');
    },
    onError: () => {
      toast.error('Ошибка обновления профиля');
    },
  });
}

export function useUser(userId: string) {
  return useQuery({
    queryKey: ['users', userId],
    queryFn: async () => {
      const response = await gatewayApi.get<PublicUser>(`/api/v1/users/${userId}`);
      return response.data;
    },
    enabled: !!userId,
  });
}
