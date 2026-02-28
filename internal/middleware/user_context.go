package middleware

import (
	"TaskForge/internal/contextkeys"
	"context"

	"github.com/gin-gonic/gin"
)

type UserContextMiddleware struct {
}

func (uc *UserContextMiddleware) UserContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if exists {
			ctx := context.WithValue(c.Request.Context(), contextkeys.UserIDKey, userID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
