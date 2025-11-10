package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func (RefreshToken) TableName() string {
	return "auth.refresh_tokens"
}

func (r *RefreshToken) Revoke() {
	now := time.Now()
	r.RevokedAt = &now
}

func (r *RefreshToken) IsRevoked() bool {
	return r.RevokedAt != nil
}

func (r *RefreshToken) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}

func (r *RefreshToken) IsValid() bool {
	return !r.IsRevoked() && !r.IsExpired()
}
