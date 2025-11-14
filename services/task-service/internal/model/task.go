package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	Status      string     `gorm:"type:varchar(50);not null;default:'todo'" json:"status"`
	Priority    *string    `gorm:"type:varchar(50);default:'medium'" json:"priority,omitempty"`
	DueDate     *time.Time `gorm:"type:timestamp" json:"due_date,omitempty"`
	CompletedAt *time.Time `gorm:"type:timestamp" json:"completed_at,omitempty"`
	CreatedAt   time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Task) TableName() string {
	return "tasks.tasks"
}

func (t Task) Validate() error {
	if t.Title == "" {
		return errors.New("title field is required")
	}

	if len(t.Title) > 255 {
		return errors.New("title field length must be 255 characters or less")
	}

	if t.Status == "" {
		return errors.New("title field is required")
	}

	if len(t.Status) > 50 {
		return errors.New("status field length must be 50 characters or less")
	}

	if t.Priority != nil && (*t.Priority != "low" && *t.Priority != "medium" && *t.Priority != "high") {
		*t.Priority = "medium"
	}

	if t.Status != "todo" && t.Status != "in-progress" && t.Status != "done" {
		t.Status = "todo"
	}

	if t.DueDate != nil {
		if (*t.DueDate).Before(time.Now()) {
			return errors.New("due date cannot be in the past")
		}
	}

	if (t.Status == "done" && t.CompletedAt == nil) || (t.Status != "done" && t.CompletedAt != nil) {
		return errors.New("status must be done and completedAt must be set")
	}

	return nil
}
