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

func (h *Handler) GetAllOutgoingSeedOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	seedOrders, err := h.seedOrderService.GetAllOutgoingSeedOrder(channelID, orgID)
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

func (h *Handler) GetAllIncomingRiceGrainOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceGrainOrders, err := h.riceGrainOrderService.GetAllIncomingRiceGrainOrder(channelID, orgID)
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

func (h *Handler) GetAllAcceptedIncomingRiceGrainOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceGrainOrders, err := h.riceGrainOrderService.GetAllAcceptedIncomingRiceGrainOrder(channelID, orgID)
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

func (h *Handler) GetSeedOrder(c *gin.Context) {
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	seedOrder, err := h.seedOrderService.GetSeedOrderByID(channelID, orderID)
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
		"data":    response.SeedOrderResponse(seedOrder),
	})
}

func (h *Handler) CreateSeedOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")

	if role != "producer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only producer role can create this order, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	var req request.CreateSeedOrderRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	supplier, err := h.supplierService.GetSupplierByID(channelID, req.SupplierID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	seed, err := h.seedService.GetSeedByID(channelID, req.SeedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if seed.SupplierID != supplier.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("this seed does not belongs to %s", supplier.Name).Error(),
			"data":    nil,
		})
		return
	}

	riceGrainOrder, err := h.riceGrainOrderService.GetRiceGrainOrderByID(channelID, req.RiceGrainOrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceGrainOrder.Status != enum.OrderAccepted {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("this rice grain order does not accepted").Error(),
			"data":    nil,
		})
		return
	}

	if riceGrainOrder.SellerID != orgID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller of this order").Error(),
			"data":    nil,
		})
		return
	}

	input := &model.SeedOrder{
		OrdererID:        orgID,
		SellerID:         req.SupplierID,
		SeedID:           req.SeedID,
		RiceGrainOrderID: req.RiceGrainOrderID,
		Weight:           req.Weight,
		Order:            model.NewOrder(time.Now()),
	}
	seedOrder, err := h.seedOrderService.CreateSeedOrder(channelID, input)
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
		"data":    response.SeedOrderResponse(seedOrder),
	})
}

func (h *Handler) ReceiveSeedOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if role != "producer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only producer role can receive seed order, you are %s", role).Error(),
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

	if seedOrder.OrdererID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("this order is not yours").Error(),
			"data":    nil,
		})
		return
	}

	err = h.seedOrderService.ReceiveSeedOrder(channelID, seedOrder, time.Now())
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

func (h *Handler) AcceptRiceGrainOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if role != "producer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only producer role can accept rice grain order, you are %s", role).Error(),
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

	if riceGrainOrder.SellerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller that manufacturer ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceGrainOrderService.AcceptRiceGrainOrder(channelID, riceGrainOrder, time.Now())
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

func (h *Handler) RejectRiceGrainOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if role != "producer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only producer role can reject rice grain order, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	var req request.RejectRiceGrainOrderRequest
	if ok := bindData(c, &req); !ok {
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

	if riceGrainOrder.SellerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller that manufacturer ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceGrainOrderService.RejectRiceGrainOrder(channelID, riceGrainOrder, time.Now(), req.Reason)
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

func (h *Handler) ShipRiceGrainOrder(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")
	orderID := c.Param("orderID")

	if role != "producer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only producer role can ship rice grain order, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	var req request.ShipRiceGrainOrderRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	riceGrainOrder, err := h.riceGrainOrderService.GetRiceGrainOrderByID(channelID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceGrainOrder.SellerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("you are not the seller that manufacturer ordered").Error(),
			"data":    nil,
		})
		return
	}

	err = h.riceGrainOrderService.ShipRiceGrainOrder(channelID, riceGrainOrder, time.Now(), req.PlowMethod, req.SowMethod, req.Irrigation, req.Fertilization, req.PlantDate, req.HarvestDate, req.StorageTemperature, req.StorageHumidity)
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
