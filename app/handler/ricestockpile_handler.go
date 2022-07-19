package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/meneketehe/hehe/app/response/organization"
	"github.com/skip2/go-qrcode"
)

func (h *Handler) GetAllRiceStockpiles(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceStockpiles, err := h.riceStockpileService.GetAllRiceStockpile(channelID, orgID)
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
		"data":    response.RiceStockpilesResponse(riceStockpiles),
	})
}

func (h *Handler) GetRiceStockpile(c *gin.Context) {
	channelID := c.Param("channelID")
	stockID := c.Param("stockID")

	pile, err := h.riceStockpileService.GetRiceStockpileByID(channelID, stockID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	sacks, err := h.riceSackService.GetAllRiceSack(channelID, pile.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": nil,
		"data":    response.RiceStockpileResponse(pile, sacks),
	})
}

func (h *Handler) GetAllRiceSack(c *gin.Context) {
	channelID := c.Param("channelID")
	stockID := c.Param("stockID")

	pile, err := h.riceStockpileService.GetRiceStockpileByID(channelID, stockID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	sacks, err := h.riceSackService.GetAllRiceSack(channelID, pile.ID)
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
		"data":    response.RiceSacksResponse(sacks),
	})
}

func (h *Handler) PrintRiceSackQRCode(c *gin.Context) {
	channelID := c.Param("channelID")
	sackID := c.Param("sackID")

	sack, err := h.riceSackService.GetRiceSack(channelID, sackID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	qrCode, err := qrcode.Encode(sack.Code, qrcode.Highest, 256)
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
		"data":    response.RiceSackQrCodeResponse(qrCode),
	})

}
