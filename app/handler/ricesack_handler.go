package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userresponse "github.com/meneketehe/hehe/app/response/user"
)

func (h *Handler) GetRiceSack(c *gin.Context) {
	code := c.Param("code")

	riceSack, err := h.userRiceSackService.FindRiceSackByCode(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"succesS": true,
		"message": nil,
		"data":    userresponse.RiceSackResponse(riceSack),
	})
}
