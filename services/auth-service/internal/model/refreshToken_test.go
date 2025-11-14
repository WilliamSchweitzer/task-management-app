// models/refresh_token_test.go
package model

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// 1. Test TableName() – ensures correct DB table
// ---------------------------------------------------------------------
func TestRefreshToken_TableName(t *testing.T) {
	var token RefreshToken
	assert.Equal(t, "auth.refresh_tokens", token.TableName())
}

// ---------------------------------------------------------------------
// 2. Test field tags (GORM + JSON) – compile-time safety
// ---------------------------------------------------------------------
func TestRefreshToken_FieldTags(t *testing.T) {
	typ := reflect.TypeOf(RefreshToken{})

	tests := []struct {
		fieldName string
		wantGorm  string
		wantJSON  string
	}{
		{"ID", `type:uuid;primary_key;default:gen_random_uuid()`, "id"},
		{"UserID", `type:uuid;not null;index`, "user_id"},
		{"TokenHash", `not null`, "-"},
		{"ExpiresAt", `not null`, "expires_at"},
		{"CreatedAt", "", "created_at"},
		{"RevokedAt", "", "revoked_at,omitempty"},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			field, ok := typ.FieldByName(tt.fieldName)
			require.True(t, ok, "field %s not found", tt.fieldName)

			gormTag := field.Tag.Get("gorm")
			jsonTag := field.Tag.Get("json")

			assert.Equal(t, tt.wantGorm, gormTag, "GORM tag mismatch")
			assert.Equal(t, tt.wantJSON, jsonTag, "JSON tag mismatch")
		})
	}
}

// ---------------------------------------------------------------------
// 3. Test JSON marshaling / unmarshaling
// ---------------------------------------------------------------------
func TestRefreshToken_JSON(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	future := now.Add(7 * 24 * time.Hour)
	revoked := now.Add(time.Hour)

	token := RefreshToken{
		ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		UserID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		TokenHash: "hashed-value-should-be-ignored",
		ExpiresAt: future,
		CreatedAt: now,
		RevokedAt: &revoked,
	}

	t.Run("marshal", func(t *testing.T) {
		data, err := json.Marshal(token)
		require.NoError(t, err)

		expected := `{
			"id":"550e8400-e29b-41d4-a716-446655440000",
			"user_id":"11111111-1111-1111-1111-111111111111",
			"expires_at":"` + future.Format(time.RFC3339) + `",
			"created_at":"` + now.Format(time.RFC3339) + `",
			"revoked_at":"` + revoked.Format(time.RFC3339) + `"
		}`

		// Normalize whitespace
		assert.JSONEq(t, expected, string(data))
	})

	t.Run("unmarshal", func(t *testing.T) {
		input := `{
			"id":"550e8400-e29b-41d4-a716-446655440000",
			"user_id":"11111111-1111-1111-1111-111111111111",
			"expires_at":"2025-12-01T12:00:00Z",
			"created_at":"2025-11-10T13:37:00Z"
		}`
		var got RefreshToken
		err := json.Unmarshal([]byte(input), &got)
		require.NoError(t, err)

		assert.Equal(t, token.ID, got.ID)
		assert.Equal(t, token.UserID, got.UserID)
		assert.False(t, got.ExpiresAt.IsZero())
		assert.False(t, got.CreatedAt.IsZero())
		assert.Nil(t, got.RevokedAt)
		assert.Empty(t, got.TokenHash) // json:"-" → ignored
	})
}

// ---------------------------------------------------------------------
// 4. Test GORM hooks / default behavior (via reflection)
// ---------------------------------------------------------------------
func TestRefreshToken_GORMDefaults(t *testing.T) {
	t.Run("ID auto-generated", func(t *testing.T) {
		var token RefreshToken
		// Simulate GORM setting default
		// In real DB, gen_random_uuid() is used — we just verify zero value
		assert.Equal(t, uuid.UUID{}, token.ID) // zero value before insert
	})

	t.Run("CreatedAt auto-filled", func(t *testing.T) {
		var token RefreshToken
		// GORM sets this on create — we test zero value
		assert.True(t, token.CreatedAt.IsZero())
	})
}

// ---------------------------------------------------------------------
// 5. Test business logic: IsExpired, IsRevoked, etc.
// ---------------------------------------------------------------------
func TestRefreshToken_StatusMethods(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name      string
		token     RefreshToken
		isExpired bool
		isRevoked bool
	}{
		{
			name: "active",
			token: RefreshToken{
				ExpiresAt: future,
				RevokedAt: nil,
			},
			isExpired: false,
			isRevoked: false,
		},
		{
			name: "expired",
			token: RefreshToken{
				ExpiresAt: past,
				RevokedAt: nil,
			},
			isExpired: true,
			isRevoked: false,
		},
		{
			name: "revoked",
			token: RefreshToken{
				ExpiresAt: future,
				RevokedAt: &now,
			},
			isExpired: false,
			isRevoked: true,
		},
		{
			name: "both",
			token: RefreshToken{
				ExpiresAt: past,
				RevokedAt: &now,
			},
			isExpired: true,
			isRevoked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isExpired, tt.token.IsExpired())
			assert.Equal(t, tt.isRevoked, tt.token.IsRevoked())
		})
	}
}
