package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	request "github.com/meneketehe/hehe/app/request/organization"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetAllIncomingSeedOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	seedOrders, err := h.seedOrderService.GetAllIncomingSeedOrder(channelID, orgID)
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
		"data":    response.SeedOrdersResponse(seedOrders),
	})
}

func (h *Handler) AcceptSeedOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if me, err := h.organizationService.GetMe(); err != nil || me.Type != "supplier" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only supplier role can accept seed order, you are %s", me.Type).Error(),
			"data":    nil,
		})
		return
	}

	seedOrder, err := h.seedOrderService.GetSeedOrderByID(channelID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if seedOrder.SellerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller that producer ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.seedOrderService.AcceptSeedOrder(channelID, seedOrder, time.Now())
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

func (h *Handler) RejectSeedOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if me, err := h.organizationService.GetMe(); err != nil || me.Type != "supplier" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only supplier role can reject seed order, you are %s", me.Type).Error(),
			"data":    nil,
		})
		return
	}

	var req request.RejectProducerOrderRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	seedOrder, err := h.seedOrderService.GetSeedOrderByID(channelID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if seedOrder.SellerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller that producer ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.seedOrderService.RejectSeedOrder(channelID, seedOrder, time.Now(), req.Reason)
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
		"message": "order rejected",
		"data":    nil,
	})
}

func (h *Handler) ShipSeedOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if me, err := h.organizationService.GetMe(); err != nil || me.Type != "supplier" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only supplier role can ship seed order, you are %s", me.Type).Error(),
			"data":    nil,
		})
		return
	}

	seedOrder, err := h.seedOrderService.GetSeedOrderByID(channelID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if seedOrder.SellerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller that producer ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.seedOrderService.ShipSeedOrder(channelID, seedOrder, time.Now())
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
