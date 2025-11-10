package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/williamschweitzer/task-management-app/services/auth-service/internal/database"
	"github.com/williamschweitzer/task-management-app/services/auth-service/internal/models"
	"github.com/williamschweitzer/task-management-app/services/auth-service/internal/service"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	User         models.User `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

var cfg service.JWTConfig = service.DefaultJWTConfig

func Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password, and name are required", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := database.DB.Where("email = ?", strings.ToLower(req.Email)).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user
	user := models.User{
		Email:        strings.ToLower(req.Email),
		PasswordHash: hashedPassword,
		Name:         req.Name,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate tokens
	accessToken, err := service.GenerateAccessToken(cfg, user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshTokenExpiry, err := service.GenerateRefreshToken(cfg, user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	resp := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes
		User:         user,
	}

	// Hash refresh token
	hashedRefreshToken, err := service.HashToken(refreshToken)
	if err != nil {
		http.Error(w, "Failed to hash refresh token", http.StatusInternalServerError)
		return
	}

	// Store the refresh token to auth.refresh_tokens

	err = service.StoreRefreshToken(user.ID, hashedRefreshToken, refreshTokenExpiry)

	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Find user
	var user models.User
	result := database.DB.Where("email = ?", strings.ToLower(req.Email)).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if !service.CheckPassword(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate tokens
	accessToken, err := service.GenerateAccessToken(cfg, user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshTokenExpiry, err := service.GenerateRefreshToken(cfg, user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes
		User:         user,
	}

	// Hash refresh token
	hashedRefreshToken, err := service.HashToken(refreshToken)
	if err != nil {
		http.Error(w, "Failed to hash refresh token", http.StatusInternalServerError)
		return
	}

	// Store the refresh token to auth.refresh_tokens
	err = service.StoreRefreshToken(user.ID, hashedRefreshToken, refreshTokenExpiry)

	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Accept refresh token from request body
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := models.ValidateEmail(req.Email); err != nil {
		http.Error(w, "Email is invalid", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Hash refresh token
	hashedRefreshToken, err := service.HashToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to hash refresh token", http.StatusInternalServerError)
		return
	}

	var refreshToken *models.RefreshToken
	// Lookup refresh token in database
	refreshToken, err = service.LookupRefreshToken(hashedRefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusInternalServerError)
		return
	}

	if refreshToken.IsExpired() {
		http.Error(w, "Refresh token has expired", http.StatusUnauthorized)
		return
	}

	if refreshToken.IsRevoked() {
		http.Error(w, "Refresh token is already revoked", http.StatusUnauthorized)
		return
	}

	// Generate new tokens
	accessToken, err := service.GenerateAccessToken(cfg, refreshToken.UserID, req.Email)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, newRefreshTokenExpiry, err := service.GenerateRefreshToken(cfg, refreshToken.UserID, req.Email)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Revoke old refresh token
	if err := service.RevokeRefreshToken(refreshToken); err != nil {
		http.Error(w, "Failed to revoke old refresh token", http.StatusInternalServerError)
		return
	}

	// Hash new refresh token
	hashedNewRefreshToken, err := service.HashToken(newRefreshToken)
	if err != nil {
		http.Error(w, "Failed to hash new refresh token", http.StatusInternalServerError)
		return
	}

	// Store new refresh token
	err = service.StoreRefreshToken(refreshToken.UserID, hashedNewRefreshToken, newRefreshTokenExpiry)
	if err != nil {
		http.Error(w, "Failed to store new refresh token", http.StatusInternalServerError)
		return
	}

	resp := RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(resp)
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	// Validate token
	claims, err := service.ValidateToken(cfg, tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Return user info
	response := map[string]interface{}{
		"valid":   true,
		"user_id": claims.UserID,
		"email":   claims.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Hash refresh token
	hashedRefreshToken, err := service.HashToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to hash refresh token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := service.LookupRefreshToken(hashedRefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	if refreshToken.RevokedAt != nil {
		http.Error(w, "Already logged out", http.StatusUnauthorized)
		return
	}

	service.RevokeRefreshToken(refreshToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
