package repository

import (
	"context"

	"github.com/ductong169z/blog-api/internal/auth"
	"github.com/ductong169z/blog-api/internal/models"
	"github.com/ductong169z/blog-api/pkg/cache/redis"
)

// News redis repository
type redisRepo struct {
	rdb redis.Client
}

// News redis repository constructor
func NewRedisRepo(rdb redis.Client) auth.RedisRepository {
	return &redisRepo{rdb: rdb}
}

// GetUserByIDCtx implements auth.RedisRepository.
func (r *redisRepo) GetUserByIDCtx(ctx context.Context, key string) (*models.User, error) {
	panic("unimplemented")
}
