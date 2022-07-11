package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	request "github.com/meneketehe/hehe/app/request/organization"
)

func (h *Handler) GetSeeds(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	seeds, err := h.seedService.GetAllSeeds(channelID, orgID)
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
	channelID := c.Param("channelID")
	seedID := c.Param("seedID")

	seed, err := h.seedService.GetSeedByID(channelID, seedID)
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

	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")

	input := &model.Seed{
		SupplierID:  orgID,
		VarietyName: req.VarietyName,
		PlantAge:    req.PlantAge,
		PlantShape:  req.PlantShape,
		PlantHeight: req.PlantHeight,
		LeafShape:   req.LeafShape,
	}

	seed, err := h.seedService.CreateSeed(channelID, input)
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

	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	seedID := c.Param("seedID")

	seed, err := h.seedService.GetSeedByID(channelID, seedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if seed.SupplierID != orgID {
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

	if err := h.seedService.UpdateSeed(channelID, seed); err != nil {
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
	orgID := c.MustGet("orgID").(string)
	channelID := c.Param("channelID")
	seedID := c.Param("seedID")

	seed, err := h.seedService.GetSeedByID(channelID, seedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if seed.SupplierID != orgID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": errors.New("this seed is not yours"),
			"data":    nil,
		})
		return
	}

	if err := h.seedService.DeleteSeed(channelID, seedID); err != nil {
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
