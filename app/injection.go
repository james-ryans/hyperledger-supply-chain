package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/handler"
	repository "github.com/meneketehe/hehe/app/repository/organization"
	service "github.com/meneketehe/hehe/app/service/organization"
	cors "github.com/rs/cors/wrapper/gin"
)

func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting services")

	organizationRepository := repository.NewOrganizationRepository(d.Gateway)
	seedRepository := repository.NewSeedRepository(d.Gateway)
	riceGrainRepository := repository.NewRiceGrainRepository(d.Gateway)
	riceRepository := repository.NewRiceRepository(d.Gateway)

	organizationService := service.NewOrganizationService(&service.OrganizationServiceConfig{
		OrganizationRepository: organizationRepository,
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
		R:                   router,
		OrganizationService: organizationService,
		SeedService:         seedService,
		RiceGrainService:    riceGrainService,
		RiceService:         riceService,
		MaxBodyBytes:        mbb,
	})

	return router, nil
}
