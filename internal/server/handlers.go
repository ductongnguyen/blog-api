package server

import (
	"context"

	authHttp "github.com/ductong169z/shorten-url/internal/auth/delivery/http"
	authGraphQL "github.com/ductong169z/shorten-url/internal/auth/delivery/graphql"
	authRepository "github.com/ductong169z/shorten-url/internal/auth/repository"
	authUseCase "github.com/ductong169z/shorten-url/internal/auth/usecase"
	apiMiddlewares "github.com/ductong169z/shorten-url/internal/middleware"

	shortHttp "github.com/ductong169z/shorten-url/internal/shortener/delivery/http"
	shortGraphQL "github.com/ductong169z/shorten-url/internal/shortener/delivery/graphql"
	shortRepository "github.com/ductong169z/shorten-url/internal/shortener/repository"
	shortUseCase "github.com/ductong169z/shorten-url/internal/shortener/usecase"

	"github.com/ductong169z/shorten-url/pkg/metric"
	"github.com/gin-contrib/requestid"

	// Swagger UI imports
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	authRepo := authRepository.NewRepository(s.db)
	authRedisRepo := authRepository.NewRedisRepo(s.redis)

	shortRepo := shortRepository.NewRepository(s.db)
	shortRedisRepo := shortRepository.NewRedisRepo(s.redis)

	// Init useCases
	authUC := authUseCase.NewUseCase(s.cfg, authRepo, authRedisRepo, s.logger)

	shortUC := shortUseCase.NewUseCase(s.cfg, shortRepo, shortRedisRepo, s.logger)

	// Init handlers
	authHandlers := authHttp.NewHandlers(s.cfg, authUC, s.logger)
	shortHandlers := shortHttp.NewHandlers(s.cfg, shortUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)

	s.gin.Use(requestid.New())
	s.gin.Use(mw.MetricsMiddleware(metrics))
	s.gin.Use(mw.LoggerMiddleware(s.logger))

	// Swagger docs endpoint
	s.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := s.gin.Group("/api/v1")
	noPrefixGroup := s.gin.Group("")
	authGroup := v1.Group("/auth")
	shortGroup := noPrefixGroup.Group("")
	
	// Create a separate group for GraphQL that doesn't have auth middleware
	graphqlGroup := v1.Group("/graphql")

	// Register HTTP routes
	authHttp.MapRoutes(authGroup, authHandlers, mw)
	shortHttp.MapRoutes(shortGroup, shortHandlers)
	
	// Register GraphQL routes - using a separate group that bypasses auth
	authGraphQL.RegisterGraphQLRoutes(graphqlGroup, s.cfg, authUC, s.logger)
	
	// Create a separate group for shortener GraphQL
	shortGraphQLGroup := v1.Group("/graphql/shortener")
	shortGraphQL.RegisterGraphQLRoutes(shortGraphQLGroup, s.cfg, shortUC, s.logger)

	return nil
}
