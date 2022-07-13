package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetSuppliers(c *gin.Context) {
	channelID := c.Param("channelID")

	suppliers, err := h.supplierService.GetAllSuppliers(channelID)
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
		"data":    response.SuppliersResponse(suppliers),
	})
}

func (h *Handler) GetSupplier(c *gin.Context) {
	supplierID := c.Param("supplierID")
	channelID := c.Param("channelID")

	supplier, err := h.supplierService.GetSupplierByID(channelID, supplierID)
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
		"data":    response.SupplierResponse(supplier),
	})
}
