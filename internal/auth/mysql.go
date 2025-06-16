//go:generate mockgen -source mysql.go -destination mock/mysql_repository_mock.go -package mock

package auth

import (
	"context"

	"github.com/ductong169z/shorten-url/internal/models"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, userId int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}
