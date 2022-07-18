package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", os.Getenv("USER_ID"))
		c.Next()
	}
}
