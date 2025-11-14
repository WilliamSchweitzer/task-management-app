package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/williamschweitzer/task-management-app/services/auth-service/internal/database"
	"github.com/williamschweitzer/task-management-app/services/auth-service/internal/model"
)

type JWTConfig struct {
	Secret              string
	Issuer              string
	AccessTokenDuration time.Duration
}

var DefaultJWTConfig = JWTConfig{
	Secret:              string(os.Getenv("JWT_SECRET")),
	Issuer:              "task-management-auth",
	AccessTokenDuration: 15 * time.Minute,
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(cfg JWTConfig, userID uuid.UUID, email string) (string, error) {
	if cfg.Secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	if err := model.ValidateEmail(email); err != nil {
		return "", err
	}

	if userID == uuid.Nil {
		return "", fmt.Errorf("userID cannot be nil")
	}

	expiryStr := os.Getenv("ACCESS_TOKEN_EXPIRY")
	if expiryStr == "" {
		expiryStr = "15m"
	}

	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		expiry = 15 * time.Minute
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "task-management-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func GenerateRefreshToken(cfg JWTConfig, userID uuid.UUID, email string) (string, time.Time, error) {
	if cfg.Secret == "" {
		return "", time.Now(), fmt.Errorf("JWT_SECRET is not set")
	}

	if err := model.ValidateEmail(email); err != nil {
		return "", time.Now(), err
	}

	if userID == uuid.Nil {
		return "", time.Now(), fmt.Errorf("userID cannot be nil")
	}

	expiryStr := os.Getenv("REFRESH_TOKEN_EXPIRY")
	if expiryStr == "" {
		expiryStr = "7d"
	}

	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		expiry = 7 * 24 * time.Hour
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "task-management-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Secret))

	return tokenString, claims.ExpiresAt.Time, err
}

func StoreRefreshToken(userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	refreshTokenEntry := model.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&refreshTokenEntry).Error; err != nil {
		return err
	}

	return nil
}

func LookupRefreshToken(tokenHash string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	err := database.DB.Where("token_hash = ?",
		tokenHash,
	).First(&refreshToken).Error

	if err != nil {
		return nil, err
	}

	return &refreshToken, err
}

func RevokeRefreshToken(refreshToken *model.RefreshToken) error {
	now := time.Now()
	refreshToken.RevokedAt = &now

	if err := database.DB.Save(refreshToken).Error; err != nil {
		return err
	}

	return nil
}

func ValidateToken(cfg JWTConfig, tokenStr string) (*Claims, error) {
	if cfg.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func HashToken(token string) (string, error) {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:]), nil
}
