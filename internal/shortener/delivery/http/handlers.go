package http

import (
	"net/http"

	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/ductong169z/shorten-url/internal/shortener"
	"github.com/ductong169z/shorten-url/pkg/logger"
	"github.com/ductong169z/shorten-url/pkg/response"
	"github.com/gin-gonic/gin"
)

// News handlers
type handlers struct {
	cfg     *config.Config
	usecase shortener.UseCase
	logger  logger.Logger
}

// NewNewsHandlers News handlers constructor
func NewHandlers(cfg *config.Config, usecase shortener.UseCase, logger logger.Logger) shortener.Handlers {
	return &handlers{cfg: cfg, usecase: usecase, logger: logger}
}

// Shorten godoc
// @Summary      Create a shortened URL
// @Description  Generate a short URL for the given original URL
// @Tags         shortener
// @Accept       json
// @Produce      json
// @Param        shortenRequest  body      ShortenRequest  true  "Original URL to shorten"
// @Success      200            {object}  ShortenResponse
// @Failure      400,409        {object}  response.Response
// @Router       /shorten [post]
func (h *handlers) Shorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WithMappedError(c, err, shortener.MapError)
		return
	}
	if err := req.Validate(); err != nil {
		response.WithMappedError(c, err, shortener.MapError)
		return
	}

	creatorIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	shortUrl := models.ShortURL{
		OriginalURL: req.OriginalURL,
		ShortCode:   req.ShortCode,
		CreatorIP:   &creatorIP,
		UserAgent:   &userAgent,
	}
	shortURL, err := h.usecase.ShortenURL(c.Request.Context(), &shortUrl)
	if err != nil {
		response.WithMappedError(c, err, shortener.MapError)
		return
	}

	response.WithOK(c, FromShortURLModel(shortURL, h.cfg.Server.AppDomain))
}

// Resolve godoc
// @Summary      Redirect to original URL
// @Description  Resolve a short code and redirect to the original URL
// @Tags         shortener
// @Param        code   path      string  true  "Short code"
// @Success      302
// @Failure      404    {object}  response.Response
// @Router       /{code} [get]
func (h *handlers) Resolve(c *gin.Context) {
	code := c.Param("code")
	shortURL, err := h.usecase.ResolveShortCode(c.Request.Context(), code)
	if err != nil {
		response.WithMappedError(c, err, shortener.MapError)
		return
	}
	c.Redirect(http.StatusFound, shortURL.OriginalURL)
}
