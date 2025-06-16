package response

import (
	"net/http"

	"github.com/ductong169z/shorten-url/pkg/errors"
	"github.com/gin-gonic/gin"
)

const (
	CodeOK = 0
)

const (
	MessageOK = "Success"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Result  any    `json:"result,omitempty"`
}

func WithOK(c *gin.Context, data any) {
	WithCode(c, http.StatusOK, data)
}

func WithNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func WithCode(c *gin.Context, code int, data any) {
	c.JSON(code, Response{
		Message: MessageOK,
		Result:  data,
	})
}

func WithError(c *gin.Context, err error) {
	c.JSON(errors.HTTPErrorResponse(err))
}

func WithErrorCode(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Message: message,
	})
}

// You call this in handlers for domain-layer errors
func WithMappedError(c *gin.Context, err error, mapFunc func(error) (int, string)) {
	code, msg := mapFunc(err)
	WithErrorCode(c, code, msg)
}
