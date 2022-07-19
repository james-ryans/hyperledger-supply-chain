package userservice

import (
	"encoding/json"
	"log"

	"github.com/meneketehe/hehe/app/model"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type riceSackService struct {
	SupplierRepository       model.SupplierRepository
	ProducerRepository       model.ProducerRepository
	ManufacturerRepository   model.ManufacturerRepository
	DistributorRepository    model.DistributorRepository
	RetailerRepository       model.RetailerRepository
	SeedRepository           model.SeedRepository
	RiceGrainRepository      model.RiceGrainRepository
	RiceRepository           model.RiceRepository
	SeedOrderRepository      model.SeedOrderRepository
	RiceGrainOrderRepository model.RiceGrainOrderRepository
	RiceOrderRepository      model.RiceOrderRepository
	RiceSackRepository       usermodel.RiceSackRepository
}

type RiceSackServiceConfig struct {
	SupplierRepository       model.SupplierRepository
	ProducerRepository       model.ProducerRepository
	ManufacturerRepository   model.ManufacturerRepository
	DistributorRepository    model.DistributorRepository
	RetailerRepository       model.RetailerRepository
	SeedRepository           model.SeedRepository
	RiceGrainRepository      model.RiceGrainRepository
	RiceRepository           model.RiceRepository
	SeedOrderRepository      model.SeedOrderRepository
	RiceGrainOrderRepository model.RiceGrainOrderRepository
	RiceOrderRepository      model.RiceOrderRepository
	RiceSackRepository       usermodel.RiceSackRepository
}

func NewRiceSackService(c *RiceSackServiceConfig) usermodel.RiceSackService {
	return &riceSackService{
		SupplierRepository:       c.SupplierRepository,
		ProducerRepository:       c.ProducerRepository,
		ManufacturerRepository:   c.ManufacturerRepository,
		DistributorRepository:    c.DistributorRepository,
		RetailerRepository:       c.RetailerRepository,
		SeedRepository:           c.SeedRepository,
		RiceGrainRepository:      c.RiceGrainRepository,
		RiceRepository:           c.RiceRepository,
		SeedOrderRepository:      c.SeedOrderRepository,
		RiceGrainOrderRepository: c.RiceGrainOrderRepository,
		RiceOrderRepository:      c.RiceOrderRepository,
		RiceSackRepository:       c.RiceSackRepository,
	}
}

func (s *riceSackService) GetRiceSackByCode(userID, code string) (*usermodel.RiceSack, error) {
	return s.RiceSackRepository.FindByCode(userID, code)
}

func (s *riceSackService) CreateRiceSack(riceSack *usermodel.RiceSack) (*usermodel.RiceSack, error) {
	if err := s.RiceSackRepository.Create(riceSack); err != nil {
		return nil, err
	}

	return riceSack, nil
}

func (s *riceSackService) TraceRiceSack(channelID string, riceSack *model.RiceSack) (*usermodel.RiceSack, error) {
	riceOrder1, err := s.RiceOrderRepository.FindByID(channelID, riceSack.RiceOrderID[1])
	if err != nil {
		return nil, err
	}

	rice1, err := s.RiceRepository.FindByID(channelID, riceOrder1.RiceID)
	if err != nil {
		return nil, err
	}

	retailer, err := s.RetailerRepository.FindByID(channelID, riceOrder1.OrdererID)
	if err != nil {
		return nil, err
	}

	distributor, err := s.DistributorRepository.FindByID(channelID, riceOrder1.SellerID)
	if err != nil {
		return nil, err
	}

	riceOrder0, err := s.RiceOrderRepository.FindByID(channelID, riceSack.RiceOrderID[0])
	if err != nil {
		return nil, err
	}

	rice0, err := s.RiceRepository.FindByID(channelID, riceOrder0.RiceID)
	if err != nil {
		return nil, err
	}

	manufacturer, err := s.ManufacturerRepository.FindByID(channelID, riceOrder0.SellerID)
	if err != nil {
		return nil, err
	}

	riceGrainOrder, err := s.RiceGrainOrderRepository.FindByRiceOrderID(channelID, riceOrder0.ID)
	if err != nil {
		return nil, err
	}

	riceGrain, err := s.RiceGrainRepository.FindByID(channelID, riceGrainOrder.RiceGrainID)
	if err != nil {
		return nil, err
	}

	producer, err := s.ProducerRepository.FindByID(channelID, riceGrainOrder.SellerID)
	if err != nil {
		return nil, err
	}

	seedOrder, err := s.SeedOrderRepository.FindByRiceGrainOrderID(channelID, riceGrainOrder.ID)
	if err != nil {
		return nil, err
	}

	seed, err := s.SeedRepository.FindByID(channelID, seedOrder.SeedID)
	if err != nil {
		return nil, err
	}

	supplier, err := s.SupplierRepository.FindByID(channelID, seedOrder.SellerID)
	if err != nil {
		return nil, err
	}

	retailerOrder := usermodel.OrderFromRiceOrgModel(nil, riceOrder1.GetReceivedAt(), &retailer.Organization, rice1, riceOrder1.RiceInstance)
	distributorOrder := usermodel.OrderFromRiceOrgModel(riceOrder1.GetShippedAt(), riceOrder0.GetReceivedAt(), &distributor.Organization, rice1, riceOrder1.RiceInstance)
	manufacturerOrder := usermodel.OrderFromRiceOrgModel(riceOrder0.GetShippedAt(), riceGrainOrder.GetReceivedAt(), &manufacturer.Organization, rice0, riceOrder0.RiceInstance)
	producerOrder := usermodel.OrderFromRiceGrainOrgModel(riceGrainOrder.GetShippedAt(), seedOrder.GetReceivedAt(), &producer.Organization, riceGrain, riceGrainOrder.RiceGrainInstance)
	supplierOrder := usermodel.OrderFromSeedOrgModel(seedOrder.GetShippedAt(), nil, &supplier.Organization, seed, seedOrder.SeedInstance)
	traces := []usermodel.Trace{
		{Order: supplierOrder},
		{Order: producerOrder},
		{Order: manufacturerOrder},
		{Order: distributorOrder},
		{Order: retailerOrder},
	}

	userRiceSack := &usermodel.RiceSack{
		Code:   riceSack.Code,
		RiceID: rice0.ID,
		Traces: traces,
	}

	JSON, err := json.Marshal(userRiceSack)
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", JSON)

	return userRiceSack, nil
}
