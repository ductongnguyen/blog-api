package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/models"
	"github.com/ductong169z/shorten-url/pkg/errors"
	"github.com/ductong169z/shorten-url/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// JWT way of auth using cookie or Authorization header
func (mw *MiddlewareManager) AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		ctx := c.Request.Context()
		mw.logger.Infof(ctx, "auth middleware header %s", tokenString)
		if err := mw.validateJWTToken(tokenString, c, mw.cfg); err != nil {
			mw.logger.Error(ctx, "middleware validateJWTToken", zap.String("headerJWT", err.Error()))
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
			c.Abort()
		}
		c.Next()
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, c *gin.Context, cfg *config.Config) error {
	if tokenString == "" {
		return errors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIdFloat, ok := claims["id"].(float64)
		if !ok {
			return errors.InvalidJWTClaims
		}
		userId := int(userIdFloat)
		
		userName, ok := claims["username"].(string)
		if !ok {
			return errors.InvalidJWTClaims
		}
		email, ok := claims["email"].(string)
		if !ok {
			return errors.InvalidJWTClaims
		}
		role, ok := claims["role"].(string)
		if !ok {
			return errors.InvalidJWTClaims
		}
		roleStr, err := models.ParseUserRole(role)
		if err != nil {
			return errors.InvalidJWTClaims
		}

		userData := &models.User{
			ID:       userId,
			Username: userName,
			Email:    email,
			Role:     roleStr,
		}

		ctx := context.WithValue(c.Request.Context(), utils.UserCtxKey{}, userData)
		c.Request = c.Request.WithContext(ctx)
	}
	return nil
}
