'use client';

import { create } from 'zustand';
import { Task, TaskStatus, CreateTaskRequest, UpdateTaskRequest } from '@/types';
import { taskService } from '@/services/task-service';
import { generateTempId } from '@/lib/utils';

interface TaskState {
  tasks: Task[];
  isLoading: boolean;
  error: string | null;
  selectedTask: Task | null;
  isModalOpen: boolean;
  modalMode: 'create' | 'edit' | 'view';
  
  // Actions
  fetchTasks: () => Promise<void>;
  createTask: (task: CreateTaskRequest) => Promise<Task>;
  updateTask: (id: string, updates: UpdateTaskRequest) => Promise<Task>;
  updateTaskStatus: (id: string, status: TaskStatus) => Promise<void>;
  deleteTask: (id: string) => Promise<void>;
  completeTask: (id: string) => Promise<void>;
  
  // Modal actions
  openModal: (mode: 'create' | 'edit' | 'view', task?: Task) => void;
  closeModal: () => void;
  setSelectedTask: (task: Task | null) => void;
  
  // Error handling
  clearError: () => void;
  
  // Computed getters
  getTasksByStatus: (status: TaskStatus) => Task[];
}

export const useTaskStore = create<TaskState>((set, get) => ({
  tasks: [],
  isLoading: false,
  error: null,
  selectedTask: null,
  isModalOpen: false,
  modalMode: 'create',

  fetchTasks: async () => {
    set({ isLoading: true, error: null });
    try {
      const tasks = await taskService.listTasks();
      set({ tasks, isLoading: false });
    } catch (err) {
      const error = err as { message?: string };
      set({
        isLoading: false,
        error: error.message || 'Failed to fetch tasks',
      });
    }
  },

  createTask: async (taskData: CreateTaskRequest) => {
    set({ error: null });
    
    // Optimistic update with temp ID
    const tempTask: Task = {
      id: generateTempId(),
      user_id: '',
      title: taskData.title,
      description: taskData.description,
      status: taskData.status || 'todo',
      priority: taskData.priority || 'medium',
      due_date: taskData.due_date,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    
    set((state) => ({
      tasks: [...state.tasks, tempTask],
    }));

    try {
      const newTask = await taskService.createTask(taskData);
      // Replace temp task with real task
      set((state) => ({
        tasks: state.tasks.map((t) => (t.id === tempTask.id ? newTask : t)),
      }));
      return newTask;
    } catch (err) {
      // Rollback on failure
      set((state) => ({
        tasks: state.tasks.filter((t) => t.id !== tempTask.id),
        error: (err as { message?: string }).message || 'Failed to create task',
      }));
      throw err;
    }
  },

  updateTask: async (id: string, updates: UpdateTaskRequest) => {
    set({ error: null });
    
    // Store original for rollback
    const originalTask = get().tasks.find((t) => t.id === id);
    if (!originalTask) throw new Error('Task not found');

    // Optimistic update
    set((state) => ({
      tasks: state.tasks.map((t) =>
        t.id === id ? { ...t, ...updates, updated_at: new Date().toISOString() } : t
      ),
    }));

    try {
      const updatedTask = await taskService.updateTask(id, updates);
      set((state) => ({
        tasks: state.tasks.map((t) => (t.id === id ? updatedTask : t)),
        selectedTask: state.selectedTask?.id === id ? updatedTask : state.selectedTask,
      }));
      return updatedTask;
    } catch (err) {
      // Rollback on failure
      set((state) => ({
        tasks: state.tasks.map((t) => (t.id === id ? originalTask : t)),
        error: (err as { message?: string }).message || 'Failed to update task',
      }));
      throw err;
    }
  },

  updateTaskStatus: async (id: string, status: TaskStatus) => {
    const originalTask = get().tasks.find((t) => t.id === id);
    if (!originalTask) return;

    // Optimistic update
    set((state) => ({
      tasks: state.tasks.map((t) =>
        t.id === id
          ? {
              ...t,
              status,
              updated_at: new Date().toISOString(),
              completed_at: status === 'done' ? new Date().toISOString() : undefined,
            }
          : t
      ),
    }));

    try {
      await taskService.updateTask(id, { status });
    } catch (err) {
      // Rollback on failure
      set((state) => ({
        tasks: state.tasks.map((t) => (t.id === id ? originalTask : t)),
        error: (err as { message?: string }).message || 'Failed to update task status',
      }));
    }
  },

  deleteTask: async (id: string) => {
    set({ error: null });
    
    // Store original for rollback
    const originalTasks = get().tasks;
    
    // Optimistic update
    set((state) => ({
      tasks: state.tasks.filter((t) => t.id !== id),
    }));

    try {
      await taskService.deleteTask(id);
      // Close modal if deleting the selected task
      if (get().selectedTask?.id === id) {
        set({ selectedTask: null, isModalOpen: false });
      }
    } catch (err) {
      // Rollback on failure
      set({
        tasks: originalTasks,
        error: (err as { message?: string }).message || 'Failed to delete task',
      });
      throw err;
    }
  },

  completeTask: async (id: string) => {
    const task = get().tasks.find((t) => t.id === id);
    if (!task) return;

    // Optimistic update
    set((state) => ({
      tasks: state.tasks.map((t) =>
        t.id === id
          ? {
              ...t,
              status: 'done' as TaskStatus,
              completed_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            }
          : t
      ),
    }));

    try {
      const completedTask = await taskService.completeTask(id);
      set((state) => ({
        tasks: state.tasks.map((t) => (t.id === id ? completedTask : t)),
      }));
    } catch (err) {
      // Rollback on failure
      set((state) => ({
        tasks: state.tasks.map((t) => (t.id === id ? task : t)),
        error: (err as { message?: string }).message || 'Failed to complete task',
      }));
    }
  },

  openModal: (mode: 'create' | 'edit' | 'view', task?: Task) => {
    set({
      isModalOpen: true,
      modalMode: mode,
      selectedTask: task || null,
    });
  },

  closeModal: () => {
    set({
      isModalOpen: false,
      selectedTask: null,
    });
  },

  setSelectedTask: (task: Task | null) => {
    set({ selectedTask: task });
  },

  clearError: () => {
    set({ error: null });
  },

  getTasksByStatus: (status: TaskStatus) => {
    return get().tasks.filter((task) => task.status === status);
  },
}));
