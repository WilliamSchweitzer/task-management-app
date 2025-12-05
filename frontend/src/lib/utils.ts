import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { TaskPriority, TaskStatus } from '@/types';

// Merge Tailwind classes with clsx
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Get border color class based on priority
export function getPriorityBorderColor(priority?: TaskPriority): string {
  switch (priority) {
    case 'high':
      return 'border-l-priority-high';
    case 'medium':
      return 'border-l-priority-medium';
    case 'low':
      return 'border-l-priority-low';
    default:
      return 'border-l-priority-medium';
  }
}

// Get background color class based on priority (for badges)
export function getPriorityBgColor(priority?: TaskPriority): string {
  switch (priority) {
    case 'high':
      return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300';
    case 'medium':
      return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-300';
    case 'low':
      return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300';
    default:
      return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-300';
  }
}

// Get status color class (for column headers)
export function getStatusColor(status: TaskStatus): string {
  switch (status) {
    case 'todo':
      return 'bg-indigo-500';
    case 'in-progress':
      return 'bg-amber-500';
    case 'done':
      return 'bg-emerald-500';
    default:
      return 'bg-gray-500';
  }
}

// Get status label
export function getStatusLabel(status: TaskStatus): string {
  switch (status) {
    case 'todo':
      return 'To Do';
    case 'in-progress':
      return 'In Progress';
    case 'done':
      return 'Done';
    default:
      return status;
  }
}

// Get priority label
export function getPriorityLabel(priority?: TaskPriority): string {
  if (!priority) return 'Medium';
  return priority.charAt(0).toUpperCase() + priority.slice(1);
}

// Format date for display
export function formatDate(dateString?: string): string {
  if (!dateString) return '';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  }).format(date);
}

// Format date for input
export function formatDateForInput(dateString?: string): string {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toISOString().split('T')[0];
}

// Check if date is past due
export function isPastDue(dateString?: string): boolean {
  if (!dateString) return false;
  const dueDate = new Date(dateString);
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return dueDate < today;
}

// Generate a unique ID (for optimistic updates)
export function generateTempId(): string {
  return `temp_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;
}
