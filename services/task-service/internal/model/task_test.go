package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTaskValidate(t *testing.T) {
	userID := uuid.New()
	validTitle := "Test Task"
	validStatus := "todo"
	lowPriority := "low"
	mediumPriority := "medium"
	highPriority := "high"
	invalidPriority := "urgent"
	futureDate := time.Now().Add(24 * time.Hour)
	pastDate := time.Now().Add(-24 * time.Hour)
	now := time.Now()

	tests := []struct {
		name    string
		task    Task
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid task with required fields only",
			task: Task{
				UserID: userID,
				Title:  validTitle,
				Status: validStatus,
			},
			wantErr: false,
		},
		{
			name: "valid task with all fields",
			task: Task{
				UserID:      userID,
				Title:       validTitle,
				Status:      "in-progress",
				Priority:    &mediumPriority,
				DueDate:     &futureDate,
				CompletedAt: nil,
			},
			wantErr: false,
		},
		{
			name: "valid completed task",
			task: Task{
				UserID:      userID,
				Title:       validTitle,
				Status:      "done",
				Priority:    &highPriority,
				CompletedAt: &now,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			task: Task{
				UserID: userID,
				Status: validStatus,
			},
			wantErr: true,
			errMsg:  "title field is required",
		},
		{
			name: "title too long",
			task: Task{
				UserID: userID,
				Title:  string(make([]byte, 256)),
				Status: validStatus,
			},
			wantErr: true,
			errMsg:  "title field length must be 255 characters or less",
		},
		{
			name: "missing status",
			task: Task{
				UserID: userID,
				Title:  validTitle,
				Status: "",
			},
			wantErr: true,
			errMsg:  "status field is required",
		},
		{
			name: "invalid status",
			task: Task{
				UserID: userID,
				Title:  validTitle,
				Status: "invalid-status",
			},
			wantErr: true,
			errMsg:  "status must be one of: todo, in-progress, done",
		},
		{
			name: "invalid priority",
			task: Task{
				UserID:   userID,
				Title:    validTitle,
				Status:   validStatus,
				Priority: &invalidPriority,
			},
			wantErr: true,
			errMsg:  "priority must be one of: low, medium, high",
		},
		{
			name: "valid low priority",
			task: Task{
				UserID:   userID,
				Title:    validTitle,
				Status:   validStatus,
				Priority: &lowPriority,
			},
			wantErr: false,
		},
		{
			name: "valid high priority",
			task: Task{
				UserID:   userID,
				Title:    validTitle,
				Status:   validStatus,
				Priority: &highPriority,
			},
			wantErr: false,
		},
		{
			name: "due date in past",
			task: Task{
				UserID:  userID,
				Title:   validTitle,
				Status:  validStatus,
				DueDate: &pastDate,
			},
			wantErr: true,
			errMsg:  "due date cannot be in the past",
		},
		{
			name: "due date in future",
			task: Task{
				UserID:  userID,
				Title:   validTitle,
				Status:  validStatus,
				DueDate: &futureDate,
			},
			wantErr: false,
		},
		{
			name: "completed task without completed_at",
			task: Task{
				UserID:      userID,
				Title:       validTitle,
				Status:      "done",
				CompletedAt: nil,
			},
			wantErr: true,
			errMsg:  "status must be done and completedAt must be set",
		},
		{
			name: "non-completed task with completed_at",
			task: Task{
				UserID:      userID,
				Title:       validTitle,
				Status:      "todo",
				CompletedAt: &now,
			},
			wantErr: true,
			errMsg:  "status must be done and completedAt must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Task.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestTaskTableName(t *testing.T) {
	task := Task{}
	expected := "tasks.tasks"
	if got := task.TableName(); got != expected {
		t.Errorf("Task.TableName() = %v, want %v", got, expected)
	}
}

func TestTaskValidateDoesNotMutate(t *testing.T) {
	invalidPriority := "urgent"
	invalidStatus := "invalid"
	task := Task{
		UserID:   uuid.New(),
		Title:    "Test",
		Status:   invalidStatus,
		Priority: &invalidPriority,
	}

	originalStatus := task.Status
	originalPriority := *task.Priority

	_ = task.Validate()

	// Verify that validation didn't mutate the task
	if task.Status != originalStatus {
		t.Errorf("Task.Validate() mutated Status: got %v, want %v", task.Status, originalStatus)
	}
	if *task.Priority != originalPriority {
		t.Errorf("Task.Validate() mutated Priority: got %v, want %v", *task.Priority, originalPriority)
	}
}
