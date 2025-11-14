package model

import (
	"errors"
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
	// ------------------------------------------------------------------
	// 1. Basic length checks (RFC 5321: local-part ≤64, domain ≤255)
	// ------------------------------------------------------------------
	const (
		minLen = 3   // a@b
		maxLen = 254 // local 64 + @ + domain 255 - 1 for the dot
	)
	if len(email) < minLen || len(email) > maxLen {
		return errors.New("invalid email length")
	}

	// ------------------------------------------------------------------
	// 2. Must contain exactly one '@'
	// ------------------------------------------------------------------
	atIdx := strings.IndexByte(email, '@')
	if atIdx < 1 || atIdx == len(email)-1 || strings.Contains(email[atIdx+1:], "@") {
		return errors.New("invalid email format")
	}

	// ------------------------------------------------------------------
	// 3. Local part (before @) – only printable ASCII, no control chars
	// ------------------------------------------------------------------
	local := email[:atIdx]
	if !isLocalPartValid(local) {
		return errors.New("invalid characters in local part")
	}

	// ------------------------------------------------------------------
	// 4. Domain part (after @) – at least one dot, no leading/trailing dot
	// ------------------------------------------------------------------
	domain := email[atIdx+1:]
	if !isDomainValid(domain) {
		return errors.New("invalid domain")
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

func isDomainValid(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '.' || s[len(s)-1] == '.' {
		return false
	}
	dotSeen := false
	for _, r := range s {
		if r == '.' {
			dotSeen = true
			continue
		}
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return dotSeen
}

func isLocalPartValid(s string) bool {
	for i, r := range s {
		if r <= 0x1F || r == 0x7F { // control chars
			return false
		}
		// Allowed specials inside quotes are skipped – we keep it simple
		if r == ' ' || strings.ContainsRune("()<>[]:;@\\,\"", r) {
			return false
		}
		if i == 0 && (r == '.' || r == '@') {
			return false // cannot start with dot or @
		}
		if i == len(s)-1 && r == '.' {
			return false // cannot end with dot
		}
	}
	return len(s) > 0
}
