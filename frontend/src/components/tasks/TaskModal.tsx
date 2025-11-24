'use client';

import { useState, useEffect, FormEvent } from 'react';
import { Modal, Button, Input, Textarea, Select } from '@/components/ui';
import { useTaskStore } from '@/stores/task-store';
import { CreateTaskRequest, UpdateTaskRequest, TaskStatus, TaskPriority } from '@/types';
import {
  cn,
  formatDateForInput,
  formatDate,
  getPriorityBgColor,
  getStatusColor,
  getStatusLabel,
  getPriorityLabel,
} from '@/lib/utils';
import { Calendar, Clock, User, Edit2, Trash2 } from 'lucide-react';

const STATUS_OPTIONS = [
  { value: 'todo', label: 'To Do' },
  { value: 'in-progress', label: 'In Progress' },
  { value: 'done', label: 'Done' },
];

const PRIORITY_OPTIONS = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
];

export function TaskModal() {
  const {
    isModalOpen,
    modalMode,
    selectedTask,
    closeModal,
    createTask,
    updateTask,
    deleteTask,
  } = useTaskStore();

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isEditing, setIsEditing] = useState(modalMode === 'edit' || modalMode === 'create');
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    status: 'todo' as TaskStatus,
    priority: 'medium' as TaskPriority,
    due_date: '',
  });

  // Initialize form data when modal opens
  useEffect(() => {
    if (selectedTask) {
      setFormData({
        title: selectedTask.title,
        description: selectedTask.description || '',
        status: selectedTask.status,
        priority: (selectedTask.priority || 'medium') as TaskPriority,
        due_date: formatDateForInput(selectedTask.due_date),
      });
      setIsEditing(modalMode === 'edit');
    } else {
      setFormData({
        title: '',
        description: '',
        status: selectedTask?.status || 'todo',
        priority: 'medium',
        due_date: '',
      });
      setIsEditing(true);
    }
  }, [selectedTask, modalMode]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!formData.title.trim()) return;

    setIsSubmitting(true);

    try {
      if (selectedTask && selectedTask.id) {
        // Update existing task
        const updates: UpdateTaskRequest = {
          title: formData.title,
          description: formData.description || undefined,
          status: formData.status,
          priority: formData.priority,
          due_date: formData.due_date || undefined,
        };
        await updateTask(selectedTask.id, updates);
      } else {
        // Create new task
        const newTask: CreateTaskRequest = {
          title: formData.title,
          description: formData.description || undefined,
          status: formData.status,
          priority: formData.priority,
          due_date: formData.due_date || undefined,
        };
        await createTask(newTask);
      }
      closeModal();
    } catch {
      // Error is handled by the store
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedTask?.id) return;
    if (!confirm('Are you sure you want to delete this task?')) return;

    setIsSubmitting(true);
    try {
      await deleteTask(selectedTask.id);
      closeModal();
    } catch {
      // Error handled by store
    } finally {
      setIsSubmitting(false);
    }
  };

  const getModalTitle = () => {
    if (modalMode === 'create') return 'Create Task';
    if (isEditing) return 'Edit Task';
    return 'Task Details';
  };

  return (
    <Modal isOpen={isModalOpen} onClose={closeModal} title={getModalTitle()} size="lg">
      {isEditing ? (
        // Edit/Create Form
        <form onSubmit={handleSubmit} className="p-6 space-y-5">
          <Input
            label="Title"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            placeholder="Enter task title"
            required
            autoFocus
          />

          <Textarea
            label="Description"
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
            placeholder="Enter task description (optional)"
            rows={4}
          />

          <div className="grid grid-cols-2 gap-4">
            <Select
              label="Status"
              value={formData.status}
              onChange={(e) =>
                setFormData({ ...formData, status: e.target.value as TaskStatus })
              }
              options={STATUS_OPTIONS}
            />

            <Select
              label="Priority"
              value={formData.priority}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  priority: e.target.value as TaskPriority,
                })
              }
              options={PRIORITY_OPTIONS}
            />
          </div>

          <Input
            label="Due Date"
            type="date"
            value={formData.due_date}
            onChange={(e) =>
              setFormData({ ...formData, due_date: e.target.value })
            }
          />

          {/* Actions */}
          <div className="flex items-center justify-between pt-4 border-t border-gray-200 dark:border-gray-700">
            <div>
              {selectedTask?.id && (
                <Button
                  type="button"
                  variant="danger"
                  onClick={handleDelete}
                  disabled={isSubmitting}
                >
                  <Trash2 className="w-4 h-4 mr-2" />
                  Delete
                </Button>
              )}
            </div>
            <div className="flex gap-3">
              <Button type="button" variant="secondary" onClick={closeModal}>
                Cancel
              </Button>
              <Button type="submit" isLoading={isSubmitting}>
                {selectedTask?.id ? 'Save Changes' : 'Create Task'}
              </Button>
            </div>
          </div>
        </form>
      ) : (
        // View Mode
        <div className="p-6">
          {selectedTask && (
            <>
              {/* Title and Actions */}
              <div className="flex items-start justify-between mb-4">
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                  {selectedTask.title}
                </h3>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setIsEditing(true)}
                >
                  <Edit2 className="w-4 h-4 mr-2" />
                  Edit
                </Button>
              </div>

              {/* Status and Priority Badges */}
              <div className="flex flex-wrap gap-2 mb-4">
                <span
                  className={cn(
                    'inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium text-white',
                    getStatusColor(selectedTask.status)
                  )}
                >
                  <span
                    className={cn(
                      'w-2 h-2 rounded-full bg-white/30'
                    )}
                  />
                  {getStatusLabel(selectedTask.status)}
                </span>
                <span
                  className={cn(
                    'inline-flex items-center px-3 py-1 rounded-full text-sm font-medium',
                    getPriorityBgColor(selectedTask.priority)
                  )}
                >
                  {getPriorityLabel(selectedTask.priority)} Priority
                </span>
              </div>

              {/* Description */}
              {selectedTask.description && (
                <div className="mb-6">
                  <h4 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">
                    Description
                  </h4>
                  <p className="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
                    {selectedTask.description}
                  </p>
                </div>
              )}

              {/* Metadata */}
              <div className="space-y-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                {selectedTask.due_date && (
                  <div className="flex items-center gap-2 text-sm">
                    <Calendar className="w-4 h-4 text-gray-400" />
                    <span className="text-gray-500 dark:text-gray-400">Due:</span>
                    <span className="text-gray-700 dark:text-gray-300">
                      {formatDate(selectedTask.due_date)}
                    </span>
                  </div>
                )}
                {selectedTask.completed_at && (
                  <div className="flex items-center gap-2 text-sm">
                    <Clock className="w-4 h-4 text-emerald-500" />
                    <span className="text-gray-500 dark:text-gray-400">
                      Completed:
                    </span>
                    <span className="text-emerald-600 dark:text-emerald-400">
                      {formatDate(selectedTask.completed_at)}
                    </span>
                  </div>
                )}
                <div className="flex items-center gap-2 text-sm">
                  <Clock className="w-4 h-4 text-gray-400" />
                  <span className="text-gray-500 dark:text-gray-400">
                    Created:
                  </span>
                  <span className="text-gray-700 dark:text-gray-300">
                    {formatDate(selectedTask.created_at)}
                  </span>
                </div>
              </div>

              {/* Actions */}
              <div className="flex justify-between items-center pt-6 mt-6 border-t border-gray-200 dark:border-gray-700">
                <Button variant="danger" onClick={handleDelete}>
                  <Trash2 className="w-4 h-4 mr-2" />
                  Delete
                </Button>
                <Button variant="secondary" onClick={closeModal}>
                  Close
                </Button>
              </div>
            </>
          )}
        </div>
      )}
    </Modal>
  );
}
