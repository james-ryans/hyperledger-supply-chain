package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("userId")

		if id == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": fmt.Errorf("user not authorized").Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		userId := id.(string)
		c.Set("userID", userId)

		session.Set("userId", userId)
		if err := session.Save(); err != nil {
			log.Printf("Failed recreate the session %v\n", err.Error())
		}

		c.Next()
	}
}
