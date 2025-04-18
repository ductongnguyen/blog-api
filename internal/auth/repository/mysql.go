package repository

import (
	"context"

	"github.com/ductong169z/blog-api/internal/auth"
	"github.com/ductong169z/blog-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// News Repository
type repo struct {
	db *gorm.DB
}

// News repository constructor
func NewRepository(db *gorm.DB) auth.Repository {
	return &repo{db: db}
}

// GetUserByID implements auth.Repository.
func (r *repo) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	panic("unimplemented")
}

// Login implements auth.Repository.
func (r *repo) Login(ctx context.Context, user *models.User) (*models.User, error) {
	panic("unimplemented")
}

// Register implements auth.Repository.
func (r *repo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	panic("unimplemented")
}
