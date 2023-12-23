package http

import (
	"net/http"

	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case typesystem.Http:
				c.AbortWithStatusJSON(e.StatusCode, e)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
			}
		}
	}
}
