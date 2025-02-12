package auth

import "github.com/gin-gonic/gin"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", "abc123")

		c.Next()
	}
}
