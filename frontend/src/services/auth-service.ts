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

    // Calculate expires_at from expires_in (seconds from now)
    const expiresAt = Math.floor(Date.now() / 1000) + response.expires_in;

    // Store tokens and user
    apiClient.storeTokens(
      response.access_token,
      response.refresh_token,
      expiresAt
    );
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(response.user));

    return response;
  },

  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(
      API_CONFIG.AUTH.LOGIN,
      credentials
    );

    // Calculate expires_at from expires_in (seconds from now)
    const expiresAt = Math.floor(Date.now() / 1000) + response.expires_in;

    // Store tokens and user
    apiClient.storeTokens(
      response.access_token,
      response.refresh_token,
      expiresAt
    );
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(response.user));

    return response;
  },

  async logout(): Promise<void> {
  try {
    const refreshToken = localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
    if (refreshToken) {
      await apiClient.post(API_CONFIG.AUTH.LOGOUT, {
        refresh_token: refreshToken,  // Send refresh token in body
      });
    }
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

  async refreshToken(): Promise<AuthResponse> {
  const refreshToken = localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
  const user = this.getStoredUser();
  
  if (!refreshToken || !user?.email) {
    throw new Error('No refresh token or email available');
  }

  const response = await apiClient.post<AuthResponse>(
    API_CONFIG.AUTH.REFRESH,
    {
      email: user.email,
      refresh_token: refreshToken,
    }
  );

  // Calculate expires_at from expires_in
  const expiresAt = Math.floor(Date.now() / 1000) + response.expires_in;

  // Store new tokens
  apiClient.storeTokens(
    response.access_token,
    response.refresh_token,
    expiresAt
  );
  
  if (response.user) {
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(response.user));
  }

  return response;
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