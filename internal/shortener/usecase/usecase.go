package usecase

import (
	"context"
	"log"
	"time"

	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/ductong169z/shorten-url/internal/shortener"
	"github.com/ductong169z/shorten-url/pkg/logger"
)

type usecase struct {
	cfg    *config.Config
	repo   shortener.Repository
	cache  shortener.Cache
	logger logger.Logger
}

const DefaultCacheTTL = 1 * time.Hour

// News UseCase constructor
func NewUseCase(cfg *config.Config, repo shortener.Repository, cache shortener.Cache, logger logger.Logger) shortener.UseCase {
	return &usecase{cfg: cfg, repo: repo, cache: cache, logger: logger}
}

func (u *usecase) ShortenURL(ctx context.Context, shortURL *models.ShortURL) (*models.ShortURL, error) {

	if shortURL.ShortCode == "" {
		shortURL.ShortCode = generateShortCode(8)
	} else {
		exist, err := u.repo.IsShortCodeExist(ctx, shortURL.ShortCode)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, shortener.ErrShortCodeAlreadyExists
		}
	}

	now := time.Now()
	shortURL.CreatedAt = now

	// Set expiration time
	if u.cfg.Server.ShortURLExpiredAt > 0 {
		expiredAt := now.AddDate(0, 0, u.cfg.Server.ShortURLExpiredAt)
		shortURL.ExpiredAt = &expiredAt
	}

	shortURL.ClickCount = 0

	// Store to DB
	if err := u.repo.CreateShortURL(ctx, shortURL); err != nil {
		return nil, err
	}

	// Store to cache
	if err := u.cache.SetShortURLByCode(ctx, shortURL.ShortCode, shortURL, DefaultCacheTTL); err != nil {
		u.logger.Errorf(ctx, "Failed to set short URL %s in cache: %v", shortURL.ShortCode, err)
	}

	return shortURL, nil
}

func (u *usecase) ResolveShortCode(ctx context.Context, code string) (*models.ShortURL, error) {
	// Try cache first
	url, err := u.cache.GetShortURLByCode(ctx, code)
	if err != nil {
		// Log cache failure but continue to DB
		log.Printf("cache error: %v", err)
	}
	if url != nil {
		// Optionally increment click count asynchronously
		go u.repo.IncrementClickCount(context.Background(), code)
		url.ClickCount++
		return url, nil
	}

	// Fallback to DB
	url, err = u.repo.GetShortURLByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, shortener.ErrShortCodeNotFound
	}

	// Check expiration
	if url.ExpiredAt != nil && url.ExpiredAt.Before(time.Now()) {
		return nil, shortener.ErrShortCodeExpired
	}

	// Save to cache
	_ = u.cache.SetShortURLByCode(ctx, code, url, DefaultCacheTTL)

	// Async increment
	go u.repo.IncrementClickCount(context.Background(), code)
	url.ClickCount++

	return url, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
