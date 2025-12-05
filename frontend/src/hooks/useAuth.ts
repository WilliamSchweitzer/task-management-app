'use client';

import { useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuthStore } from '@/stores/auth-store';

const PUBLIC_ROUTES = ['/login', '/signup', '/'];

export function useAuth() {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, isLoading, checkAuth, user } = useAuthStore();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  useEffect(() => {
    // Skip redirects while still loading
    if (isLoading) return;

    const isPublicRoute = PUBLIC_ROUTES.includes(pathname);

    if (!isAuthenticated && !isPublicRoute) {
      // Redirect to login if not authenticated and trying to access protected route
      router.push('/login');
    } else if (isAuthenticated && isPublicRoute && pathname !== '/') {
      // Redirect to dashboard if authenticated and trying to access auth pages
      router.push('/dashboard');
    }
  }, [isAuthenticated, isLoading, pathname, router]);

  return {
    isAuthenticated,
    isLoading,
    user,
  };
}
