/**
 * Маппинг функций на этапы проекта 3 (Event-Driven Architecture)
 * 
 * Этот проект добавляет асинхронную обработку обложек:
 * - books-service публикует сообщения в RabbitMQ
 * - worker-service обрабатывает изображения
 * - MinIO хранит файлы
 * 
 * Базовые сервисы (auth, books, reviews) уже готовы из Project 2.
 */

export interface StageInfo {
  stage: number;
  name: string;
  description: string;
  hint: string;
  icon: string;
  service: 'auth' | 'books' | 'worker';
}

export const FEATURE_STAGES: Record<string, StageInfo> = {
  auth: {
    stage: 0,
    name: 'Авторизация',
    description: 'Готово из Project 2',
    hint: 'auth-service уже реализован',
    icon: '👤',
    service: 'auth',
  },
  books: {
    stage: 0,
    name: 'Каталог книг',
    description: 'Готово из Project 2',
    hint: 'books-service уже реализован',
    icon: '📚',
    service: 'books',
  },
  reviews: {
    stage: 0,
    name: 'Рецензии',
    description: 'Готово из Project 2',
    hint: 'reviews уже реализованы',
    icon: '⭐',
    service: 'books',
  },
  coverUpload: {
    stage: 13,
    name: 'Загрузка обложки',
    description: 'Загрузка изображения обложки книги',
    hint: 'Реализуйте POST /api/v1/books/{id}/cover в CoverHandler',
    icon: '📤',
    service: 'books',
  },
  coverStatus: {
    stage: 18,
    name: 'Статус обложки',
    description: 'Получение статуса обработки и URL обложки',
    hint: 'Реализуйте GET /api/v1/books/{id}/cover/status в CoverHandler',
    icon: '🔄',
    service: 'books',
  },
  coverProcessing: {
    stage: 15,
    name: 'Обработка обложки',
    description: 'Worker обрабатывает изображения (resize, thumbnail)',
    hint: 'Реализуйте обработку изображений в worker-service',
    icon: '⚙️',
    service: 'worker',
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
  
  // Network error - сервис не запущен
  if (axiosError.code === 'ERR_NETWORK') return true;
  if (axiosError.code === 'ERR_CONNECTION_REFUSED') return true;
  
  return false;
}

/**
 * Определяет, является ли ошибка сетевой (сервис не запущен)
 */
export function isNetworkError(error: unknown): boolean {
  if (!error || typeof error !== 'object') return false;
  
  const axiosError = error as { code?: string; message?: string };
  
  if (axiosError.code === 'ERR_NETWORK') return true;
  if (axiosError.code === 'ERR_CONNECTION_REFUSED') return true;
  if (axiosError.message?.includes('Network Error')) return true;
  
  return false;
}
