package http

import (
	"github.com/ductong169z/blog-api/config"
	"github.com/ductong169z/blog-api/internal/models"
	"github.com/ductong169z/blog-api/internal/news"
	"github.com/ductong169z/blog-api/pkg/logger"
	"github.com/ductong169z/blog-api/pkg/response"
	"github.com/ductong169z/blog-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// News handlers
type newsHandlers struct {
	cfg    *config.Config
	newsUC news.UseCase
	logger logger.Logger
}

// NewNewsHandlers News handlers constructor
func NewNewsHandlers(cfg *config.Config, newsUC news.UseCase, logger logger.Logger) news.Handlers {
	return &newsHandlers{cfg: cfg, newsUC: newsUC, logger: logger}
}

// Create godoc
// @Summary Create news
// @Description Create news handler
// @Tags News
// @Accept json
// @Produce json
// @Success 201 {object} models.News
// @Router /news/create [post]
func (h newsHandlers) Create(c *gin.Context) {
	n := &models.News{}
	if err := c.Bind(n); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	ctx := c.Request.Context()
	createdNews, err := h.newsUC.Create(ctx, n)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, createdNews)
}

// Update godoc
// @Summary Update news
// @Description Update news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {object} models.News
// @Router /news/{id} [put]
func (h newsHandlers) Update(c *gin.Context) {
	newsUUID, err := uuid.Parse(c.Param("newsId"))
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	n := &models.News{}
	if err = c.Bind(n); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}
	n.NewsID = newsUUID

	ctx := c.Request.Context()
	updatedNews, err := h.newsUC.Update(ctx, n)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, updatedNews)
}

// GetByID godoc
// @Summary Get by id news
// @Description Get by id news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {object} models.News
// @Router /news/{id} [get]
func (h newsHandlers) GetByID(c *gin.Context) {
	newsUUID, err := uuid.Parse(c.Param("newsId"))
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	ctx := c.Request.Context()
	newsByID, err := h.newsUC.GetNewsByID(ctx, newsUUID)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, newsByID)
}

// Delete godoc
// @Summary Delete news
// @Description Delete by id news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {string} string	"ok"
// @Router /news/{id} [delete]
func (h newsHandlers) Delete(c *gin.Context) {
	newsUUID, err := uuid.Parse(c.Param("newsId"))
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	ctx := c.Request.Context()
	if err = h.newsUC.Delete(ctx, newsUUID); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithNoContent(c)
}

// GetNews godoc
// @Summary Get all news
// @Description Get all news with pagination
// @Tags News
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.NewsList
// @Router /news [get]
func (h newsHandlers) GetNews(c *gin.Context) {
	pq, err := utils.GetPaginationFromCtx(c)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	ctx := c.Request.Context()
	newsList, err := h.newsUC.GetNews(ctx, pq)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, newsList)
}
