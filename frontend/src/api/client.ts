import axios, { AxiosInstance } from 'axios';
import { useAuthStore } from '@/stores/auth';

// Project 2: Multi-service configuration
const AUTH_URL = import.meta.env.VITE_AUTH_URL || 'http://localhost:8081';
const BOOKS_URL = import.meta.env.VITE_BOOKS_URL || 'http://localhost:8082';

function createApiClient(baseURL: string): AxiosInstance {
  const client = axios.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Request interceptor - add auth token
  client.interceptors.request.use(
    (config) => {
      const token = useAuthStore.getState().token;
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  return client;
}

// Auth service client
export const authApi = createApiClient(AUTH_URL);

// Books service client  
export const booksApi = createApiClient(BOOKS_URL);

// Token refresh interceptor for booksApi
booksApi.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      const refreshToken = useAuthStore.getState().refreshToken;
      if (!refreshToken) {
        useAuthStore.getState().logout();
        return Promise.reject(error);
      }

      try {
        const { data } = await authApi.post('/api/v1/auth/refresh', {
          refresh_token: refreshToken,
        });
        
        useAuthStore.getState().setAuth(
          data.access_token,
          data.refresh_token,
          data.user
        );
        
        originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
        return booksApi(originalRequest);
      } catch {
        useAuthStore.getState().logout();
        return Promise.reject(error);
      }
    }
    
    return Promise.reject(error);
  }
);

// Token refresh interceptor for authApi
authApi.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    // Don't retry refresh endpoint itself
    if (originalRequest.url?.includes('/auth/refresh')) {
      return Promise.reject(error);
    }
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      const refreshToken = useAuthStore.getState().refreshToken;
      if (!refreshToken) {
        useAuthStore.getState().logout();
        return Promise.reject(error);
      }

      try {
        const { data } = await authApi.post('/api/v1/auth/refresh', {
          refresh_token: refreshToken,
        });
        
        useAuthStore.getState().setAuth(
          data.access_token,
          data.refresh_token,
          data.user
        );
        
        originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
        return authApi(originalRequest);
      } catch {
        useAuthStore.getState().logout();
        return Promise.reject(error);
      }
    }
    
    return Promise.reject(error);
  }
);

// Legacy export for compatibility
export const api = booksApi;
