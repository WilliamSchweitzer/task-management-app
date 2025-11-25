import { create } from 'zustand';
import { authService } from '@/services/auth-service';
import type { User, LoginCredentials, SignupCredentials } from '@/types';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isCheckingAuth: boolean; // For initial auth check
  isLoading: boolean; // For login/signup actions
  error: string | null;
  login: (credentials: LoginCredentials) => Promise<void>;
  signup: (credentials: SignupCredentials) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isCheckingAuth: true, // Start as true for initial check
  isLoading: false, // Start as false for form actions
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
    try {
      await authService.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      set({
        user: null,
        isAuthenticated: false,
        error: null,
      });
    }
  },

checkAuth: async () => {
  set({ isCheckingAuth: true });
  
  const storedUser = authService.getStoredUser();
  const hasToken = authService.isAuthenticated();
  
  if (!hasToken || !storedUser) {
    set({
      user: null,
      isAuthenticated: false,
      isCheckingAuth: false,
    });
    return;
  }

  // Add a timeout - if verification takes > 5 seconds, just trust local storage
  const timeoutPromise = new Promise((resolve) => setTimeout(resolve, 5000));
  
  try {
    await Promise.race([
      authService.verifyToken(),
      timeoutPromise
    ]);
    
    set({
      user: storedUser,
      isAuthenticated: true,
      isCheckingAuth: false,
    });
  } catch (error) {
    console.error('Auth verification error:', error);
    set({
      user: storedUser,
      isAuthenticated: true,
      isCheckingAuth: false,
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