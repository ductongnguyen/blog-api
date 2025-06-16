package server

import (
	"github.com/ductong169z/shorten-url/pkg/cache/redis"
	"github.com/ductong169z/shorten-url/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Option -.
type Option func(*Server)

func FiberEngine(gin *gin.Engine) Option {
	return func(s *Server) {
		s.gin = gin
	}
}

func Redis(rdb redis.Client) Option {
	return func(s *Server) {
		s.redis = rdb
	}
}

func Logger(logger logger.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}
