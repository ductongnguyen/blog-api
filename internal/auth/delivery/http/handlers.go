package http

import (
	"github.com/ductong169z/blog-api/config"
	"github.com/ductong169z/blog-api/internal/auth"
	"github.com/ductong169z/blog-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// News handlers
type handlers struct {
	cfg     *config.Config
	usecase auth.UseCase
	logger  logger.Logger
}

// NewNewsHandlers News handlers constructor
func NewHandlers(cfg *config.Config, usecase auth.UseCase, logger logger.Logger) auth.Handlers {
	return &handlers{cfg: cfg, usecase: usecase, logger: logger}
}

// GetUserByID implements auth.Handlers.
func (h *handlers) GetUserByID(c *gin.Context) {
	panic("unimplemented")
}

// Login implements auth.Handlers.
func (h *handlers) Login(c *gin.Context) {
	panic("unimplemented")
}

// Register implements auth.Handlers.
func (h *handlers) Register(c *gin.Context) {
	panic("unimplemented")
}
