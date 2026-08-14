/**
 * Маппинг функций на этапы проекта 4 (API Gateway)
 * 
 * В этом проекте frontend работает через единую точку входа — API Gateway (порт 8000).
 * Базовые сервисы (auth, books, reviews, covers) уже готовы из предыдущих проектов.
 * 
 * Новые функции:
 * - Rate Limiting (Этапы 13-15)
 * - Кэширование (Этапы 16-18)  
 * - Dashboard (Этапы 19-22)
 * - Глобальный поиск (Этап 23)
 */

export interface StageInfo {
  stage: number;
  name: string;
  description: string;
  hint: string;
  icon: string;
}

export const FEATURE_STAGES: Record<string, StageInfo> = {
  // Базовые — уже готовы из предыдущих проектов
  gateway: {
    stage: 4,
    name: 'API Gateway',
    description: 'Единая точка входа для всех запросов',
    hint: 'Реализуйте проксирующие хендлеры в api-gateway',
    icon: '🚪',
  },
  auth: {
    stage: 0,
    name: 'Авторизация',
    description: 'Готово из Project 1-3',
    hint: 'auth-service уже реализован',
    icon: '👤',
  },
  books: {
    stage: 0,
    name: 'Каталог книг',
    description: 'Готово из Project 1-3',
    hint: 'books-service уже реализован',
    icon: '📚',
  },
  reviews: {
    stage: 0,
    name: 'Рецензии',
    description: 'Готово из Project 1-3',
    hint: 'reviews уже реализованы',
    icon: '⭐',
  },
  covers: {
    stage: 0,
    name: 'Обложки',
    description: 'Готово из Project 3',
    hint: 'cover upload уже реализован',
    icon: '🖼️',
  },
  
  // Новые функции Project 4
  rateLimit: {
    stage: 14,
    name: 'Rate Limiting',
    description: 'Защита API от перегрузки',
    hint: 'Реализуйте RateLimitMiddleware с Redis',
    icon: '⚡',
  },
  caching: {
    stage: 16,
    name: 'Кэширование',
    description: 'Кэширование ответов в Redis',
    hint: 'Реализуйте CacheMiddleware для GET запросов',
    icon: '💾',
  },
  dashboard: {
    stage: 19,
    name: 'Dashboard',
    description: 'Агрегированные данные из нескольких сервисов',
    hint: 'Реализуйте GET /api/v1/dashboard с агрегацией',
    icon: '📊',
  },
  globalSearch: {
    stage: 23,
    name: 'Глобальный поиск',
    description: 'Поиск по всем сервисам одновременно',
    hint: 'Реализуйте GET /api/v1/search с параллельными запросами',
    icon: '🔍',
  },
};

/**
 * Определяет, является ли ошибка признаком нереализованного endpoint
 */
export function isFeatureNotImplemented(error: unknown): boolean {
  if (!error || typeof error !== 'object') return false;
  
  const axiosError = error as { response?: { status?: number }; code?: string };
  
  // 404 - endpoint не существует
  if (axiosError.response?.status === 404) return true;
  
  // 502 - gateway не может достучаться до сервиса
  if (axiosError.response?.status === 502) return true;
  
  // Network error - gateway не запущен
  if (axiosError.code === 'ERR_NETWORK') return true;
  if (axiosError.code === 'ERR_CONNECTION_REFUSED') return true;
  
  return false;
}

/**
 * Определяет, является ли ошибка rate limit
 */
export function isRateLimited(error: unknown): boolean {
  if (!error || typeof error !== 'object') return false;
  
  const axiosError = error as { response?: { status?: number } };
  return axiosError.response?.status === 429;
}

/**
 * Определяет, является ли ошибка сетевой (gateway не запущен)
 */
export function isNetworkError(error: unknown): boolean {
  if (!error || typeof error !== 'object') return false;
  
  const axiosError = error as { code?: string; message?: string };
  
  if (axiosError.code === 'ERR_NETWORK') return true;
  if (axiosError.code === 'ERR_CONNECTION_REFUSED') return true;
  if (axiosError.message?.includes('Network Error')) return true;
  
  return false;
}
