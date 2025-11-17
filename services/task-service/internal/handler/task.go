package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/database"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/model"
)

type CreateTaskRequest struct {
	// UserID      uuid.UUID  `json:"user_id"` - Kong provides UserID
	Title       string     `json:"title"`
	Descrption  *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    *string    `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
}

type GetTaskRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetTaskResponse struct {
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

func ListTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"List tasks endpoint - to be implemented"}`))
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	// Check Kong set the header (JWT Verified by Kong -> auth-service)
	userIDStr := r.Header.Get("X-User-Id")
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
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
		Description: req.Descrption,
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
	database.CreateTask(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Task created successfully!"}`))
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	// Get Tasks by UUID from DB
	var req GetTaskRequest

	// Decode Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate ID foramt
	if err := uuid.Validate(req.ID.String()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch Task from DB
	task, err := database.GetTask(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create response
	resp := GetTaskResponse{
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"Update task endpoint - to be implemented"}`))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"Delete task endpoint - to be implemented"}`))
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"Complete task endpoint - to be implemented"}`))
}
