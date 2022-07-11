package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	request "github.com/meneketehe/hehe/app/request/organization"
)

func (h *Handler) GetSeeds(c *gin.Context) {
	orgId := c.MustGet("orgId").(string)
	channelId := c.Param("channelId")

	seeds, err := h.seedService.GetAllSeeds(channelId, orgId)
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
		"data":    model.ToSeedsResponse(seeds),
	})
}

func (h *Handler) GetSeed(c *gin.Context) {
	channelId := c.Param("channelId")
	seedId := c.Param("seedId")

	seed, err := h.seedService.GetSeedByID(channelId, seedId)
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
		"data":    model.ToSeedResponse(seed),
	})
}

func (h *Handler) CreateSeed(c *gin.Context) {
	var req request.SeedRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	req.Sanitize()

	orgId := c.MustGet("orgId").(string)
	channelId := c.Param("channelId")

	input := &model.Seed{
		SupplierID:  orgId,
		VarietyName: req.VarietyName,
		PlantAge:    req.PlantAge,
		PlantShape:  req.PlantShape,
		PlantHeight: req.PlantHeight,
		LeafShape:   req.LeafShape,
	}

	seed, err := h.seedService.CreateSeed(channelId, input)
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
		"message": "create seed success",
		"data":    model.ToSeedResponse(seed),
	})
}

func (h *Handler) UpdateSeed(c *gin.Context) {
	var req request.SeedRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	req.Sanitize()

	orgId := c.MustGet("orgId").(string)
	channelId := c.Param("channelId")
	seedId := c.Param("seedId")

	seed, err := h.seedService.GetSeedByID(channelId, seedId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if seed.SupplierID != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": errors.New("this seed is not yours"),
			"data":    nil,
		})
		return
	}

	seed.VarietyName = req.VarietyName
	seed.PlantAge = req.PlantAge
	seed.PlantShape = req.PlantShape
	seed.PlantHeight = req.PlantHeight
	seed.LeafShape = req.LeafShape

	if err := h.seedService.UpdateSeed(channelId, seed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "update seed success",
		"data":    model.ToSeedResponse(seed),
	})
}

func (h *Handler) DeleteSeed(c *gin.Context) {
	orgId := c.MustGet("orgId").(string)
	channelId := c.Param("channelId")
	seedId := c.Param("seedId")

	seed, err := h.seedService.GetSeedByID(channelId, seedId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if seed.SupplierID != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": errors.New("this seed is not yours"),
			"data":    nil,
		})
		return
	}

	if err := h.seedService.DeleteSeed(channelId, seedId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "delete seed success",
		"data":    nil,
	})
}
