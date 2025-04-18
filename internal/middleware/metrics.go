package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/ductong169z/blog-api/pkg/metric"
)

// Prometheus metrics middleware
func (mw *MiddlewareManager) MetricsMiddleware(metrics metric.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// CODE TODO ...
	}
}
