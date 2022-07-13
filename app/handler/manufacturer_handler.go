package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetManufacturers(c *gin.Context) {
	channelID := c.Param("channelID")

	manufacturers, err := h.manufacturerService.GetAllManufacturers(channelID)
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
		"data":    response.ManufacturersResponse(manufacturers),
	})
}

func (h *Handler) GetManufacturer(c *gin.Context) {
	manufacturerID := c.Param("manufacturerID")
	channelID := c.Param("channelID")

	manufacturer, err := h.manufacturerService.GetManufacturerByID(channelID, manufacturerID)
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
		"data":    response.ManufacturerResponse(manufacturer),
	})
}
