package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var validate = validator.New()

func ValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("validation", validate)
		c.Next()
	}
}
