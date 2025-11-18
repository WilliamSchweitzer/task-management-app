package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/database"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/model"
)

// Mock database functions for testing
var (
	mockCreateTask   func(task model.Task) error
	mockGetTask      func(taskID uuid.UUID) (*model.Task, error)
	mockUpdateTask   func(taskID uuid.UUID, updates map[string]interface{}) (*model.Task, error)
	mockDeleteTask   func(taskID uuid.UUID) error
	mockCompleteTask func(taskID uuid.UUID) (*model.Task, error)
)

func setupTestRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/tasks", CreateTask)
	r.Get("/tasks/{taskID}", GetTask)
	r.Put("/tasks/{taskID}", UpdateTask)
	r.Delete("/tasks/{taskID}", DeleteTask)
	r.Patch("/tasks/{taskID}/complete", CompleteTask)
	return r
}

func TestCreateTask(t *testing.T) {
	router := setupTestRouter()

	t.Run("successful task creation", func(t *testing.T) {
		userID := uuid.New()
		reqBody := CreateTaskRequest{
			Title:  "Test Task",
			Status: "todo",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", userID.String())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response model.Task
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Task", response.Title)
		assert.Equal(t, "todo", response.Status)
		assert.Equal(t, userID, response.UserID)
	})

	t.Run("missing user ID header", func(t *testing.T) {
		reqBody := CreateTaskRequest{
			Title:  "Test Task",
			Status: "todo",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		reqBody := CreateTaskRequest{
			Title:  "Test Task",
			Status: "todo",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", "invalid-uuid")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid User ID")
	})

	t.Run("invalid request payload", func(t *testing.T) {
		userID := uuid.New()
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", userID.String())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid request payload")
	})

	t.Run("validation error - missing title", func(t *testing.T) {
		userID := uuid.New()
		reqBody := CreateTaskRequest{
			Status: "todo",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", userID.String())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "title field is required")
	})

	t.Run("validation error - invalid status", func(t *testing.T) {
		userID := uuid.New()
		reqBody := CreateTaskRequest{
			Title:  "Test Task",
			Status: "invalid-status",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", userID.String())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "status must be one of")
	})
}

func TestGetTask(t *testing.T) {
	router := setupTestRouter()

	t.Run("successful task retrieval", func(t *testing.T) {
		// This test would need a real database or more sophisticated mocking
		// For now, we'll test the error cases
	})

	t.Run("invalid task ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks/invalid-uuid", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid task ID")
	})
}

func TestUpdateTask(t *testing.T) {
	router := setupTestRouter()

	t.Run("invalid task ID", func(t *testing.T) {
		reqBody := UpdateTaskRequest{
			Title: stringPtr("Updated Title"),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/tasks/invalid-uuid", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid task ID")
	})

	t.Run("invalid request payload", func(t *testing.T) {
		taskID := uuid.New()
		req := httptest.NewRequest("PUT", "/tasks/"+taskID.String(), bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid request payload")
	})
}

func TestDeleteTask(t *testing.T) {
	router := setupTestRouter()

	t.Run("invalid task ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/tasks/invalid-uuid", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid task ID")
	})
}

func TestCompleteTask(t *testing.T) {
	router := setupTestRouter()

	t.Run("invalid task ID", func(t *testing.T) {
		req := httptest.NewRequest("PATCH", "/tasks/invalid-uuid/complete", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid task ID")
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
