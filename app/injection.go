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

	seedRepository := repository.NewSeedRepository(d.Gateway)
	organizationRepository := repository.NewOrganizationRepository(d.Gateway)

	seedService := service.NewSeedService(&service.SeedServiceConfig{
		SeedRepository: seedRepository,
	})

	organizationService := service.NewOrganizationService(&service.OrganizationServiceConfig{
		OrganizationRepository: organizationRepository,
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
		SeedService:         seedService,
		OrganizationService: organizationService,
		MaxBodyBytes:        mbb,
	})

	return router, nil
}
