package middleware

import (
	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct{}

// CorsMiddleware оборачивает обычный http.Handler и добавляет CORS заголовки
func (m *CorsMiddleware) CorsMiddleware(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if origin != "" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD,PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "authorization,content-type,content-length")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	}

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
