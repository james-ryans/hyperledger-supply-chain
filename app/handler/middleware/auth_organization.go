package middleware

import "github.com/gin-gonic/gin"

func AuthOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("orgID", "28899ce0-b79d-497b-ae78-fd3b896e0429")
		c.Next()
	}
}
