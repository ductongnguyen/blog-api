package graphql

import (
	"context"

	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/ductong169z/shorten-url/internal/shortener"
	"github.com/ductong169z/shorten-url/pkg/logger"
)

// Resolver is the GraphQL resolver for shortener operations
type Resolver struct {
	cfg     *config.Config
	usecase shortener.UseCase
	logger  logger.Logger
}

// NewResolver creates a new shortener GraphQL resolver
func NewResolver(cfg *config.Config, usecase shortener.UseCase, logger logger.Logger) *Resolver {
	return &Resolver{
		cfg:     cfg,
		usecase: usecase,
		logger:  logger,
	}
}

// ResolveShortCode resolves a short code to its original URL
func (r *Resolver) ResolveShortCode(ctx context.Context, code string) (*ShortURLResponse, error) {
	shortURL, err := r.usecase.ResolveShortCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return FromShortURLModel(shortURL, r.cfg.Server.AppDomain), nil
}

// ShortenURL creates a shortened URL
func (r *Resolver) ShortenURL(ctx context.Context, input ShortenURLInput) (*ShortURLResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Create a new ShortURL model
	shortURL := &models.ShortURL{
		OriginalURL: input.OriginalURL,
		ShortCode:   input.ShortCode,
	}

	// Get client IP and user agent if available from context
	// In a real implementation, you would extract these from the request
	// For now, we'll use placeholders
	clientIP := "127.0.0.1"
	userAgent := "GraphQL Client"
	shortURL.CreatorIP = &clientIP
	shortURL.UserAgent = &userAgent

	result, err := r.usecase.ShortenURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	return FromShortURLModel(result, r.cfg.Server.AppDomain), nil
}
