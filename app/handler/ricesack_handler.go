package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	userresponse "github.com/meneketehe/hehe/app/response/user"
)

func (h *Handler) GetRiceSack(c *gin.Context) {
	session := sessions.Default(c)

	code := c.Param("code")
	fromHistory := c.Query("from_history")

	userID := ""
	if fromHistory != "true" && session.Get("userId") != nil {
		userID = session.Get("userId").(string)
	}

	riceSack, err := h.userRiceSackService.GetRiceSackByCode(userID, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": nil,
		"data":    userresponse.RiceSackResponse(riceSack),
	})
}
