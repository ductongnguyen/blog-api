package graphql

import (
	"context"
	"time"

	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/auth"
	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/ductong169z/shorten-url/pkg/logger"
	"github.com/ductong169z/shorten-url/pkg/utils"
)

// Resolver is the GraphQL resolver for auth operations
type Resolver struct {
	cfg     *config.Config
	usecase auth.UseCase
	logger  logger.Logger
}

// NewResolver creates a new auth GraphQL resolver
func NewResolver(cfg *config.Config, usecase auth.UseCase, logger logger.Logger) *Resolver {
	return &Resolver{
		cfg:     cfg,
		usecase: usecase,
		logger:  logger,
	}
}

// User resolves the user query
func (r *Resolver) User(ctx context.Context, id int) (*UserResponse, error) {
	user, err := r.usecase.GetUserByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return fromUserModel(user), nil
}

// Login resolves the login mutation
func (r *Resolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	loginAttempt := &models.User{
		Username: input.Username,
		Password: input.Password,
	}

	user, err := r.usecase.Login(ctx, loginAttempt)
	if err != nil {
		return nil, mapError(err)
	}

	tokenString, expiredAt, err := utils.GenerateJWTToken(user, r.cfg)
	if err != nil {
		return nil, mapError(err)
	}

	refreshToken, refreshTokenExpiresAt, err := r.usecase.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, mapError(err)
	}

	return &AuthResponse{
		Token:                 tokenString,
		ExpiresAt:             formatTime(expiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: formatTime(refreshTokenExpiresAt),
		User:                  fromUserModel(user),
	}, nil
}

// Register resolves the register mutation
func (r *Resolver) Register(ctx context.Context, input RegisterInput) (*UserResponse, error) {
	newUser := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Role:     models.UserRole(input.Role),
	}

	user, err := r.usecase.Register(ctx, newUser)
	if err != nil {
		return nil, mapError(err)
	}

	return fromUserModel(user), nil
}

// RefreshToken resolves the refreshToken mutation
func (r *Resolver) RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenResponse, error) {
	rt, err := r.usecase.ValidateRefreshToken(ctx, input.RefreshToken)
	if err != nil || rt.Revoked {
		return nil, mapError(auth.ErrInvalidToken)
	}

	user, err := r.usecase.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, mapError(err)
	}

	tokenString, expiredAt, err := utils.GenerateJWTToken(user, r.cfg)
	if err != nil {
		return nil, mapError(err)
	}

	newRefreshToken, refreshExp, err := r.usecase.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, mapError(err)
	}

	return &RefreshTokenResponse{
		Token:           tokenString,
		ExpiresAt:       formatTime(expiredAt),
		RefreshToken:    newRefreshToken,
		RefreshExpiresAt: formatTime(refreshExp),
	}, nil
}

// Helper functions
func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func mapError(err error) error {
	// You can add custom error mapping here if needed
	return err
}
