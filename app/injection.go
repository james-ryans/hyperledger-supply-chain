package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	organizationAccountRepository := repository.NewOrganizationAccountRepository(d.Couch)
	globalOrganizationRepository := repository.NewGlobalOrganizationRepository(d.Couch)
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

	organizationAccountService := service.NewOrganizationAccountService(&service.OrganizationAccountServiceConfig{
		OrganizationAccountRepository: organizationAccountRepository,
	})
	globalOrganizationService := service.NewGlobalOrganizationService(&service.GlobalOrganizationServiceConfig{
		GlobalOrganizationRepository: globalOrganizationRepository,
	})
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

	userAccountRepository := userrepository.NewUserAccountRepository(d.Couch)
	userRepository := userrepository.NewUserRepository(d.Gateway)
	userRiceSackRepository := userrepository.NewRiceSackRepository(d.Gateway)
	scanHistoryRepository := userrepository.NewScanHistoryRepository(d.Gateway)
	commentRepository := userrepository.NewCommentRepository(d.Gateway)

	userAccountService := userservice.NewUserAccountService(&userservice.UserAccountServiceConfig{
		UserAccountRepository: userAccountRepository,
	})
	userService := userservice.NewUserService(&userservice.UserServiceConfig{
		UserRepository: userRepository,
	})
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
	scanHistoryService := userservice.NewScanHistoryService(&userservice.ScanHistoryServiceConfig{
		ScanHistoryRepository: scanHistoryRepository,
	})
	commentService := userservice.NewCommentService(&userservice.CommentServiceConfig{
		CommentRepository: commentRepository,
	})

	router := gin.Default()

	origin := strings.Split(os.Getenv("CORS_ORIGIN"), ",")
	log.Printf("Allow CORS to: %s", origin)
	c := cors.New(cors.Options{
		AllowedOrigins:   origin,
		AllowCredentials: true,
		AllowedMethods:   []string{"OPTION", "GET", "POST", "PUT", "DELETE"},
	})
	router.Use(c)

	secret := os.Getenv("SESSION_SECRET")
	secure, err := strconv.ParseBool(os.Getenv("SESSION_SECURE"))
	if err != nil {
		return nil, err
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	router.Use(sessions.Sessions(os.Getenv("APP_NAME"), store))

	JSON, _ := json.Marshal(sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	log.Printf("sessions: %s\n", JSON)

	maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
	mbb, err := strconv.ParseInt(maxBodyBytes, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse MAX_BODY_BYTES as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		R:                          router,
		OrganizationAccountService: organizationAccountService,
		GlobalOrganizationService:  globalOrganizationService,
		SupplierService:            supplierService,
		ProducerService:            producerService,
		ManufacturerService:        manufacturerService,
		DistributorService:         distributorService,
		RetailerService:            retailerService,
		SeedService:                seedService,
		RiceGrainService:           riceGrainService,
		RiceService:                riceService,
		RiceStockpileService:       riceStockpileService,
		RiceSackService:            riceSackService,
		SeedOrderService:           seedOrderService,
		RiceGrainOrderService:      riceGrainOrderService,
		RiceOrderService:           riceOrderService,
		UserAccountService:         userAccountService,
		UserService:                userService,
		UserRiceSackService:        userRiceSackService,
		ScanHistoryService:         scanHistoryService,
		CommentService:             commentService,
		MaxBodyBytes:               mbb,
	})

	return router, nil
}
