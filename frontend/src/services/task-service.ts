import { apiClient } from '@/lib/api-client';
import { API_CONFIG } from '@/lib/config';
import { Task, CreateTaskRequest, UpdateTaskRequest } from '@/types';

export const taskService = {
  async listTasks(): Promise<Task[]> {
    return apiClient.get<Task[]>(API_CONFIG.TASKS.LIST);
  },

  async getTask(id: string): Promise<Task> {
    return apiClient.get<Task>(API_CONFIG.TASKS.GET(id));
  },

  async createTask(task: CreateTaskRequest): Promise<Task> {
    return apiClient.post<Task>(API_CONFIG.TASKS.CREATE, task);
  },

  async updateTask(id: string, updates: UpdateTaskRequest): Promise<Task> {
    return apiClient.put<Task>(API_CONFIG.TASKS.UPDATE(id), updates);
  },

  async deleteTask(id: string): Promise<void> {
    return apiClient.delete(API_CONFIG.TASKS.DELETE(id));
  },

  async completeTask(id: string): Promise<Task> {
    return apiClient.patch<Task>(API_CONFIG.TASKS.COMPLETE(id));
  },
};
