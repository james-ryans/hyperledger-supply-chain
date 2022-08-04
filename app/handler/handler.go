package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/handler/middleware"
	"github.com/meneketehe/hehe/app/model"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type Handler struct {
	organizationAccountService model.OrganizationAccountService
	globalOrganizationService  model.GlobalOrganizationService
	globalChannelService       model.GlobalChannelService
	supplierService            model.SupplierService
	producerService            model.ProducerService
	manufacturerService        model.ManufacturerService
	distributorService         model.DistributorService
	retailerService            model.RetailerService
	seedService                model.SeedService
	riceGrainService           model.RiceGrainService
	riceService                model.RiceService
	riceStockpileService       model.RiceStockpileService
	riceSackService            model.RiceSackService
	seedOrderService           model.SeedOrderService
	riceGrainOrderService      model.RiceGrainOrderService
	riceOrderService           model.RiceOrderService
	userAccountService         usermodel.UserAccountService
	userService                usermodel.UserService
	userRiceSackService        usermodel.RiceSackService
	scanHistoryService         usermodel.ScanHistoryService
	commentService             usermodel.CommentService
	MaxBodyBytes               int64
}

type Config struct {
	R                          *gin.Engine
	OrganizationAccountService model.OrganizationAccountService
	GlobalOrganizationService  model.GlobalOrganizationService
	GlobalChannelService       model.GlobalChannelService
	SupplierService            model.SupplierService
	ProducerService            model.ProducerService
	ManufacturerService        model.ManufacturerService
	DistributorService         model.DistributorService
	RetailerService            model.RetailerService
	SeedService                model.SeedService
	RiceGrainService           model.RiceGrainService
	RiceService                model.RiceService
	RiceStockpileService       model.RiceStockpileService
	RiceSackService            model.RiceSackService
	SeedOrderService           model.SeedOrderService
	RiceGrainOrderService      model.RiceGrainOrderService
	RiceOrderService           model.RiceOrderService
	UserAccountService         usermodel.UserAccountService
	UserService                usermodel.UserService
	UserRiceSackService        usermodel.RiceSackService
	ScanHistoryService         usermodel.ScanHistoryService
	CommentService             usermodel.CommentService
	MaxBodyBytes               int64
}

func NewHandler(c *Config) {
	h := &Handler{
		organizationAccountService: c.OrganizationAccountService,
		globalOrganizationService:  c.GlobalOrganizationService,
		globalChannelService:       c.GlobalChannelService,
		supplierService:            c.SupplierService,
		producerService:            c.ProducerService,
		manufacturerService:        c.ManufacturerService,
		distributorService:         c.DistributorService,
		retailerService:            c.RetailerService,
		seedService:                c.SeedService,
		riceGrainService:           c.RiceGrainService,
		riceService:                c.RiceService,
		riceStockpileService:       c.RiceStockpileService,
		riceSackService:            c.RiceSackService,
		seedOrderService:           c.SeedOrderService,
		riceGrainOrderService:      c.RiceGrainOrderService,
		riceOrderService:           c.RiceOrderService,
		userAccountService:         c.UserAccountService,
		userService:                c.UserService,
		userRiceSackService:        c.UserRiceSackService,
		scanHistoryService:         c.ScanHistoryService,
		commentService:             c.CommentService,
		MaxBodyBytes:               c.MaxBodyBytes,
	}

	c.R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "no route found.",
			"data":    nil,
		})
	})

	uag := c.R.Group("api/users/account")
	uag.POST("/register", h.RegisterUser)
	uag.POST("/login", h.LoginUser)
	uag.POST("/logout", h.LogoutUser)

	uag.Use(middleware.AuthUser())
	uag.GET("/me", h.GetMeAsUser)
	uag.PUT("/edit-profile", h.EditUserProfile)

	ursg := c.R.Group("api/users/rice-sacks")
	ursg.Use(middleware.AuthUser())
	ursg.GET("/:code", h.GetRiceSack)

	ushg := c.R.Group("api/users/scan-histories")
	ushg.Use(middleware.AuthUser())
	ushg.GET("", h.GetScanHistories)

	ucg := c.R.Group("api/users/comments")
	ucg.GET("/:riceID", h.GetComments)
	ucg.Use(middleware.AuthUser())
	ucg.POST("/:riceID", h.WriteComment)

	ag := c.R.Group("api/organizations/account")
	ag.POST("/login", h.LoginOrganization)
	ag.POST("/logout", h.LogoutOrganization)

	ag.Use(middleware.AuthOrganization())
	ag.GET("/me", h.GetMeAsOrganization)

	osg := c.R.Group("api/organizations/organizations")
	osg.Use(middleware.AuthOrganization(), middleware.SuperadminRole())
	osg.GET("", h.GetAllOrganizations)
	osg.GET("/:ID", h.GetOrganization)
	osg.POST("", h.CreateOrganization)

	csg := c.R.Group("api/organizations/channels")
	csg.Use(middleware.AuthOrganization(), middleware.SuperadminRole())
	csg.GET("", h.GetAllChannels)
	csg.GET("/:channelID", h.GetChannel)
	csg.POST("", h.CreateChannel)

	usg := c.R.Group("api/organizations/users")
	usg.Use(middleware.AuthOrganization(), middleware.AuthOrganization())
	usg.GET("", h.GetAllUsers)
	usg.GET("/:ID", h.GetUser)
	usg.POST("", h.CreateUser)
	usg.PUT("/:ID", h.UpdateUser)
	usg.DELETE("/:ID", h.DeleteUser)

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

	rsg := c.R.Group("api/organizations/channels/:channelID/rice-stockpiles")
	rsg.Use(middleware.AuthOrganization())
	rsg.GET("", h.GetAllRiceStockpiles)
	rsg.GET("/:stockID", h.GetRiceStockpile)
	rsg.GET("/rice-sacks/:sackID/print", h.PrintRiceSackQRCode)

	sog := c.R.Group("api/organizations/channels/:channelID/supplier-orders")
	sog.Use(middleware.AuthOrganization())
	sog.GET("incoming", h.GetAllIncomingSeedOrder)
	sog.POST("/:orderID/accept", h.AcceptSeedOrder)
	sog.POST("/:orderID/reject", h.RejectSeedOrder)
	sog.POST("/:orderID/ship", h.ShipSeedOrder)

	pog := c.R.Group("api/organizations/channels/:channelID/producer-orders")
	pog.Use(middleware.AuthOrganization())
	pog.GET("/outgoing", h.GetAllOutgoingSeedOrder)
	pog.GET("/incoming", h.GetAllIncomingRiceGrainOrder)
	pog.GET("/accepted-incoming", h.GetAllAcceptedIncomingRiceGrainOrder)
	pog.GET("/:orderID", h.GetSeedOrder)
	pog.POST("", h.CreateSeedOrder)
	pog.POST("/:orderID/receive", h.ReceiveSeedOrder)
	pog.POST("/:orderID/accept", h.AcceptRiceGrainOrder)
	pog.POST("/:orderID/reject", h.RejectRiceGrainOrder)
	pog.POST("/:orderID/ship", h.ShipRiceGrainOrder)

	mog := c.R.Group("api/organizations/channels/:channelID/manufacturer-orders")
	mog.Use(middleware.AuthOrganization())
	mog.GET("/outgoing", h.GetAllOutgoingRiceGrainOrder)
	mog.GET("/incoming", h.GetAllIncomingRiceOrder)
	mog.GET("/accepted-incoming", h.GetAllAcceptedIncomingRiceOrder)
	mog.GET("/:orderID", h.GetRiceGrainOrder)
	mog.POST("", h.CreateRiceGrainOrder)
	mog.POST("/:orderID/receive", h.ReceiveRiceGrainOrder)
	mog.POST("/:orderID/accept", h.AcceptRiceOrder)
	mog.POST("/:orderID/reject", h.RejectRiceOrder)
	mog.POST("/:orderID/ship", h.ShipRiceOrder)

	dog := c.R.Group("api/organizations/channels/:channelID/distributor-orders")
	dog.Use(middleware.AuthOrganization())
	dog.GET("/outgoing", h.GetAllOutgoingRiceOrder)
	dog.GET("/incoming", h.GetAllIncomingDistributedRiceOrder)
	dog.GET("/:orderID", h.GetRiceOrder)
	dog.POST("", h.CreateRiceOrder)
	dog.POST("/:orderID/receive", h.ReceiveRiceOrder)
	dog.POST("/:orderID/accept", h.AcceptDistributedRiceOrder)
	dog.POST("/:orderID/ship", h.ShipDistributedRiceOrder)

	rog := c.R.Group("api/organizations/channels/:channelID/retailer-orders")
	rog.Use(middleware.AuthOrganization())
	rog.GET("/outgoing", h.GetAllOutgoingDistributedRiceOrder)
	rog.GET("/:orderID", h.GetDistributedRiceOrder)
	rog.POST("", h.CreateDistributedRiceOrder)
	rog.POST("/:orderID/receive", h.ReceiveDistributedRiceOrder)
}
