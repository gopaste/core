package middlewares

import (
	"net/http"
	"strings"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/token"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func AuthPostMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		parts := strings.Fields(authorizationHeader)
		if len(parts) != 2 {
			c.Set("x-user-id", nil)
			c.Next()
			return
		}

		tokenString := parts[1]

		payload, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.Response{Message: err.Error()})
			return
		}

		c.Set("x-user-id", payload.UserID.String())
		c.Next()
	}
}
