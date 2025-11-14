package model

import (
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
