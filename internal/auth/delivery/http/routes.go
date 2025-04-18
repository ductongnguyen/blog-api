package http

import (
	"github.com/ductong169z/blog-api/internal/middleware"
	"github.com/ductong169z/blog-api/internal/auth"
	"github.com/gin-gonic/gin"
)

// Map news routes
func MapRoutes(group *gin.RouterGroup, h auth.Handlers, mw *middleware.MiddlewareManager) {
	group.Use(mw.AuthJWTMiddleware())
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.GET("/user/:userId", h.GetUserByID)
}
