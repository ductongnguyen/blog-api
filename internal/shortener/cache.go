//go:generate mockgen -source cache.go -destination mock/cache_mock.go -package mock
package shortener

import (
	"context"
	"time"

	"github.com/ductong169z/shorten-url/internal/models"
)

type Cache interface {
	GetShortURLByCode(ctx context.Context, code string) (*models.ShortURL, error)
	SetShortURLByCode(ctx context.Context, code string, url *models.ShortURL, ttl time.Duration) error
}
