package middleware

import (
	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, origins: origins, logger: logger}
}
