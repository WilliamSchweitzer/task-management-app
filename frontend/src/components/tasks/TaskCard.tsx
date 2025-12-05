'use client';

import { Task, TaskStatus } from '@/types';
import {
  cn,
  getPriorityBorderColor,
  getPriorityBgColor,
  getPriorityLabel,
  formatDate,
  isPastDue,
} from '@/lib/utils';
import { useTaskStore } from '@/stores/task-store';
import { Calendar, Clock, CheckCircle2, MoreVertical, Trash2, Edit2 } from 'lucide-react';
import { useState } from 'react';

interface TaskCardProps {
  task: Task;
  isDragging?: boolean;
}

export function TaskCard({ task, isDragging }: TaskCardProps) {
  const { openModal, deleteTask, updateTaskStatus } = useTaskStore();
  const [showMenu, setShowMenu] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ top: 0, left: 0 });

  const handleCardClick = (e: React.MouseEvent) => {
    if (showMenu) return; // Don't open modal if menu is showing
    openModal('view', task);
  };

  const handleStatusChange = (newStatus: TaskStatus) => {
    updateTaskStatus(task.id, newStatus);
    setShowMenu(false);
  };

  const handleDelete = () => {
    if (confirm('Are you sure you want to delete this task?')) {
      deleteTask(task.id);
    }
    setShowMenu(false);
  };

  const handleEdit = () => {
    openModal('edit', task);
    setShowMenu(false);
  };

  const pastDue = isPastDue(task.due_date) && task.status !== 'done';

  return (
    <>
      <div
        onClick={handleCardClick}
        className={cn(
          'group relative bg-white dark:bg-gray-800 rounded-lg border-l-4 shadow-card cursor-pointer',
          'transition-all duration-200 hover:shadow-card-hover hover:-translate-y-0.5',
          'border border-gray-200 dark:border-gray-700',
          getPriorityBorderColor(task.priority),
          isDragging && 'opacity-50 rotate-2 scale-105'
        )}
      >
        {/* Card Content */}
        <div className="p-4">
          {/* Header with title and menu */}
          <div className="flex items-start justify-between gap-2 mb-2">
            <h3 className="font-medium text-gray-900 dark:text-white line-clamp-2">
              {task.title}
            </h3>
            <button
              onClick={(e) => {
                e.stopPropagation();
                const rect = e.currentTarget.getBoundingClientRect();
                setMenuPosition({
                  top: rect.bottom + window.scrollY + 4,
                  left: rect.right - 192 + window.scrollX,
                });
                setShowMenu(!showMenu);
              }}
              className={cn(
                "p-1 rounded-md text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-opacity",
                showMenu ? "opacity-100" : "opacity-0 group-hover:opacity-100"
              )}
            >
              <MoreVertical className="w-4 h-4" />
            </button>
          </div>

          {/* Description */}
          {task.description && (
            <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-2 mb-3">
              {task.description}
            </p>
          )}

          {/* Footer with metadata */}
          <div className="flex flex-wrap items-center gap-2">
            {/* Priority Badge */}
            <span
              className={cn(
                'inline-flex items-center px-2 py-0.5 rounded text-xs font-medium',
                getPriorityBgColor(task.priority)
              )}
            >
              {getPriorityLabel(task.priority)}
            </span>

            {/* Due Date */}
            {task.due_date && (
              <span
                className={cn(
                  'inline-flex items-center gap-1 text-xs',
                  pastDue
                    ? 'text-red-600 dark:text-red-400'
                    : 'text-gray-500 dark:text-gray-400'
                )}
              >
                <Calendar className="w-3 h-3" />
                {formatDate(task.due_date)}
              </span>
            )}

            {/* Completed indicator */}
            {task.completed_at && (
              <span className="inline-flex items-center gap-1 text-xs text-emerald-600 dark:text-emerald-400">
                <CheckCircle2 className="w-3 h-3" />
                Completed
              </span>
            )}
          </div>
        </div>

        {/* Quick complete button for non-done tasks */}
        {task.status !== 'done' && (
          <button
            onClick={(e) => {
              e.stopPropagation();
              handleStatusChange('done');
            }}
            className="absolute bottom-2 right-2 p-1.5 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400 opacity-0 group-hover:opacity-100 transition-opacity hover:bg-emerald-200 dark:hover:bg-emerald-900/50"
            title="Mark as complete"
          >
            <CheckCircle2 className="w-4 h-4" />
          </button>
        )}
      </div>

      {/* Dropdown Menu - Rendered outside card */}
      {showMenu && (
        <>
          {/* Backdrop */}
          <div 
            className="fixed inset-0 z-40" 
            onClick={(e) => {
              e.stopPropagation();
              setShowMenu(false);
            }}
          />
          
          {/* Menu */}
          <div 
            className="fixed w-48 bg-white dark:bg-gray-800 rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 py-1 z-50"
            style={{
              top: `${menuPosition.top}px`,
              left: `${menuPosition.left}px`,
            }}
          >
            <button
              onClick={handleEdit}
              className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center gap-2"
            >
              <Edit2 className="w-4 h-4" />
              Edit Task
            </button>

            <div className="border-t border-gray-200 dark:border-gray-700 my-1" />
            <div className="px-4 py-1">
              <span className="text-xs font-medium text-gray-500 dark:text-gray-400">
                Move to
              </span>
            </div>
            {task.status !== 'todo' && (
              <button
                onClick={() => handleStatusChange('todo')}
                className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                To Do
              </button>
            )}
            {task.status !== 'in-progress' && (
              <button
                onClick={() => handleStatusChange('in-progress')}
                className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                In Progress
              </button>
            )}
            {task.status !== 'done' && (
              <button
                onClick={() => handleStatusChange('done')}
                className="w-full px-4 py-2 text-left text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                Done
              </button>
            )}

            <div className="border-t border-gray-200 dark:border-gray-700 my-1" />
            <button
              onClick={handleDelete}
              className="w-full px-4 py-2 text-left text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 flex items-center gap-2"
            >
              <Trash2 className="w-4 h-4" />
              Delete Task
            </button>
          </div>
        </>
      )}
    </>
  );
}