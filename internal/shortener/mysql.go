//go:generate mockgen -source mysql.go -destination mock/mysql_repository_mock.go -package mock
package shortener

import (
	"context"

	"github.com/ductong169z/shorten-url/internal/models"
)

type Repository interface {
	CreateShortURL(ctx context.Context, url *models.ShortURL) error
	GetShortURLByCode(ctx context.Context, code string) (*models.ShortURL, error)
	IncrementClickCount(ctx context.Context, code string) error
	IsShortCodeExist(ctx context.Context, code string) (bool, error)
}
