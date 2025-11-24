'use client';

import { create } from 'zustand';
import { User, LoginCredentials, SignupCredentials } from '@/types';
import { authService } from '@/services/auth-service';

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  error: string | null;
  
  // Actions
  login: (credentials: LoginCredentials) => Promise<void>;
  signup: (credentials: SignupCredentials) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  isLoading: true,
  isAuthenticated: false,
  error: null,

  login: async (credentials: LoginCredentials) => {
    set({ isLoading: true, error: null });
    try {
      const response = await authService.login(credentials);
      set({
        user: response.user,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (err) {
      const error = err as { message?: string };
      set({
        isLoading: false,
        error: error.message || 'Login failed. Please check your credentials.',
      });
      throw err;
    }
  },

  signup: async (credentials: SignupCredentials) => {
    set({ isLoading: true, error: null });
    try {
      const response = await authService.signup(credentials);
      set({
        user: response.user,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (err) {
      const error = err as { message?: string };
      set({
        isLoading: false,
        error: error.message || 'Signup failed. Please try again.',
      });
      throw err;
    }
  },

  logout: async () => {
    set({ isLoading: true });
    try {
      await authService.logout();
    } finally {
      set({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: null,
      });
    }
  },

  checkAuth: async () => {
    // Don't check if already loading
    if (get().isLoading && get().user !== null) return;
    
    set({ isLoading: true });
    
    // First check if we have stored user
    const storedUser = authService.getStoredUser();
    const hasToken = authService.isAuthenticated();
    
    if (!hasToken || !storedUser) {
      set({
        user: null,
        isAuthenticated: false,
        isLoading: false,
      });
      return;
    }

    // Verify token with server
    try {
      const user = await authService.verifyToken();
      if (user) {
        set({
          user,
          isAuthenticated: true,
          isLoading: false,
        });
      } else {
        set({
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });
      }
    } catch {
      // Token verification failed, but keep user logged in if we have stored data
      // This allows offline functionality
      set({
        user: storedUser,
        isAuthenticated: true,
        isLoading: false,
      });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));

// Listen for token expiration events
if (typeof window !== 'undefined') {
  window.addEventListener('auth:token-expired', () => {
    useAuthStore.getState().logout();
  });
}
