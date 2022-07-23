package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		adminId := session.Get("adminId")
		orgId := session.Get("orgId")

		if adminId == nil || orgId == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": fmt.Errorf("user not authorized").Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Set("adminID", adminId)
		c.Set("orgID", orgId)

		session.Set("adminId", adminId)
		session.Set("orgId", orgId)
		if err := session.Save(); err != nil {
			log.Printf("Failed to recreate the session %v\n", err.Error())
		}

		c.Next()
	}
}
