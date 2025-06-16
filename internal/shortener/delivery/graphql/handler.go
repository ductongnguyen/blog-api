package graphql

import (
	"context"
	"errors"

	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/shortener"
	"github.com/ductong169z/shorten-url/pkg/logger"
)

// Handler is the GraphQL handler for shortener operations
type Handler struct {
	resolver *Resolver
	logger   logger.Logger
}

// NewHandler creates a new shortener GraphQL handler
func NewHandler(cfg *config.Config, usecase shortener.UseCase, logger logger.Logger) *Handler {
	resolver := NewResolver(cfg, usecase, logger)
	return &Handler{
		resolver: resolver,
		logger:   logger,
	}
}

// handleGraphQLOperation processes GraphQL queries and mutations
func (h *Handler) handleGraphQLOperation(ctx context.Context, query string, variables map[string]interface{}, operationName string) (interface{}, error) {
	// Simple operation detection based on query string
	if operationName == "resolveShortCode" || query == "query resolveShortCode" {
		// Extract code from variables
		code, ok := variables["code"].(string)
		if !ok {
			return nil, errors.New("invalid code parameter")
		}
		return h.resolver.ResolveShortCode(ctx, code)
	} else if operationName == "shortenURL" || query == "mutation shortenURL" {
		// Extract input from variables
		inputVar, ok := variables["input"].(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid input parameter")
		}
		
		// Extract originalURL and shortCode from input
		originalURL, ok := inputVar["originalURL"].(string)
		if !ok {
			return nil, errors.New("originalURL is required")
		}
		
		// ShortCode is optional
		var shortCode string
		if sc, ok := inputVar["shortCode"].(string); ok {
			shortCode = sc
		}
		
		input := ShortenURLInput{
			OriginalURL: originalURL,
			ShortCode:   shortCode,
		}
		
		return h.resolver.ShortenURL(ctx, input)
	}
	
	return nil, errors.New("invalid GraphQL operation")
}
