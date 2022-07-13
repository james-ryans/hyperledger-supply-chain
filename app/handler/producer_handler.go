package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetProducers(c *gin.Context) {
	channelID := c.Param("channelID")

	producers, err := h.producerService.GetAllProducers(channelID)
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
		"data":    response.ProducersResponse(producers),
	})
}

func (h *Handler) GetProducer(c *gin.Context) {
	producerID := c.Param("producerID")
	channelID := c.Param("channelID")

	producer, err := h.producerService.GetProducerByID(channelID, producerID)
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
		"data":    response.ProducerResponse(producer),
	})
}
