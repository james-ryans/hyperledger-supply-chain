package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	"github.com/meneketehe/hehe/app/model/enum"
	request "github.com/meneketehe/hehe/app/request/organization"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetAllOutgoingRiceGrainOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceGrainOrders, err := h.riceGrainOrderService.GetAllOutgoingRiceGrainOrder(channelID, orgID)
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
		"data":    response.RiceGrainOrdersResponse(riceGrainOrders),
	})
}

func (h *Handler) GetAllIncomingRiceOrder(c *gin.Context) {
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

func (h *Handler) GetAllAcceptedIncomingRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceOrders, err := h.riceOrderService.GetAllAcceptedIncomingRiceOrder(channelID, orgID)
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

func (h *Handler) GetIncomingRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceOrders, err := h.riceOrderService.GetAllAcceptedIncomingRiceOrder(channelID, orgID)
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

func (h *Handler) GetRiceGrainOrder(c *gin.Context) {
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	riceGrainOrder, err := h.riceGrainOrderService.GetRiceGrainOrderByID(channelID, orderID)
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
		"data":    response.RiceGrainOrderResponse(riceGrainOrder),
	})
}

func (h *Handler) CreateRiceGrainOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	if me, err := h.organizationService.GetMe(); err != nil || me.Type != "manufacturer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only manufacturer role can create this order, you are %s", me.Type).Error(),
			"data":    nil,
		})
		return
	}

	var req request.CreateRiceGrainOrderRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	producer, err := h.producerService.GetProducerByID(channelID, req.ProducerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	riceGrain, err := h.riceGrainService.GetRiceGrainByID(channelID, req.RiceGrainID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if riceGrain.ProducerID != producer.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("this rice grain does not belongs to %s", producer.Name).Error(),
			"data":    nil,
		})
		return
	}

	riceOrder, err := h.riceOrderService.GetRiceOrderByID(channelID, req.RiceOrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceOrder.Status != enum.OrderAccepted {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("this rice order does not accepted").Error(),
			"data":    nil,
		})
		return
	}

	if riceOrder.SellerID != orgID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller of this order").Error(),
			"data":    nil,
		})
		return
	}

	input := &model.RiceGrainOrder{
		OrdererID:   orgID,
		SellerID:    req.ProducerID,
		RiceGrainID: req.RiceGrainID,
		RiceOrderID: req.RiceOrderID,
		Weight:      req.Weight,
		Order:       model.NewOrder(time.Now()),
	}
	riceGrainOrder, err := h.riceGrainOrderService.CreateRiceGrainOrder(channelID, input)
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
		"data":    response.RiceGrainOrderResponse(riceGrainOrder),
	})
}

func (h *Handler) ReceiveRiceGrainOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if me, err := h.organizationService.GetMe(); err != nil || me.Type != "manufacturer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only manufacturer role can receive rice grain order, you are %s", me.Type).Error(),
			"data":    nil,
		})
		return
	}

	riceGrainOrder, err := h.riceGrainOrderService.GetRiceGrainOrderByID(channelID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceGrainOrder.OrdererID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("this order is not yours").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceGrainOrderService.ReceiveRiceGrainOrder(channelID, riceGrainOrder, time.Now())
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

func (h *Handler) AcceptRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if me, err := h.organizationService.GetMe(); err != nil || me.Type != "manufacturer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only manufacturer role can accept rice order, you are %s", me.Type).Error(),
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
			"message": fmt.Errorf("you are not the seller that distributor ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceOrderService.AcceptRiceOrder(channelID, riceOrder, time.Now())
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

func (h *Handler) RejectRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if me, err := h.organizationService.GetMe(); err != nil || me.Type != "manufacturer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only manufacturer role can reject rice order, you are %s", me.Type).Error(),
			"data":    nil,
		})
		return
	}

	var req request.RejectRiceOrderRequest
	if ok := bindData(c, &req); !ok {
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
			"message": fmt.Errorf("you are not the seller that distributor ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceOrderService.RejectRiceOrder(channelID, riceOrder, time.Now(), req.Reason)
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

func (h *Handler) ShipRiceOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if me, err := h.organizationService.GetMe(); err != nil || me.Type != "manufacturer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only manufacturer role can ship rice order, you are %s", me.Type).Error(),
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
			"message": fmt.Errorf("you are not the seller that distributor ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceOrderService.ShipRiceOrder(channelID, riceOrder, time.Now(), req.Grade, req.MillingDate, req.StorageTemperature, req.StorageHumidity)
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
