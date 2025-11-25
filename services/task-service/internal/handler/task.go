package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/database"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/model"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/utils"
)

// TODO: Add authorization for task endpoints based on UserID so each user has their own tasks

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    *string    `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type GetTaskResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	Priority    *string    `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if userID.String() == "" {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var req CreateTaskRequest

	// Decode request and create request object
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create Task object
	task := model.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CompletedAt: req.CompletedAt,
	}

	// Validate task input
	if err := task.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Store task in tasks.tasks service.StoreTask
	if err := database.CreateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Task created successfully!"}`))
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := r.Context().Value(utils.UserIDKey).(uuid.UUID)

	// Fetch Tasks for THIS USER from DB
	tasks, err := database.GetTasksByUserID(userID)
	if err != nil {
		if err.Error() == "no tasks found for user" {
			// Return empty array for no tasks (not an error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]GetTaskResponse{})
			return
		}
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	// Create response
	resp := make([]GetTaskResponse, len(tasks))
	for i, task := range tasks {
		resp[i] = GetTaskResponse{
			ID:          task.ID,
			UserID:      task.UserID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			DueDate:     task.DueDate,
			CompletedAt: task.CompletedAt,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	// Return Task data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	// Get task ID from URL path
	taskIDStr := chi.URLParam(r, "taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Fetch Task from DB
	task, err := database.GetTask(taskID)
	if err != nil {
		if err.Error() == "task not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	// Create response
	resp := GetTaskResponse{
		ID:          task.ID, // Don't forget this!
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		CompletedAt: task.CompletedAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	// Return Task data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var req UpdateTaskRequest

	// Decode Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get task ID from URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Build updates map with only provided fields
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.DueDate != nil {
		updates["due_date"] = *req.DueDate
	}
	if req.CompletedAt != nil {
		updates["completed_at"] = *req.CompletedAt
	}

	// Update in database
	task, err := database.UpdateTask(taskID, updates)
	if err != nil {
		if err.Error() == "failed to update task: task not found" || strings.Contains(err.Error(), "task not found") {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Return updated task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get task ID from URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Delete from database
	err = database.DeleteTask(taskID)
	if err != nil {
		if strings.Contains(err.Error(), "task not found") {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// Return 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	// Get task ID from URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Mark task as completed
	task, err := database.CompleteTask(taskID)
	if err != nil {
		if strings.Contains(err.Error(), "task not found") {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to complete task", http.StatusInternalServerError)
		return
	}

	// Return completed task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
