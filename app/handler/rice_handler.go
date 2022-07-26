package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	request "github.com/meneketehe/hehe/app/request/organization"
)

func (h *Handler) GetRices(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	rices, err := h.riceService.GetAllRices(channelID, orgID)
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
		"data":    model.ToRicesResponse(rices),
	})
}

func (h *Handler) GetRice(c *gin.Context) {
	channelID := c.Param("channelID")
	riceID := c.Param("riceID")

	rice, err := h.riceService.GetRiceByID(channelID, riceID)
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
		"data":    model.ToRiceResponse(rice),
	})
}

func (h *Handler) CreateRice(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")

	if role != "manufacturer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only manufacturer role can create rice asset, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	var req request.RiceRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	req.Sanitize()

	input := &model.Rice{
		ManufacturerID: orgID,
		BrandName:      req.BrandName,
		Weight:         req.Weight,
		Texture:        req.Texture,
		AmyloseRate:    req.AmyloseRate,
	}

	rice, err := h.riceService.CreateRice(channelID, input)
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
		"message": "create rice success",
		"data":    model.ToRiceResponse(rice),
	})
}

func (h *Handler) UpdateRice(c *gin.Context) {
	var req request.RiceRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	req.Sanitize()

	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	riceID := c.Param("riceID")

	rice, err := h.riceService.GetRiceByID(channelID, riceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if rice.ManufacturerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": errors.New("this rice is not yours").Error(),
			"data":    nil,
		})
		return
	}

	rice.BrandName = req.BrandName
	rice.Weight = req.Weight
	rice.Texture = req.Texture
	rice.AmyloseRate = req.AmyloseRate

	if err := h.riceService.UpdateRice(channelID, rice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "update rice success",
		"data":    model.ToRiceResponse(rice),
	})
}

func (h *Handler) DeleteRice(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	riceID := c.Param("riceID")

	rice, err := h.riceService.GetRiceByID(channelID, riceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if rice.ManufacturerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": errors.New("this rice is not yours"),
			"data":    nil,
		})
		return
	}

	if err := h.riceService.DeleteRice(channelID, riceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "delete rice success",
		"data":    nil,
	})
}
