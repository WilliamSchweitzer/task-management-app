'use client';

import { TaskStatus, Task } from '@/types';
import { cn, getStatusColor, getStatusLabel } from '@/lib/utils';
import { TaskCard } from './TaskCard';
import { useDroppable } from '@dnd-kit/core';
import {
  SortableContext,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { Plus } from 'lucide-react';
import { useTaskStore } from '@/stores/task-store';

interface TaskColumnProps {
  status: TaskStatus;
  tasks: Task[];
}

export function TaskColumn({ status, tasks }: TaskColumnProps) {
  const { openModal } = useTaskStore();
  const { setNodeRef, isOver } = useDroppable({
    id: status,
  });

  const handleAddTask = () => {
    // Open modal with default status set to this column
    openModal('create', {
      id: '',
      user_id: '',
      title: '',
      status: status,
      created_at: '',
      updated_at: '',
    });
  };

  return (
    <div
      className={cn(
        'flex flex-col min-h-[calc(100vh-220px)] bg-gray-100/50 dark:bg-gray-800/30 rounded-xl',
        'border border-gray-200 dark:border-gray-700/50',
        isOver && 'ring-2 ring-brand-500 ring-inset'
      )}
    >
      {/* Column Header */}
      <div className="sticky top-0 z-10 p-4 bg-gray-100/80 dark:bg-gray-800/80 backdrop-blur-sm rounded-t-xl border-b border-gray-200 dark:border-gray-700/50">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className={cn('w-3 h-3 rounded-full', getStatusColor(status))} />
            <h2 className="font-semibold text-gray-900 dark:text-white">
              {getStatusLabel(status)}
            </h2>
            <span className="inline-flex items-center justify-center w-6 h-6 text-xs font-medium rounded-full bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
              {tasks.length}
            </span>
          </div>
          <button
            onClick={handleAddTask}
            className="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
            title="Add task"
          >
            <Plus className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* Tasks Container */}
      <div
        ref={setNodeRef}
        className={cn(
          'flex-1 p-3 space-y-3 overflow-y-auto',
          isOver && 'bg-brand-50 dark:bg-brand-900/10'
        )}
      >
        <SortableContext
          items={tasks.map((t) => t.id)}
          strategy={verticalListSortingStrategy}
        >
          {tasks.map((task) => (
            <TaskCard key={task.id} task={task} />
          ))}
        </SortableContext>

        {/* Empty state */}
        {tasks.length === 0 && (
          <div className="flex flex-col items-center justify-center h-32 text-center">
            <p className="text-sm text-gray-500 dark:text-gray-400 mb-2">
              No tasks yet
            </p>
            <button
              onClick={handleAddTask}
              className="text-sm text-brand-600 hover:text-brand-700 dark:text-brand-400 dark:hover:text-brand-300 font-medium"
            >
              Add a task
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
