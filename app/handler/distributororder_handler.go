package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	request "github.com/meneketehe/hehe/app/request/organization"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetAllOutgoingRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceOrders, err := h.riceOrderService.GetAllOutgoingRiceOrder(channelID, orgID)
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
		"data":    response.RiceOrdersResponse(riceOrders),
	})
}

func (h *Handler) GetAllIncomingDistributedRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceOrders, err := h.riceOrderService.GetAllIncomingRiceOrder(channelID, orgID)
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
		"data":    response.RiceOrdersResponse(riceOrders),
	})
}

func (h *Handler) GetRiceOrder(c *gin.Context) {
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	riceOrder, err := h.riceOrderService.GetRiceOrderByID(channelID, orderID)
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
		"data":    response.RiceOrderResponse(riceOrder),
	})
}

func (h *Handler) CreateRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")

	if role != "distributor" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only distributor role can create this order, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	var req request.CreateRiceOrderRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	manufacturer, err := h.manufacturerService.GetManufacturerByID(channelID, req.ManufacturerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	rice, err := h.riceService.GetRiceByID(channelID, req.RiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if rice.ManufacturerID != manufacturer.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("this rice does not belongs to %s", manufacturer.Name).Error(),
			"data":    nil,
		})
		return
	}

	input := &model.RiceOrder{
		OrdererID: orgID,
		SellerID:  req.ManufacturerID,
		RiceID:    req.RiceID,
		Quantity:  req.Quantity,
		Order:     model.NewOrder(time.Now()),
	}
	riceOrder, err := h.riceOrderService.CreateRiceOrder(channelID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "create order success",
		"data":    response.RiceOrderResponse(riceOrder),
	})
}

func (h *Handler) ReceiveRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if role != "distributor" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only distributor role can receive rice order, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	riceOrder, err := h.riceOrderService.GetRiceOrderByID(channelID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceOrder.OrdererID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("this order is not yours").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceOrderService.ReceiveRiceOrder(channelID, riceOrder, time.Now())
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
		"message": "order received",
		"data":    nil,
	})
}

func (h *Handler) AcceptDistributedRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if role != "distributor" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only distributor role can accept distributed rice order, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	riceOrder, err := h.riceOrderService.GetRiceOrderByID(channelID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceOrder.SellerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller that retailer ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceOrderService.AcceptDistributionRiceOrder(channelID, riceOrder, time.Now())
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
		"message": "order accepted",
		"data":    nil,
	})
}

func (h *Handler) ShipDistributedRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if role != "distributor" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only distributor role can ship distribution of rice order, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	var req request.ShipRiceOrderRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	riceOrder, err := h.riceOrderService.GetRiceOrderByID(channelID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceOrder.SellerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller that retailer ordered").Error(),
			"data":    nil,
		})
		return
	}

	riceStockpile, err := h.riceStockpileService.GetRiceStockpileByVendorIDAndRiceID(channelID, orgID, riceOrder.RiceID)
	if err != nil || riceStockpile.Stock < riceOrder.Quantity {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("rice out of stock").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceOrderService.ShipDistributionRiceOrder(channelID, riceOrder, time.Now(), req.Grade, req.MillingDate, req.StorageTemperature, req.StorageHumidity)
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
		"message": "order shipped",
		"data":    nil,
	})
}
