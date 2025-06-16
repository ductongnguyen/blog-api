package graphql

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/auth"
	"github.com/ductong169z/shorten-url/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Handler is the GraphQL handler for auth operations
type Handler struct {
	resolver *Resolver
	logger   logger.Logger
}

// NewHandler creates a new auth GraphQL handler
func NewHandler(cfg *config.Config, usecase auth.UseCase, logger logger.Logger) *Handler {
	resolver := NewResolver(cfg, usecase, logger)
	return &Handler{
		resolver: resolver,
		logger:   logger,
	}
}

// RegisterRoutes registers the GraphQL routes with the Gin router
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {

	// Create a custom resolver that will handle all GraphQL operations
	// This is a simplified approach - in a production environment, you'd use gqlgen's code generation
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

		// Handle the GraphQL operation based on the query
		response, err := h.handleGraphQLOperation(c.Request.Context(), request.Query, request.Variables, request.OperationName)
		if err != nil {
			c.JSON(500, gin.H{"errors": []gin.H{{
				"message": err.Error(),
			}}})
			return
		}

		c.JSON(200, response)
	})

	// GraphQL playground (for development)
	router.GET("/playground", func(c *gin.Context) {
		playground := playground.Handler("Auth GraphQL Playground", "/auth/graphql")
		playground.ServeHTTP(c.Writer, c.Request)
	})
}



// handleGraphQLOperation processes GraphQL queries and mutations
func (h *Handler) handleGraphQLOperation(ctx context.Context, query string, variables map[string]interface{}, operationName string) (interface{}, error) {
	// Simple operation detection based on query string
	if operationName == "user" || query == "query user" {
		// Extract user ID from variables
		userID, ok := variables["id"].(int)
		if !ok {
			// Try to convert from float64 (common in JSON)
			if idFloat, ok := variables["id"].(float64); ok {
				userID = int(idFloat)
			} else {
				return nil, errors.New("invalid user ID")
			}
		}
		return h.resolver.User(ctx, userID)
	} else if operationName == "login" || query == "mutation login" {
		// Extract login input from variables
		inputVar, ok := variables["input"].(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid login input")
		}

		input := LoginInput{
			Username: inputVar["username"].(string),
			Password: inputVar["password"].(string),
		}
		return h.resolver.Login(ctx, input)
	} else if operationName == "register" || query == "mutation register" {
		// Extract register input from variables
		inputVar, ok := variables["input"].(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid register input")
		}

		input := RegisterInput{
			Username: inputVar["username"].(string),
			Email:    inputVar["email"].(string),
			Password: inputVar["password"].(string),
			Role:     inputVar["role"].(string),
		}
		return h.resolver.Register(ctx, input)
	} else if operationName == "refreshToken" || query == "mutation refreshToken" {
		// Extract refresh token input from variables
		inputVar, ok := variables["input"].(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid refresh token input")
		}

		input := RefreshTokenInput{
			RefreshToken: inputVar["refreshToken"].(string),
		}
		return h.resolver.RefreshToken(ctx, input)
	}

	return nil, errors.New("invalid GraphQL operation")
}
