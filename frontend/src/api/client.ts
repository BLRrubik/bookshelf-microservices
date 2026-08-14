import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { useAuthStore } from '@/stores/auth';
import { useRateLimitStore } from '@/stores/rateLimit';
import { useRequestStore } from '@/stores/request';

// Project 4: Single Gateway endpoint
const GATEWAY_URL = import.meta.env.VITE_GATEWAY_URL || 'http://localhost:8000';

function createGatewayClient(): AxiosInstance {
  const client = axios.create({
    baseURL: GATEWAY_URL,
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
      
      // Track request start time for response timing
      config.metadata = { startTime: Date.now() };
      
      return config;
    },
    (error) => Promise.reject(error)
  );

  // Response interceptor - handle rate limits and track metadata
  client.interceptors.response.use(
    (response: AxiosResponse) => {
      // Track rate limit headers
      const remaining = response.headers['x-ratelimit-remaining'];
      const limit = response.headers['x-ratelimit-limit'];
      const reset = response.headers['x-ratelimit-reset'];
      
      if (remaining !== undefined && limit !== undefined) {
        const remainingNum = parseInt(remaining);
        const limitNum = parseInt(limit);
        
        useRateLimitStore.getState().updateRateLimit(
          limitNum,
          remainingNum,
          reset ? parseInt(reset) : undefined
        );
        
        // Warn when running low
        if (remainingNum < 10) {
          useRateLimitStore.getState().setWarning(true);
        } else {
          useRateLimitStore.getState().setWarning(false);
        }
      }
      
      // Track request metadata
      const requestId = response.headers['x-request-id'];
      const cacheStatus = response.headers['x-cache'];
      const startTime = (response.config as any).metadata?.startTime;
      const responseTime = startTime ? Date.now() - startTime : 0;
      
      useRequestStore.getState().setLastRequest({
        requestId,
        cacheStatus,
        responseTime,
        url: response.config.url || '',
        method: response.config.method?.toUpperCase() || 'GET',
        status: response.status,
      });
      
      return response;
    },
    async (error) => {
      const originalRequest = error.config;
      
      // Handle rate limit exceeded (429)
      if (error.response?.status === 429) {
        const retryAfter = error.response.headers['retry-after'];
        useRateLimitStore.getState().setBlocked(
          retryAfter ? parseInt(retryAfter) : 60
        );
        return Promise.reject(error);
      }
      
      // Token refresh on 401
      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;
        
        const refreshToken = useAuthStore.getState().refreshToken;
        if (!refreshToken) {
          useAuthStore.getState().logout();
          return Promise.reject(error);
        }

        try {
          const { data } = await gatewayApi.post('/api/v1/auth/refresh', {
            refresh_token: refreshToken,
          });
          
          useAuthStore.getState().setAuth(
            data.access_token,
            data.refresh_token,
            data.user
          );
          
          originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
          return gatewayApi(originalRequest);
        } catch {
          useAuthStore.getState().logout();
          return Promise.reject(error);
        }
      }
      
      return Promise.reject(error);
    }
  );

  return client;
}

export const gatewayApi = createGatewayClient();

// Aliases for compatibility
export const api = gatewayApi;
export const authApi = gatewayApi;
export const booksApi = gatewayApi;

// Extend axios config type
declare module 'axios' {
  interface AxiosRequestConfig {
    metadata?: {
      startTime: number;
    };
  }
}
