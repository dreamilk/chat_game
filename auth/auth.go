package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authToken = "cg_token"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		if c.Request.Header.Get("Upgrade") == "websocket" {
			token = c.Query(authToken)
		} else {
			t, err := c.Cookie(authToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			token = t
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_id", token)

		c.Next()
	}
}
