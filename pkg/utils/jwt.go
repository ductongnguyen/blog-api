package utils

import (
	"time"

	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/golang-jwt/jwt"
)

// JWT Claims struct
type Claims struct {
	Id       int    `json:"id"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// Generate new JWT Token
func GenerateJWTToken(user *models.User, config *config.Config) (string, time.Time, error) {
	// Register the JWT claims, which includes the username and expiry time
	expiredAt := time.Now().Add(time.Minute * 60)
	claims := &Claims{
		Id:       user.ID,
		Role:     user.Role.String(),
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiredAt, nil
}
