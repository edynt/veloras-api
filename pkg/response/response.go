package response

import (
	"net/http"

	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/gin-gonic/gin"
)

// StandardRes response format
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`          // data return null not show
	Error     interface{} `json:"error_details,omitempty"` // Error return null not show
	ErrorFlag bool        `json:"error"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{
		Code:      200,
		Message:   msg.Success,
		Data:      data,
		ErrorFlag: false,
	})
}

func ErrorResponse(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, APIResponse{
		Code:      code,
		Message:   message,
		Error:     err,
		ErrorFlag: true,
	})
}

type HandlerFunc func(ctx *gin.Context) (res interface{}, err error)

func Wrap(handler HandlerFunc) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		res, err := handler(ctx)
		if err != nil {
			if apiErr, ok := err.(*APIError); ok {
				ErrorResponse(ctx, apiErr.StatusCode, apiErr.Message, apiErr.Err)
			} else {
				ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err)
			}
			return
		}
		SuccessResponse(ctx, res)
	}
}
