package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	request "github.com/meneketehe/hehe/app/request/organization"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetAllChannels(c *gin.Context) {
	chs, err := h.globalChannelService.GetAllChannels()
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
		"data":    response.GlobalChannelsResponse(chs),
	})
}

func (h *Handler) GetChannel(c *gin.Context) {
	ID := c.Param("channelID")

	ch, err := h.globalChannelService.GetChannel(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	orgIDs := append(ch.SuppliersID, ch.ProducersID...)
	orgIDs = append(orgIDs, ch.ManufacturersID...)
	orgIDs = append(orgIDs, ch.DistributorsID...)
	orgIDs = append(orgIDs, ch.RetailersID...)
	orgs, err := h.globalOrganizationService.GetOrganizationsByIDs(orgIDs)
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
		"data":    response.GlobalChannelResponse(ch, orgs),
	})
}

func (h *Handler) CreateChannel(c *gin.Context) {
	var req request.CreateChannelRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	exists, err := h.globalChannelService.CheckNameExists(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("organization with %s name already exists", req.Name).Error(),
			"data":    nil,
		})
		return
	}

	input := &model.GlobalChannel{
		Name:            req.Name,
		SuppliersID:     req.SuppliersID,
		ProducersID:     req.ProducersID,
		ManufacturersID: req.ManufacturersID,
		DistributorsID:  req.DistributorsID,
		RetailersID:     req.RetailersID,
	}
	ch, err := h.globalChannelService.CreateChannel(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	orgIDs := append(ch.SuppliersID, ch.ProducersID...)
	orgIDs = append(orgIDs, ch.ManufacturersID...)
	orgIDs = append(orgIDs, ch.DistributorsID...)
	orgIDs = append(orgIDs, ch.RetailersID...)
	orgs, err := h.globalOrganizationService.GetOrganizationsByIDs(orgIDs)
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
		"message": "create channel success",
		"data":    response.GlobalChannelResponse(ch, orgs),
	})
}
