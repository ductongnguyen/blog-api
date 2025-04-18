package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ductong169z/blog-api/config"
	"github.com/ductong169z/blog-api/internal/models"
	"github.com/ductong169z/blog-api/internal/news"
	"github.com/ductong169z/blog-api/pkg/errors"
	"github.com/ductong169z/blog-api/pkg/logger"
	"github.com/ductong169z/blog-api/pkg/utils"
	"github.com/google/uuid"
)

const (
	basePrefix    = "api-news:"
	cacheDuration = 3600
)

// News UseCase
type newsUC struct {
	cfg       *config.Config
	newsRepo  news.Repository
	redisRepo news.RedisRepository
	logger    logger.Logger
}

// News UseCase constructor
func NewNewsUseCase(cfg *config.Config, newsRepo news.Repository, redisRepo news.RedisRepository, logger logger.Logger) news.UseCase {
	return &newsUC{cfg: cfg, newsRepo: newsRepo, redisRepo: redisRepo, logger: logger}
}

// Create news
func (u *newsUC) Create(ctx context.Context, news *models.News) (*models.News, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, errors.NewUnauthorizedError(errors.WithMessage(err, "newsUC.Create.GetUserFromCtx"))
	}

	news.AuthorID = user.UserID
	if err = utils.ValidateStruct(ctx, news); err != nil {
		return nil, errors.NewBadRequestError(errors.WithMessage(err, "newsUC.Create.ValidateStruct"))
	}

	n, err := u.newsRepo.Create(ctx, news)
	if err != nil {
		return nil, err
	}

	return n, err
}

// Update news item
func (u *newsUC) Update(ctx context.Context, news *models.News) (*models.News, error) {
	newsByID, err := u.newsRepo.GetNewsByID(ctx, news.NewsID)
	if err != nil {
		return nil, err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String(), u.logger); err != nil {
		return nil, errors.NewError(http.StatusForbidden, errors.ErrForbidden, errors.WithMessage(err, "newsUC.Update.ValidateIsOwner"))
	}

	updatedUser, err := u.newsRepo.Update(ctx, news)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.DeleteNewsCtx(ctx, u.getKeyWithPrefix(news.NewsID.String())); err != nil {
		u.logger.Errorf(ctx, "newsUC.Update.DeleteNewsCtx: %v", err)
	}

	return updatedUser, nil
}

// Get news by id
func (u *newsUC) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error) {
	newsBase, err := u.redisRepo.GetNewsByIDCtx(ctx, u.getKeyWithPrefix(newsID.String()))
	if err != nil {
		u.logger.Errorf(ctx, "newsUC.GetNewsByID.GetNewsByIDCtx: %v", err)
	}
	if newsBase != nil {
		return newsBase, nil
	}

	n, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetNewsCtx(ctx, u.getKeyWithPrefix(newsID.String()), cacheDuration, n); err != nil {
		u.logger.Errorf(ctx, "newsUC.GetNewsByID.SetNewsCtx: %s", err)
	}

	return n, nil
}

// Delete news
func (u *newsUC) Delete(ctx context.Context, newsID uuid.UUID) error {
	newsByID, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String(), u.logger); err != nil {
		return errors.NewError(http.StatusForbidden, errors.ErrForbidden, errors.WithMessage(err, "newsUC.Delete.ValidateIsOwner"))
	}

	if err = u.newsRepo.Delete(ctx, newsID); err != nil {
		return err
	}

	if err = u.redisRepo.DeleteNewsCtx(ctx, u.getKeyWithPrefix(newsID.String())); err != nil {
		u.logger.Errorf(ctx, "newsUC.Delete.DeleteNewsCtx: %v", err)
	}

	return nil
}

// Get news
func (u *newsUC) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {
	results, err := u.newsRepo.GetNews(ctx, pq)
	if err != nil {
		u.logger.Error(ctx, err)
	}
	return results, err
}

func (u *newsUC) getKeyWithPrefix(newsID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, newsID)
}
