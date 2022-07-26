package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuperadminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role").(string)
		if role != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": fmt.Errorf("you are not superadmin").Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
