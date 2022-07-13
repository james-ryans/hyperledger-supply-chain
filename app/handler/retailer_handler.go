package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetRetailers(c *gin.Context) {
	channelID := c.Param("channelID")

	retailers, err := h.retailerService.GetAllRetailers(channelID)
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
		"data":    response.RetailersResponse(retailers),
	})
}

func (h *Handler) GetRetailer(c *gin.Context) {
	retailerID := c.Param("retailerID")
	channelID := c.Param("channelID")

	retailer, err := h.retailerService.GetRetailerByID(channelID, retailerID)
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
		"data":    response.RetailerResponse(retailer),
	})
}
