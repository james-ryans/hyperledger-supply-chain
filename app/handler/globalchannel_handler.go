package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/helper"
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
			"message": fmt.Errorf("channel with %s name already exists", req.Name).Error(),
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

	orgsID := ch.OrgsID()
	orgs, err := h.globalOrganizationService.GetOrganizationsByIDs(orgsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	err = h.globalChannelService.CreateChannelBlock(ch, orgs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	err = h.globalChannelService.JoinChannel(ch.Name, filepath.Join(helper.BasePath, "app", "storage", ch.Name, "genesis.block"), orgs)
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

func (h *Handler) GetJoinedChannels(c *gin.Context) {
	chs, err := h.globalChannelService.GetJoinedChannels()
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
		"data": gin.H{
			"channels": chs,
		},
	})
}

// func (h *Handler) JoinChannel(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
// 		return
// 	}
// 	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ".block"
// 	err = c.SaveUploadedFile(file, filename)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
// 		return
// 	}
//
// 	name, err := h.globalChannelService.GetChannelNameByFile(filename)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
// 		return
// 	}
//
// 	err = h.globalChannelService.JoinChannel(name, filepath.Join(helper.BasePath, "app", filename))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
// 		return
// 	}
//
// 	err = os.Remove(filepath.Join(helper.BasePath, "app", filename))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "join channel success",
// 		"data":    nil,
// 	})
// }
//
// func (h *Handler) CheckActivationChannel(c *gin.Context) {
// 	name := c.Param("channelID")
//
// 	active, err := h.globalChannelService.CheckActivationChannel(name)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data": gin.H{
// 				"active": active,
// 			},
// 		})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": nil,
// 		"data": gin.H{
// 			"active": active,
// 		},
// 	})
// }
//
// func (h *Handler) ActivateChannel(c *gin.Context) {
// 	name := c.Param("channelID")
//
// 	err := h.globalChannelService.ActivateChannel(name)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "activate channel success",
// 		"data":    nil,
// 	})
// }
