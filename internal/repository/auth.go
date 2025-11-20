package repository

import (
	"gorm.io/gorm"
)

// AuthGorm represents the GORM implementation for auth-related operations
// Since we're using JWT, we don't need to store tokens in the database
type AuthGorm struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository
// Since we're using JWT, this is mostly a placeholder
func NewAuthRepository(db *gorm.DB) *AuthGorm {
	return &AuthGorm{db: db}
}

// Migrate runs any necessary migrations
// Since we're using JWT, we don't need auth token tables
func (a *AuthGorm) Migrate() error {
	// No migrations needed for JWT-based auth
	return nil
}
