package http

import (
	"github.com/ductong169z/shorten-url/internal/shortener"
	"github.com/gin-gonic/gin"
)

func MapRoutes(group *gin.RouterGroup, h shortener.Handlers) {
	group.POST("/shorten", h.Shorten)
	group.GET("/:code", h.Resolve)
}
