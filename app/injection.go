package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/handler"
	repository "github.com/meneketehe/hehe/app/repository/organization"
	userrepository "github.com/meneketehe/hehe/app/repository/user"
	service "github.com/meneketehe/hehe/app/service/organization"
	userservice "github.com/meneketehe/hehe/app/service/user"
	cors "github.com/rs/cors/wrapper/gin"
)

func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting services")

	supplierRepository := repository.NewSupplierRepository(d.Gateway)
	producerRepository := repository.NewProducerRepository(d.Gateway)
	manufacturerRepository := repository.NewManufacturerRepository(d.Gateway)
	distributorRepository := repository.NewDistributorRepository(d.Gateway)
	retailerRepository := repository.NewRetailerRepository(d.Gateway)

	seedRepository := repository.NewSeedRepository(d.Gateway)
	riceGrainRepository := repository.NewRiceGrainRepository(d.Gateway)
	riceRepository := repository.NewRiceRepository(d.Gateway)
	riceStockpileRepository := repository.NewRiceStockpileRepository(d.Gateway)
	riceSackRepository := repository.NewRiceSackRepository(d.Gateway)

	seedOrderRepository := repository.NewSeedOrderRepository(d.Gateway)
	riceGrainOrderRepository := repository.NewRiceGrainOrderRepository(d.Gateway)
	riceOrderRepository := repository.NewRiceOrderRepository(d.Gateway)

	organizationService := service.NewOrganizationService()
	supplierService := service.NewSupplierService(&service.SupplierServiceConfig{
		SupplierRepository: supplierRepository,
	})
	producerService := service.NewProducerService(&service.ProducerServiceConfig{
		ProducerRepository: producerRepository,
	})
	manufacturerService := service.NewManufacturerService(&service.ManufacturerServiceConfig{
		ManufacturerRepository: manufacturerRepository,
	})
	distributorService := service.NewDistributorService(&service.DistributorServiceConfig{
		DistributorRepository: distributorRepository,
	})
	retailerService := service.NewRetailerService(&service.RetailerServiceConfig{
		RetailerRepository: retailerRepository,
	})

	seedService := service.NewSeedService(&service.SeedServiceConfig{
		SeedRepository: seedRepository,
	})
	riceGrainService := service.NewRiceGrainService(&service.RiceGrainServiceConfig{
		RiceGrainRepository: riceGrainRepository,
	})
	riceService := service.NewRiceService(&service.RiceServiceConfig{
		RiceRepository: riceRepository,
	})
	riceStockpileService := service.NewRiceStockpileService(&service.RiceStockpileServiceConfig{
		RiceStockpileRepository: riceStockpileRepository,
	})
	riceSackService := service.NewRiceSackService(&service.RiceSackServiceConfig{
		RiceSackRepository: riceSackRepository,
	})

	seedOrderService := service.NewSeedOrderService(&service.SeedOrderServiceConfig{
		SeedOrderRepository: seedOrderRepository,
	})
	riceGrainOrderService := service.NewRiceGrainOrderService(&service.RiceGrainOrderServiceConfig{
		RiceGrainOrderRepository: riceGrainOrderRepository,
	})
	riceOrderService := service.NewRiceOrderService(&service.RiceOrderServiceConfig{
		RiceOrderRepository: riceOrderRepository,
	})

	userRiceSackRepository := userrepository.NewRiceSackRepository(d.Gateway)
	userRiceSackService := userservice.NewRiceSackService(&userservice.RiceSackServiceConfig{
		SupplierRepository:       supplierRepository,
		ProducerRepository:       producerRepository,
		ManufacturerRepository:   manufacturerRepository,
		DistributorRepository:    distributorRepository,
		RetailerRepository:       retailerRepository,
		SeedRepository:           seedRepository,
		RiceGrainRepository:      riceGrainRepository,
		RiceRepository:           riceRepository,
		SeedOrderRepository:      seedOrderRepository,
		RiceGrainOrderRepository: riceGrainOrderRepository,
		RiceOrderRepository:      riceOrderRepository,
		RiceSackRepository:       userRiceSackRepository,
	})

	router := gin.Default()

	origin := os.Getenv("CORS_ORIGIN")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{origin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	router.Use(c)

	maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
	mbb, err := strconv.ParseInt(maxBodyBytes, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse MAX_BODY_BYTES as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		R:                     router,
		OrganizationService:   organizationService,
		SupplierService:       supplierService,
		ProducerService:       producerService,
		ManufacturerService:   manufacturerService,
		DistributorService:    distributorService,
		RetailerService:       retailerService,
		SeedService:           seedService,
		RiceGrainService:      riceGrainService,
		RiceService:           riceService,
		RiceStockpileService:  riceStockpileService,
		RiceSackService:       riceSackService,
		SeedOrderService:      seedOrderService,
		RiceGrainOrderService: riceGrainOrderService,
		RiceOrderService:      riceOrderService,
		UserRiceSackService:   userRiceSackService,
		MaxBodyBytes:          mbb,
	})

	return router, nil
}
