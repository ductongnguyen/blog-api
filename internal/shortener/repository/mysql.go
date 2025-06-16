package repository

import (
	"context"
	"errors"

	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/ductong169z/shorten-url/internal/shortener"
	"gorm.io/gorm"
)

// News Repository
type repo struct {
	db *gorm.DB
}

// News repository constructor
func NewRepository(db *gorm.DB) shortener.Repository {
	return &repo{db: db}
}

func (r *repo) CreateShortURL(ctx context.Context, url *models.ShortURL) error {
	return r.db.WithContext(ctx).Create(url).Error
}

func (r *repo) GetShortURLByCode(ctx context.Context, code string) (*models.ShortURL, error) {
	var url models.ShortURL
	if err := r.db.WithContext(ctx).Where("short_code = ?", code).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &url, nil
}

func (r *repo) IncrementClickCount(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).Model(&models.ShortURL{}).Where("short_code = ?", code).UpdateColumn("click_count", gorm.Expr("click_count + 1")).Error
}

func (r *repo) IsShortCodeExist(ctx context.Context, code string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.ShortURL{}).
		Where("short_code = ?", code).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
