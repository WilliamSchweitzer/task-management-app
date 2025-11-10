// service/jwt_test.go
package service

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken(t *testing.T) {
	testCfg := JWTConfig{
		Secret:              "test-secret-32-byte-key-for-hs256!!",
		Issuer:              "task-management-auth",
		AccessTokenDuration: 15 * time.Minute,
	}

	tests := []struct {
		name    string // <-- explicit name
		userID  uuid.UUID
		email   string
		wantErr bool
	}{
		{name: "valid_user_and_email", userID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), email: "valid@email.com", wantErr: false},
		{name: "nil_user_id", userID: uuid.Nil, email: "valid@email.com", wantErr: true},
		{name: "invalid_email_format", userID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), email: "invalid.com", wantErr: true},
		{name: "empty_email", userID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), email: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateAccessToken(testCfg, tt.userID, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	testCfg := JWTConfig{
		Secret:              "test-secret-32-byte-key-for-hs256!!",
		Issuer:              "task-management-auth",
		AccessTokenDuration: 15 * time.Minute,
	}

	tests := []struct {
		name    string // <-- explicit name
		userID  uuid.UUID
		email   string
		wantErr bool
	}{
		{name: "valid_user_and_email", userID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), email: "valid@email.com", wantErr: false},
		{name: "nil_user_id", userID: uuid.Nil, email: "valid@email.com", wantErr: true},
		{name: "invalid_email_format", userID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), email: "invalid.com", wantErr: true},
		{name: "empty_email", userID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), email: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _, err := GenerateRefreshToken(testCfg, tt.userID, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestStoreRefreshToken(t *testing.T) {
	// This is a placeholder for testing the StoreRefreshToken function.
	// Implementing this test would require a mock or in-memory Redis client.
	t.Skip("StoreRefreshToken test not implemented")
}

func TestLookupRefreshToken(t *testing.T) {
	// This is a placeholder for testing the LookupRefreshToken function.
	// Implementing this test would require a mock or in-memory Redis client.
	t.Skip("LookupRefreshToken test not implemented")
}

func TestRevokeRefreshToken(t *testing.T) {
	// This is a placeholder for testing the RevokeRefreshToken function.
	// Implementing this test would require a mock or in-memory Redis client.
	t.Skip("RevokeRefreshToken test not implemented")
}

func TestValidateToken(t *testing.T) {
	cfg := JWTConfig{
		Secret:              "test-secret-32-byte-key-for-hs256!!",
		Issuer:              "task-management-auth",
		AccessTokenDuration: 15 * time.Minute,
	}

	// ---- a fresh, valid token -------------------------------------------------
	validToken, err := GenerateAccessToken(cfg,
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		"valid@email.com")
	assert.NoError(t, err)

	// ---- an expired token (signed with the *same* secret) --------------------
	expiredClaims := &Claims{
		UserID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Email:  "valid@email.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    cfg.Issuer,
		},
	}
	expiredObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, _ := expiredObj.SignedString(cfg.Secret)

	// ---- a token with a tampered signature -----------------------------------
	tamperedToken := validToken[:len(validToken)-1] + "X"

	tests := []struct {
		name    string // <-- explicit name
		token   string
		wantErr bool
	}{
		{name: "valid_token", token: validToken, wantErr: false},
		{name: "expired_token", token: expiredToken, wantErr: true},
		{name: "malformed_token", token: "invalid.token.string", wantErr: true},
		{name: "empty_token", token: "", wantErr: true},
		{name: "tampered_signature", token: tamperedToken, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(cfg, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "valid@email.com", claims.Email)
			}
		})
	}
}

func TestHashToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // pre-computed SHA-256 hex
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "simple ascii",
			input:    "hello",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:     "with spaces and punctuation",
			input:    "my secret token!",
			expected: "c623603b5ba82c4f90420f7ab34b4264f3d095851981c5c01f770e4ad1e7a999",
		},
		{
			name:     "unicode characters",
			input:    "café",
			expected: "850f7dc43910ff890f8879c0ed26fe697c93a067ad93a7d50f466a7028a9bf4e",
		},
		{
			name:     "long input (64+ bytes)",
			input:    string(make([]byte, 100)), // 100 zero bytes
			expected: "cd00e292c5970d3c5e2f0ffa5171e555bc46bfc4faddfb4a418b6840b86e79a3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashToken(tt.input)

			// 1. Error should always be nil
			require.NoError(t, err, "HashToken should never return error")

			// 2. Output must match expected SHA-256 hex
			assert.Equal(t, tt.expected, got, "hash mismatch")

			// 3. Output must be valid hex (extra safety)
			assert.Equal(t, 64, len(got), "SHA-256 hex must be 64 chars")
			_, hexErr := hex.DecodeString(got)
			assert.NoError(t, hexErr, "output must be valid hex")
		})
	}
}
