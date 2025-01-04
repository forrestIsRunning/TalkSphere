package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LimitSize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5<<20) // 5MB
		c.Next()
	}
}
