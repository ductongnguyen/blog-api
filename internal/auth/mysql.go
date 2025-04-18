//go:generate mockgen -source mysql_repository.go -destination mock/mysql_repository_mock.go -package mock

package auth

import (
	"context"

	"github.com/ductong169z/blog-api/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
