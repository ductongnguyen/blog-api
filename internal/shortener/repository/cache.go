package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/ductong169z/shorten-url/internal/shortener"
	"github.com/ductong169z/shorten-url/pkg/cache/redis"
)

// News redis repository
type redisRepo struct {
	rdb redis.Client
}

// News redis repository constructor
func NewRedisRepo(rdb redis.Client) shortener.Cache {
	return &redisRepo{rdb: rdb}
}

func (r *redisRepo) GetShortURLByCode(ctx context.Context, code string) (*models.ShortURL, error) {
	data, err := r.rdb.Get(ctx, code)
	if err != nil {
		return nil, err
	}
	var url models.ShortURL
	err = json.Unmarshal([]byte(data), &url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *redisRepo) SetShortURLByCode(ctx context.Context, code string, url *models.ShortURL, ttl time.Duration) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}
	err = r.rdb.Set(ctx, code, string(data), ttl)
	if err != nil {
		return err
	}
	return nil
}
