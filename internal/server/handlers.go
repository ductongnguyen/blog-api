package server

import (
	"context"

	authHttp "github.com/ductong169z/blog-api/internal/auth/delivery/http"
	authRepository "github.com/ductong169z/blog-api/internal/auth/repository"
	authUseCase "github.com/ductong169z/blog-api/internal/auth/usecase"
	apiMiddlewares "github.com/ductong169z/blog-api/internal/middleware"
	newsHttp "github.com/ductong169z/blog-api/internal/news/delivery/http"
	newsRepository "github.com/ductong169z/blog-api/internal/news/repository"
	newsUseCase "github.com/ductong169z/blog-api/internal/news/usecase"
	"github.com/ductong169z/blog-api/pkg/metric"
	"github.com/gin-contrib/requestid"
)

// Map Server Handlers
func (s *Server) MapHandlers() error {
	ctx := context.Background()
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf(ctx, "CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		ctx,
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	metrics.SetSkipPath([]string{"readiness"})

	// Init repositories
	nRepo := newsRepository.NewNewsRepository(s.db)
	newsRedisRepo := newsRepository.NewNewsRedisRepo(s.redis)
	authRepo := authRepository.NewRepository(s.db)
	authRedisRepo := authRepository.NewRedisRepo(s.redis)

	// Init useCases
	newsUC := newsUseCase.NewNewsUseCase(s.cfg, nRepo, newsRedisRepo, s.logger)
	authUC := authUseCase.NewUseCase(s.cfg, authRepo, authRedisRepo, s.logger)

	// Init handlers
	newsHandlers := newsHttp.NewNewsHandlers(s.cfg, newsUC, s.logger)
	authHandlers := authHttp.NewHandlers(s.cfg, authUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)

	s.gin.Use(requestid.New())
	s.gin.Use(mw.MetricsMiddleware(metrics))
	s.gin.Use(mw.LoggerMiddleware(s.logger))

	v1 := s.gin.Group("/api/v1")
	newsGroup := v1.Group("/news")
	authGroup := v1.Group("/auth")

	newsHttp.MapNewsRoutes(newsGroup, newsHandlers, mw)
	authHttp.MapRoutes(authGroup, authHandlers, mw)

	return nil
}
