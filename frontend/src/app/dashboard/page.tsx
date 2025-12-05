'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/stores/auth-store';
import { useTaskStore } from '@/stores/task-store';
import { TaskBoard, TaskModal } from '@/components/tasks';

export default function DashboardPage() {
  const router = useRouter();
  const { user, isAuthenticated, isCheckingAuth, logout, checkAuth } = useAuthStore(); 
  const { fetchTasks, isLoading: tasksLoading, error: tasksError, openModal } = useTaskStore();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  // Check authentication but don't redirect based on tasks loading
  useEffect(() => {
    // Only redirect if we're done checking auth and user is not authenticated
    if (!isCheckingAuth && !isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, isCheckingAuth, router]);

  // Fetch tasks separately - failure here should NOT redirect
  useEffect(() => {
    if (isAuthenticated) {
      fetchTasks().catch((error) => {
        // Log error but don't redirect - user should still see the dashboard
        console.error('Failed to fetch tasks:', error);
      });
    }
  }, [isAuthenticated]);

  // Show loading spinner only while checking authentication
  if (isCheckingAuth) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  // If not authenticated after check, return null (redirect will happen)
  if (!isAuthenticated) {
    return null;
  }

  // Render dashboard even if tasks are loading or failed
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
<header className="bg-white dark:bg-gray-800 shadow">
  <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
    {/* Mobile layout: Stack vertically */}
    <div className="flex flex-col gap-3 sm:flex-row sm:justify-between sm:items-center">
      {/* Title section */}
      <div className="flex-shrink-0">
        <h1 className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">
          TaskBoard
        </h1>
        <p className="text-xs sm:text-sm text-gray-600 dark:text-gray-400 mt-1">
          Welcome back, {user?.name || 'User'}
        </p>
      </div>
      
      {/* Button section */}
      <div className="flex items-center gap-2 sm:gap-4">
        <button
          onClick={() => openModal('create')}
          className="flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors whitespace-nowrap"
        >
          <span className="hidden xs:inline">+ New Task</span>
          <span className="xs:hidden">+ Task</span>
        </button>
        <button
          onClick={logout}
          className="flex-1 sm:flex-none px-3 sm:px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors whitespace-nowrap"
        >
          Logout
        </button>
      </div>
    </div>
  </div>
</header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Error Message for Tasks API */}
        {tasksError && (
          <div className="mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
            <div className="flex items-start">
              <svg
                className="h-5 w-5 text-red-600 dark:text-red-400 mr-2 mt-0.5"
                fill="none"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <div className="flex-1">
                <h3 className="text-sm font-medium text-red-800 dark:text-red-300">
                  Unable to load tasks
                </h3>
                <p className="text-sm text-red-700 dark:text-red-400 mt-1">
                  {tasksError}
                </p>
                <button
                  onClick={() => fetchTasks()}
                  className="mt-2 text-sm text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300 underline"
                >
                  Try again
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Loading State for Tasks */}
        {tasksLoading && !tasksError && (
          <div className="flex items-center justify-center py-12">
            <div className="text-center">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
              <p className="text-gray-600 dark:text-gray-400">Loading tasks...</p>
            </div>
          </div>
        )}

        {/* Task Board - Always render, even if empty */}
        {!tasksLoading && <TaskBoard />}
      </main>

      {/* Task Modal */}
      <TaskModal />
    </div>
  );
}