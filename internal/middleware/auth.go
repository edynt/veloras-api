package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the request url path
		uri := c.Request.URL.Path
		log.Println("uri request: ", uri)

		// check headers authorization
		jwtToken, valid := utils.ExtractBearerToken(c)
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized, "message": msg.Unauthorized, "error_details": msg.HeaderAuthenticationNotFound, "error": true,
			})

			return
		}
		// validate jwt token by subject
		claims, err := utils.VerifyTokenSubject(jwtToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized, "message": msg.InvalidToken, "error_details": msg.VerifyTokenFailed, "error": true,
			})

			return
		}

		// update claims to context
		ctx := context.WithValue(c.Request.Context(), "subjectUUID", claims.Subject)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
