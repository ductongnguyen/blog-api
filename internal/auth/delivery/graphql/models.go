package graphql

import (
	"github.com/ductong169z/shorten-url/internal/models"
)

// Input types
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken"`
}

// Response types
type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AuthResponse struct {
	Token                 string       `json:"token"`
	ExpiresAt             string       `json:"expiresAt"`
	RefreshToken          string       `json:"refreshToken"`
	RefreshTokenExpiresAt string       `json:"refreshTokenExpiresAt"`
	User                  *UserResponse `json:"user"`
}

type RefreshTokenResponse struct {
	Token           string `json:"token"`
	ExpiresAt       string `json:"expiresAt"`
	RefreshToken    string `json:"refreshToken"`
	RefreshExpiresAt string `json:"refreshExpiresAt"`
}

// Converter functions
func fromUserModel(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.String(),
		CreatedAt: formatTime(user.CreatedAt),
		UpdatedAt: formatTime(user.UpdatedAt),
	}
}

func fromUserModelList(users []*models.User) []*UserResponse {
	if users == nil {
		return []*UserResponse{}
	}

	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = fromUserModel(user)
	}

	return userResponses
}
