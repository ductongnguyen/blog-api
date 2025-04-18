//go:generate mockgen -source mysql_repository.go -destination mock/mysql_repository_mock.go -package mock

package news

import (
	"context"

	"github.com/ductong169z/blog-api/internal/models"
	"github.com/ductong169z/blog-api/pkg/utils"
	"github.com/google/uuid"
)

// News Repository
type Repository interface {
	Create(ctx context.Context, news *models.News) (*models.News, error)
	Update(ctx context.Context, news *models.News) (*models.News, error)
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
	GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error)
}
