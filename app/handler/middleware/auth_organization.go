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
		role := session.Get("role")

		log.Printf("sessions: adminId=%s orgId=%s role=%s\n", adminId, orgId, role)

		if adminId == nil || orgId == nil || role == nil {
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
		c.Set("role", role)

		session.Set("adminId", adminId)
		session.Set("orgId", orgId)
		session.Set("role", role)
		if err := session.Save(); err != nil {
			log.Printf("Failed to recreate the session %v\n", err.Error())
		}

		c.Next()
	}
}
