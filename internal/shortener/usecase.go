//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package shortener

import (
	"context"

	"github.com/ductong169z/shorten-url/internal/models"
)

type UseCase interface {
	ShortenURL(ctx context.Context, shortURL *models.ShortURL) (*models.ShortURL, error)
	ResolveShortCode(ctx context.Context, code string) (*models.ShortURL, error)
}
