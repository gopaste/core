package middlewares

import (
	"net/http"
	"strings"

	"github.com/Caixetadev/snippet/internal/token"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func AuthPostMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		parts := strings.Fields(authorizationHeader)
		if len(parts) != 2 {
			c.Set("x-user-id", nil)
			c.Next()
			return
		}

		tokenString := parts[1]

		authorized, err := token.IsAuthorized(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if authorized {
			userID, err := token.ExtractIDFromToken(tokenString, secret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
				c.Abort()
				return
			}

			c.Set("x-user-id", userID)
			c.Next()
			return
		}
	}
}
