package http

import (
	"github.com/ductong169z/blog-api/internal/middleware"
	"github.com/ductong169z/blog-api/internal/news"
	"github.com/gin-gonic/gin"
)

// Map news routes
func MapNewsRoutes(newsGroup *gin.RouterGroup, h news.Handlers, mw *middleware.MiddlewareManager) {
	newsGroup.Use(mw.AuthJWTMiddleware())
	newsGroup.POST("/create", h.Create)
	newsGroup.PUT("/:newsId", h.Update)
	newsGroup.DELETE("/:newsId", h.Delete)
	newsGroup.GET("/:newsId", h.GetByID)
	newsGroup.GET("", h.GetNews)
}
