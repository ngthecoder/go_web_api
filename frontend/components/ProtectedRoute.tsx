'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect, ReactNode } from 'react';

interface ProtectedRouteProps {
  children: ReactNode;
  redirectTo?: string;
  showLoading?: boolean;
}

export default function ProtectedRoute({ 
  children,
  redirectTo = '/login',
  showLoading = true
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth();

  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      localStorage.setItem('redirectAfterLogin', window.location.pathname);
      router.push(redirectTo);
    }
  }, [isLoading, isAuthenticated, router, redirectTo]);

  if (isLoading) {
    return showLoading ? (
      <div className="flex justify-center items-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    ) : null;
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}
