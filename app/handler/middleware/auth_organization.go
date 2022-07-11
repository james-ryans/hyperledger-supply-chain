package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AuthOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("orgID", os.Getenv("ORG_ID"))
		c.Next()
	}
}
