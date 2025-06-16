package graphql

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/auth"
	"github.com/ductong169z/shorten-url/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RegisterGraphQLRoutes godoc
// @Summary      Register GraphQL routes for authentication
// @Description  Registers GraphQL endpoints for user authentication operations
// @Tags         auth, graphql
// @Accept       json
// @Produce      json
// @Param        query   body      string  true  "GraphQL query" example="mutation login { login(input: { username: \"testuser\", password: \"password123\" }) { token expiresAt refreshToken refreshTokenExpiresAt user { id username email role } } }"
// @Param        variables body    object  false "GraphQL variables" example="{\"input\": {\"username\": \"testuser\", \"password\": \"password123\"}}"
// @Param        operationName body string  false "GraphQL operation name" example="login"
// @Success      200      {object}  object "Successful response"
// @Failure      400,401  {object}  object "Error response"
// @Router       /api/v1/graphql [post]
func RegisterGraphQLRoutes(router *gin.RouterGroup, cfg *config.Config, usecase auth.UseCase, logger logger.Logger) {
	handler := NewHandler(cfg, usecase, logger)
	
	// Register playground route (for development)
	// @Summary      Auth GraphQL Playground
	// @Description  Interactive GraphQL playground for authentication operations
	// @Description  Example queries:
	// @Description  1. Login: mutation login { login(input: { username: "testuser", password: "password123" }) { token expiresAt refreshToken user { id username email } } }
	// @Description  2. Register: mutation register { register(input: { username: "newuser", email: "new@example.com", password: "password123", role: "user" }) { id username email role } }
	// @Description  3. Get User: query getUser { user(id: 1) { id username email role } }
	// @Description  4. Refresh Token: mutation refreshToken { refreshToken(input: { refreshToken: "your-refresh-token" }) { token expiresAt refreshToken refreshTokenExpiresAt } }
	// @Tags         auth, graphql
	// @Accept       html
	// @Produce      html
	// @Success      200  {string}  html "GraphQL Playground UI"
	// @Router       /api/v1/graphql/playground [get]
	router.GET("/playground", func(c *gin.Context) {
		playground := playground.Handler("Auth GraphQL Playground", "/api/v1/graphql")
		playground.ServeHTTP(c.Writer, c.Request)
	})
	
	// Register GraphQL endpoint for public operations (login, register, refresh token)
	router.POST("/graphql", func(c *gin.Context) {
		var request struct {
			Query         string                 `json:"query"`
			Variables     map[string]interface{} `json:"variables"`
			OperationName string                 `json:"operationName"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"errors": []gin.H{{
				"message": "Invalid request format",
			}}})
			return
		}

		// Handle the GraphQL operation
		response, err := handler.handleGraphQLOperation(c.Request.Context(), request.Query, request.Variables, request.OperationName)
		if err != nil {
			c.JSON(500, gin.H{"errors": []gin.H{{
				"message": err.Error(),
			}}})
			return
		}

		c.JSON(200, response)
	})
}
