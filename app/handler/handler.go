package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/handler/middleware"
	"github.com/meneketehe/hehe/app/model"
)

type Handler struct {
	organizationService model.OrganizationService
	seedService         model.SeedService
	riceGrainService    model.RiceGrainService
	riceService         model.RiceService
	MaxBodyBytes        int64
}

type Config struct {
	R                   *gin.Engine
	OrganizationService model.OrganizationService
	SeedService         model.SeedService
	RiceGrainService    model.RiceGrainService
	RiceService         model.RiceService
	MaxBodyBytes        int64
}

func NewHandler(c *Config) {
	h := &Handler{
		organizationService: c.OrganizationService,
		seedService:         c.SeedService,
		riceGrainService:    c.RiceGrainService,
		riceService:         c.RiceService,
		MaxBodyBytes:        c.MaxBodyBytes,
	}

	c.R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "no route found.",
			"data":    nil,
		})
	})

	ag := c.R.Group("api/organizations/account")
	ag.Use(middleware.AuthOrganization())
	ag.GET("/me", h.GetMeAsOrganization)

	sg := c.R.Group("api/organizations/channels/:channelID/seeds")
	sg.Use(middleware.AuthOrganization())
	sg.GET("", h.GetSeeds)
	sg.GET("/:seedID", h.GetSeed)
	sg.POST("", h.CreateSeed)
	sg.PUT("/:seedID", h.UpdateSeed)
	sg.DELETE("/:seedID", h.DeleteSeed)

	rgg := c.R.Group("api/organizations/channels/:channelID/rice-grains")
	rgg.Use(middleware.AuthOrganization())
	rgg.GET("", h.GetRiceGrains)
	rgg.GET("/:riceGrainID", h.GetRiceGrain)
	rgg.POST("", h.CreateRiceGrain)
	rgg.PUT("/:riceGrainID", h.UpdateRiceGrain)
	rgg.DELETE("/:riceGrainID", h.DeleteRiceGrain)

	rg := c.R.Group("api/organizations/channels/:channelID/rices")
	rg.Use(middleware.AuthOrganization())
	rg.GET("", h.GetRices)
	rg.GET("/:riceID", h.GetRice)
	rg.POST("", h.CreateRice)
	rg.PUT("/:riceID", h.UpdateRice)
	rg.DELETE("/:riceID", h.DeleteRice)
}
