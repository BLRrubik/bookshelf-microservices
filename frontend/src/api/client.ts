import axios, { AxiosInstance } from 'axios';
import { useAuthStore } from '@/stores/auth.ts';

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

  // Clear persisted auth state when the access token is no longer valid.
  client.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        useAuthStore.getState().logout();
      }

      return Promise.reject(error);
    }
  );

  return client;
}

// Auth service client
export const authApi = createApiClient(AUTH_URL);

// Books service client  
export const booksApi = createApiClient(BOOKS_URL);

// Legacy export for compatibility
export const api = booksApi;
