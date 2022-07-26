package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	request "github.com/meneketehe/hehe/app/request/organization"
)

func (h *Handler) GetRiceGrains(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	riceGrains, err := h.riceGrainService.GetAllRiceGrains(channelID, orgID)
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
		"data":    model.ToRiceGrainsResponse(riceGrains),
	})
}

func (h *Handler) GetRiceGrain(c *gin.Context) {
	channelID := c.Param("channelID")
	riceGrainID := c.Param("riceGrainID")

	riceGrain, err := h.riceGrainService.GetRiceGrainByID(channelID, riceGrainID)
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
		"data":    model.ToRiceGrainResponse(riceGrain),
	})
}

func (h *Handler) CreateRiceGrain(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)
	channelID := c.Param("channelID")

	if role != "producer" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Errorf("only producer role can create rice grain asset, you are %s", role).Error(),
			"data":    nil,
		})
		return
	}

	var req request.RiceGrainRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	req.Sanitize()

	input := &model.RiceGrain{
		ProducerID:  orgID,
		VarietyName: req.VarietyName,
		GrainShape:  req.GrainShape,
		GrainColor:  req.GrainColor,
	}

	riceGrain, err := h.riceGrainService.CreateRiceGrain(channelID, input)
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
		"message": "create rice grain success",
		"data":    model.ToRiceGrainResponse(riceGrain),
	})
}

func (h *Handler) UpdateRiceGrain(c *gin.Context) {
	var req request.RiceGrainRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	req.Sanitize()

	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	riceGrainID := c.Param("riceGrainID")

	riceGrain, err := h.riceGrainService.GetRiceGrainByID(channelID, riceGrainID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceGrain.ProducerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": errors.New("this rice grain is not yours").Error(),
			"data":    nil,
		})
		return
	}

	riceGrain.VarietyName = req.VarietyName
	riceGrain.GrainShape = req.GrainShape
	riceGrain.GrainColor = req.GrainColor

	if err := h.riceGrainService.UpdateRiceGrain(channelID, riceGrain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "update rice grain success",
		"data":    model.ToRiceGrainResponse(riceGrain),
	})
}

func (h *Handler) DeleteRiceGrain(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	riceGrainID := c.Param("riceGrainID")

	riceGrain, err := h.riceGrainService.GetRiceGrainByID(channelID, riceGrainID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if riceGrain.ProducerID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": errors.New("this rice grain is not yours"),
			"data":    nil,
		})
		return
	}

	if err := h.riceGrainService.DeleteRiceGrain(channelID, riceGrainID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "delete rice grain success",
		"data":    nil,
	})
}
