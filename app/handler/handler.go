package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/handler/middleware"
	"github.com/meneketehe/hehe/app/model"
)

type Handler struct {
	seedService  model.SeedService
	MaxBodyBytes int64
}

type Config struct {
	R            *gin.Engine
	SeedService  model.SeedService
	MaxBodyBytes int64
}

func NewHandler(c *Config) {
	h := &Handler{
		seedService:  c.SeedService,
		MaxBodyBytes: c.MaxBodyBytes,
	}

	c.R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No route found.",
			"data":    nil,
		})
	})

	sg := c.R.Group("api/organizations/channels/:channelId/seeds")
	sg.Use(middleware.AuthOrganization())
	sg.GET("", h.GetSeeds)
	sg.GET("/:seedId", h.GetSeed)
	sg.POST("", h.CreateSeed)
	sg.PUT("/:seedId", h.UpdateSeed)
	sg.DELETE("/:seedId", h.DeleteSeed)
}
