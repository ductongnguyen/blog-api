package usecase

import (
	"context"

	"github.com/ductong169z/blog-api/config"
	"github.com/ductong169z/blog-api/internal/auth"
	"github.com/ductong169z/blog-api/internal/models"
	"github.com/ductong169z/blog-api/pkg/logger"
	"github.com/google/uuid"
)

const (
	basePrefix    = "api-news:"
	cacheDuration = 3600
)

// News UseCase constructor
func NewUseCase(cfg *config.Config, repo auth.Repository, redisRepo auth.RedisRepository, logger logger.Logger) auth.UseCase {
	return &usecase{cfg: cfg, repo: repo, redisRepo: redisRepo, logger: logger}
}

// useCase
type usecase struct {
	cfg       *config.Config
	repo      auth.Repository
	redisRepo auth.RedisRepository
	logger    logger.Logger
}

// GetUserByID implements auth.UseCase.
func (u *usecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	panic("unimplemented")
}

// Login implements auth.UseCase.
func (u *usecase) Login(ctx context.Context, user *models.User) (*models.User, error) {
	panic("unimplemented")
}

// Register implements auth.UseCase.
func (u *usecase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	panic("unimplemented")
}
