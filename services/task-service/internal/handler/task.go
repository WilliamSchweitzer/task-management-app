package handler

import (
	"net/http"

	"github.com/google/uuid"
)

// CreateTask Request struct

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

	// Decode request and create request object

	// Validate access token and refresh token using VerifyToken endpoint from auth service? Or is there better way to do this?

	// Validate task input

	// Check if matching UUIDs?

	// Create task object

	// Store task in tasks.tasks service.StoreTask

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"Create task endpoint - to be implemented"}`))
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message":"Get task endpoint - to be implemented"}`))
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
