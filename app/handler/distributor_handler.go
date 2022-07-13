package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetDistributors(c *gin.Context) {
	channelID := c.Param("channelID")

	distributors, err := h.distributorService.GetAllDistributors(channelID)
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
		"data":    response.DistributorsResponse(distributors),
	})
}

func (h *Handler) GetDistributor(c *gin.Context) {
	distributorID := c.Param("distributorID")
	channelID := c.Param("channelID")

	distributor, err := h.distributorService.GetDistributorByID(channelID, distributorID)
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
		"data":    response.DistributorResponse(distributor),
	})
}
