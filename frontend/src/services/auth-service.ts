import { apiClient } from '@/lib/api-client';
import { API_CONFIG, STORAGE_KEYS } from '@/lib/config';
import {
  AuthResponse,
  LoginCredentials,
  SignupCredentials,
  User,
} from '@/types';

export const authService = {
  async signup(credentials: SignupCredentials): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(
      API_CONFIG.AUTH.SIGNUP,
      credentials
    );

    // Store tokens and user
    apiClient.storeTokens(
      response.tokens.access_token,
      response.tokens.refresh_token,
      response.tokens.expires_at
    );
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(response.user));

    return response;
  },

  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(
      API_CONFIG.AUTH.LOGIN,
      credentials
    );

    // Store tokens and user
    apiClient.storeTokens(
      response.tokens.access_token,
      response.tokens.refresh_token,
      response.tokens.expires_at
    );
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(response.user));

    return response;
  },

  async logout(): Promise<void> {
    try {
      await apiClient.post(API_CONFIG.AUTH.LOGOUT);
    } catch {
      // Continue with logout even if API call fails
    } finally {
      apiClient.removeTokens();
    }
  },

  async verifyToken(): Promise<User | null> {
    try {
      const response = await apiClient.get<{ user: User }>(API_CONFIG.AUTH.VERIFY);
      return response.user;
    } catch {
      return null;
    }
  },

  getStoredUser(): User | null {
    if (typeof window === 'undefined') return null;
    const userJson = localStorage.getItem(STORAGE_KEYS.USER);
    if (!userJson) return null;
    try {
      return JSON.parse(userJson) as User;
    } catch {
      return null;
    }
  },

  isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false;
    const token = localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
    return !!token;
  },
};
