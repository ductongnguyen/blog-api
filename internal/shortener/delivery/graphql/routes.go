package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/shortener"
	"github.com/ductong169z/shorten-url/pkg/logger"
	"github.com/gin-gonic/gin"
)

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// RegisterGraphQLRoutes godoc
// @Summary      Register GraphQL routes for URL shortener
// @Description  Registers GraphQL endpoints for URL shortening operations
// @Tags         shortener, graphql
// @Accept       json
// @Produce      json
// @Param        query   body      string  true  "GraphQL query" example="mutation shortenURL { shortenURL(input: { originalURL: \"https://example.com\", shortCode: \"example\" }) { id originalURL shortCode createdAt updatedAt fullShortURL } }"
// @Param        variables body    object  false "GraphQL variables" example="{\"input\": {\"originalURL\": \"https://example.com\", \"shortCode\": \"example\"}}"
// @Param        operationName body string  false "GraphQL operation name" example="shortenURL"
// @Success      200      {object}  object "Successful response"
// @Failure      400,404  {object}  object "Error response"
// @Router       /api/v1/graphql/shortener [post]
func RegisterGraphQLRoutes(router *gin.RouterGroup, cfg *config.Config, usecase shortener.UseCase, logger logger.Logger) {
	handler := NewHandler(cfg, usecase, logger)
	
	// Register playground route (for development)
	// @Summary      Shortener GraphQL Playground
	// @Description  Interactive GraphQL playground for URL shortening operations
	// @Description  Example queries:
	// @Description  1. Shorten URL: mutation shortenURL { shortenURL(input: { originalURL: "https://example.com", shortCode: "example" }) { id originalURL shortCode createdAt updatedAt fullShortURL } }
	// @Description  2. Resolve Short Code: query resolveShortCode { resolveShortCode(code: "example") { id originalURL shortCode createdAt updatedAt fullShortURL } }
	// @Tags         shortener, graphql
	// @Accept       html
	// @Produce      html
	// @Success      200  {string}  html "GraphQL Playground UI"
	// @Router       /api/v1/graphql/shortener/playground [get]
	router.GET("/playground", func(c *gin.Context) {
		playground := playground.Handler("Shortener GraphQL Playground", "/api/v1/graphql/shortener")
		playground.ServeHTTP(c.Writer, c.Request)
	})
	
	// Register GraphQL endpoint
	router.POST("", func(c *gin.Context) {
		var req GraphQLRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{"message": "invalid request format"}}})
			return
		}
		
		result, err := handler.handleGraphQLOperation(c.Request.Context(), req.Query, req.Variables, req.OperationName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"errors": []gin.H{{"message": err.Error()}}})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"data": map[string]interface{}{req.OperationName: result}})
	})
}
