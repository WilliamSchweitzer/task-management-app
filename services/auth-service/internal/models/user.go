package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Name         string    `gorm:"not null" json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ValidateEmail(email string) error {
	// Simple email validation logic (can be improved with regex)
	if len(email) < 3 || len(email) > 254 {
		return fmt.Errorf("invalid email length")
	} else if !strings.Contains(email, "@") || email == "" {
		return fmt.Errorf("invalid email format")
	} else if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
}

func ValidateName(name string) error {
	if len(name) < 2 || len(name) > 100 {
		return fmt.Errorf("name must be between 2 and 100 characters")
	} else if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

func (User) TableName() string {
	return "auth.users"
}
