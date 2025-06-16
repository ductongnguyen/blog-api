//go:generate mockgen -source delivery.go -destination mock/handlers_mock.go -package mock
package shortener

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	Shorten(c *gin.Context)
	Resolve(c *gin.Context)
}
