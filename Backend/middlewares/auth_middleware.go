package middlewares

import (
	"GoTransact/apps/accounts"
	models "GoTransact/apps/base"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Status:  http.StatusUnauthorized,
				Message: "error",
				Data:    map[string]interface{}{"data": "unauthorized request"},
			})
			c.Abort()
			return

		}

		user, err := utils.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Response{
				Status:  http.StatusUnauthorized,
				Message: "error",
				Data:    map[string]interface{}{"data": "unauthorized request"},
			})
			c.Abort()
			return

		}
		fmt.Println("a======================", c.Keys)
		c.Set("user", user)
		fmt.Println("b-------------------", c.Keys)
		c.Next()
	}
}
