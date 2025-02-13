package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chat_game/handlers"
	"chat_game/utils/common"
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
				c.JSON(http.StatusOK, handlers.Resp{
					Code:    -100,
					Message: "Unauthorized",
				})
				c.Abort()
				return
			}
			token = t
		}

		if token == "" {
			c.JSON(http.StatusOK, handlers.Resp{
				Code:    -100,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set(common.UserIDKey, token)

		c.Next()
	}
}
