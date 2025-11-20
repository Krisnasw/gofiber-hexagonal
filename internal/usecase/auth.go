package usecase

import (
	"app-hexagonal/internal/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthUsecaseInterface defines the interface for authentication use cases
// This helps with dependency inversion in our hexagonal architecture
type AuthUsecaseInterface interface {
	Login(credentials *domain.Credentials) (*domain.TokenResponse, error)
	RefreshToken(refreshToken string) (*domain.TokenResponse, error)
	Logout(accessToken string) error
	ValidateToken(tokenString string) (*domain.JWTClaims, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// AuthUsecase handles authentication business logic
type AuthUsecase struct {
	userRepo domain.UserRepository
	// No need for authRepo since we're using JWT
	jwtSecret []byte
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo domain.UserRepository) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  userRepo,
		jwtSecret: []byte("your-secret-key-change-this-in-production"), // In production, load from config
	}
}

// Login authenticates a user and generates JWT tokens
func (au *AuthUsecase) Login(credentials *domain.Credentials) (*domain.TokenResponse, error) {
	// Find user by email
	user, err := au.userRepo.FindByEmail(credentials.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify the password
	if !au.CheckPasswordHash(credentials.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate access token (1 hour expiry)
	accessToken, err := au.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (7 days expiry)
	refreshToken, err := au.generateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour in seconds
	}, nil
}

// RefreshToken generates a new access token using a refresh token
func (au *AuthUsecase) RefreshToken(refreshToken string) (*domain.TokenResponse, error) {
	// Parse and validate the refresh token
	claims, err := au.parseToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if it's actually a refresh token by checking expiry
	// Refresh tokens have longer expiry than access tokens
	// In a real implementation, you might want to store token type in claims

	// Generate new access token
	newAccessToken, err := au.generateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := au.generateRefreshToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	return &domain.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour in seconds
	}, nil
}

// Logout is a no-op with JWT since tokens are stateless
// In a real implementation, you might want to implement token blacklisting
func (au *AuthUsecase) Logout(accessToken string) error {
	// With JWT, logout is typically handled on the client side
	// by deleting the tokens. Server-side blacklisting would
	// require storing blacklisted tokens, which defeats the
	// purpose of stateless JWT.

	// For now, we'll just validate the token and return success
	_, err := au.parseToken(accessToken)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// ValidateToken checks if a token is valid
func (au *AuthUsecase) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	return au.parseToken(tokenString)
}

// generateAccessToken generates a JWT access token
func (au *AuthUsecase) generateAccessToken(userID, email string) (string, error) {
	claims := &domain.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(au.jwtSecret)
}

// generateRefreshToken generates a JWT refresh token
func (au *AuthUsecase) generateRefreshToken(userID, email string) (string, error) {
	claims := &domain.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(au.jwtSecret)
}

// parseToken parses and validates a JWT token
func (au *AuthUsecase) parseToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return au.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// HashPassword hashes a password using bcrypt
func (au *AuthUsecase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a password with its hash
func (au *AuthUsecase) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
