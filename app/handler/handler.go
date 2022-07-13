package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/handler/middleware"
	"github.com/meneketehe/hehe/app/model"
)

type Handler struct {
	organizationService   model.OrganizationService
	supplierService       model.SupplierService
	producerService       model.ProducerService
	manufacturerService   model.ManufacturerService
	distributorService    model.DistributorService
	retailerService       model.RetailerService
	seedService           model.SeedService
	riceGrainService      model.RiceGrainService
	riceService           model.RiceService
	seedOrderService      model.SeedOrderService
	riceGrainOrderService model.RiceGrainOrderService
	riceOrderService      model.RiceOrderService
	MaxBodyBytes          int64
}

type Config struct {
	R                     *gin.Engine
	OrganizationService   model.OrganizationService
	SupplierService       model.SupplierService
	ProducerService       model.ProducerService
	ManufacturerService   model.ManufacturerService
	DistributorService    model.DistributorService
	RetailerService       model.RetailerService
	SeedService           model.SeedService
	RiceGrainService      model.RiceGrainService
	RiceService           model.RiceService
	SeedOrderService      model.SeedOrderService
	RiceGrainOrderService model.RiceGrainOrderService
	RiceOrderService      model.RiceOrderService
	MaxBodyBytes          int64
}

func NewHandler(c *Config) {
	h := &Handler{
		organizationService:   c.OrganizationService,
		supplierService:       c.SupplierService,
		producerService:       c.ProducerService,
		manufacturerService:   c.ManufacturerService,
		distributorService:    c.DistributorService,
		retailerService:       c.RetailerService,
		seedService:           c.SeedService,
		riceGrainService:      c.RiceGrainService,
		riceService:           c.RiceService,
		seedOrderService:      c.SeedOrderService,
		riceGrainOrderService: c.RiceGrainOrderService,
		riceOrderService:      c.RiceOrderService,
		MaxBodyBytes:          c.MaxBodyBytes,
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
	ag.POST("/init", h.InitOrganization)

	sug := c.R.Group("api/organizations/channels/:channelID/suppliers")
	sug.Use(middleware.AuthOrganization())
	sug.GET("", h.GetSuppliers)
	sug.GET("/:supplierID", h.GetSupplier)

	prg := c.R.Group("api/organizations/channels/:channelID/producers")
	prg.Use(middleware.AuthOrganization())
	prg.GET("", h.GetProducers)
	prg.GET("/:producerID", h.GetProducer)

	mag := c.R.Group("api/organizations/channels/:channelID/manufacturers")
	mag.Use(middleware.AuthOrganization())
	mag.GET("", h.GetManufacturers)
	mag.GET("/:manufacturerID", h.GetManufacturer)

	dig := c.R.Group("api/organizations/channels/:channelID/distributors")
	dig.Use(middleware.AuthOrganization())
	dig.GET("", h.GetDistributors)
	dig.GET("/:distributorID", h.GetDistributor)

	reg := c.R.Group("api/organizations/channels/:channelID/retailers")
	reg.Use(middleware.AuthOrganization())
	reg.GET("", h.GetRetailers)
	reg.GET("/:retailerID", h.GetRetailer)

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

	sog := c.R.Group("api/organizations/channels/:channelID/supplier-orders")
	sog.Use(middleware.AuthOrganization())
	sog.GET("incoming", h.GetAllIncomingSeedOrder)
	sog.POST("/:orderID/accept", h.AcceptSeedOrder)
	sog.POST("/:orderID/reject", h.RejectSeedOrder)

	pog := c.R.Group("api/organizations/channels/:channelID/producer-orders")
	pog.Use(middleware.AuthOrganization())
	pog.GET("/outgoing", h.GetAllOutgoingSeedOrder)
	pog.GET("/incoming", h.GetAllIncomingRiceGrainOrder)
	pog.GET("/accepted-incoming", h.GetAllAcceptedIncomingRiceGrainOrder)
	pog.GET("/:orderID", h.GetSeedOrder)
	pog.POST("", h.CreateSeedOrder)
	pog.POST("/:orderID/accept", h.AcceptRiceGrainOrder)
	pog.POST("/:orderID/reject", h.RejectRiceGrainOrder)

	mog := c.R.Group("api/organizations/channels/:channelID/manufacturer-orders")
	mog.Use(middleware.AuthOrganization())
	mog.GET("/outgoing", h.GetAllOutgoingManufacturerOrder)
	mog.GET("/incoming", h.GetAllIncomingDistributorOrder)
	mog.GET("/accepted-incoming", h.GetAllAcceptedIncomingDistributorOrder)
	mog.GET("/:orderID", h.GetManufacturerOrder)
	mog.POST("", h.CreateManufacturerOrder)
	mog.POST("/:orderID/accept", h.AcceptDistributorOrder)
	mog.POST("/:orderID/reject", h.RejectDistributorOrder)

	dog := c.R.Group("api/organizations/channels/:channelID/distributor-orders")
	dog.Use(middleware.AuthOrganization())
	dog.GET("/outgoing", h.GetAllOutgoingRiceOrder)
	dog.GET("/:orderID", h.GetRiceOrder)
	dog.POST("", h.CreateRiceOrder)
}
