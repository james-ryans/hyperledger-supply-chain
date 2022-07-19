package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userresponse "github.com/meneketehe/hehe/app/response/user"
)

func (h *Handler) GetScanHistories(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	scanHistories, err := h.scanHistoryService.GetAllScanHistories(userID)
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
		"data":    userresponse.ScanHistoriesResponse(scanHistories),
	})
}
