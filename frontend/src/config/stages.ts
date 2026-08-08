/**
 * Маппинг функций на этапы проекта 2 (Microservices Decomposition)
 * 
 * В этом проекте два независимых сервиса:
 * - auth-service (порт 8081) — авторизация пользователей
 * - books-service (порт 8082) — книги и рецензии
 * 
 * Каждая функция становится доступной после реализации соответствующего этапа.
 */

export interface StageInfo {
  stage: number;
  name: string;
  description: string;
  hint: string;
  icon: string;
  service: 'auth' | 'books';
}

export const FEATURE_STAGES: Record<string, StageInfo> = {
  authHealth: {
    stage: 5,
    name: 'Auth Service Health',
    description: 'Проверка работоспособности auth-service',
    hint: 'Реализуйте GET /health в auth-service',
    icon: '🟢',
    service: 'auth',
  },
  booksHealth: {
    stage: 8,
    name: 'Books Service Health',
    description: 'Проверка работоспособности books-service',
    hint: 'Реализуйте GET /health в books-service',
    icon: '🟢',
    service: 'books',
  },
  auth: {
    stage: 6,
    name: 'Авторизация',
    description: 'Регистрация и вход в систему через auth-service',
    hint: 'Реализуйте POST /api/v1/auth/register и POST /api/v1/auth/login в auth-service',
    icon: '👤',
    service: 'auth',
  },
  books: {
    stage: 11,
    name: 'Каталог книг',
    description: 'Просмотр, создание и редактирование книг через books-service',
    hint: 'Реализуйте BookHandler в books-service (GET/POST /api/v1/books)',
    icon: '📚',
    service: 'books',
  },
  reviews: {
    stage: 12,
    name: 'Рецензии',
    description: 'Просмотр и написание рецензий через books-service',
    hint: 'Реализуйте ReviewHandler в books-service (GET/POST /api/v1/books/{id}/reviews)',
    icon: '⭐',
    service: 'books',
  },
  userInfo: {
    stage: 16,
    name: 'Информация о пользователях',
    description: 'Отображение имён авторов рецензий',
    hint: 'Реализуйте межсервисный вызов auth-service из books-service',
    icon: '🔗',
    service: 'books',
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

/**
 * Получить информацию о сервисе по названию функции
 */
export function getServiceInfo(feature: keyof typeof FEATURE_STAGES) {
  const info = FEATURE_STAGES[feature];
  return {
    ...info,
    port: info.service === 'auth' ? 8081 : 8082,
    serviceName: info.service === 'auth' ? 'auth-service' : 'books-service',
  };
}
