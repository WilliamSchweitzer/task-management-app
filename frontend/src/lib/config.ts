// API Configuration - Update these with your actual Kong Gateway endpoints
export const API_CONFIG = {
  // Base URL for your Kong Gateway
  BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8000',
  
  // Auth service endpoints
  AUTH: {
    SIGNUP: '/auth/signup',
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    VERIFY: '/auth/verify',
  },
  
  // Task service endpoints
  TASKS: {
    LIST: '/tasks',
    CREATE: '/tasks',
    GET: (id: string) => `/tasks/${id}`,
    UPDATE: (id: string) => `/tasks/${id}`,
    DELETE: (id: string) => `/tasks/${id}`,
    COMPLETE: (id: string) => `/tasks/${id}/complete`,
  },
};

// Storage keys
export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'task_app_access_token',
  REFRESH_TOKEN: 'task_app_refresh_token',
  USER: 'task_app_user',
  TOKEN_EXPIRY: 'task_app_token_expiry',
} as const;
