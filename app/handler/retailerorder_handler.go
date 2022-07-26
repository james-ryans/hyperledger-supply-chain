package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	usermodel "github.com/meneketehe/hehe/app/model/user"
	request "github.com/meneketehe/hehe/app/request/organization"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetAllOutgoingDistributedRiceOrder(c *gin.Context) {
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

func (h *Handler) GetDistributedRiceOrder(c *gin.Context) {
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

func (h *Handler) CreateDistributedRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")

	if role != "retailer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only retailer role can create this order, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	var req request.CreateDistributedRiceOrderRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	distributor, err := h.distributorService.GetDistributorByID(channelID, req.DistributorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	riceStockpile, err := h.riceStockpileService.GetRiceStockpileByID(channelID, req.RiceStockpileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if riceStockpile.VendorID != distributor.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("this rice does not belongs to %s", distributor.Name).Error(),
			"data":    nil,
		})
		return
	}

	input := &model.RiceOrder{
		OrdererID: orgID,
		SellerID:  req.DistributorID,
		RiceID:    riceStockpile.RiceID,
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

func (h *Handler) ReceiveDistributedRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if role != "retailer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only retailer role can receive distribution of rice order, you are %s", role).Error(),
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

	err = h.riceOrderService.ReceiveDistributionRiceOrder(channelID, riceOrder, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	riceSacks, err := h.riceSackService.GetAllRiceSackByRiceOrderID(channelID, riceOrder.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	var userRiceSacks []*usermodel.RiceSack
	for _, sack := range riceSacks {
		userRiceSack, err := h.userRiceSackService.TraceRiceSack(channelID, sack)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}

		userRiceSack, err = h.userRiceSackService.CreateRiceSack(userRiceSack)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}

		userRiceSacks = append(userRiceSacks, userRiceSack)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "order received",
		"data":    nil,
	})
}
