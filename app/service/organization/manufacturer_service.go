package service

import "github.com/meneketehe/hehe/app/model"

type manufacturerService struct {
	ManufacturerRepository model.ManufacturerRepository
}

type ManufacturerServiceConfig struct {
	ManufacturerRepository model.ManufacturerRepository
}

func NewManufacturerService(c *ManufacturerServiceConfig) model.ManufacturerService {
	return &manufacturerService{
		ManufacturerRepository: c.ManufacturerRepository,
	}
}

func (s *manufacturerService) GetAllManufacturers(channelID string) ([]*model.Manufacturer, error) {
	return s.ManufacturerRepository.FindAll(channelID)
}

func (s *manufacturerService) GetManufacturerByID(channelID, ID string) (*model.Manufacturer, error) {
	return s.ManufacturerRepository.FindByID(channelID, ID)
}

func (s *manufacturerService) CreateManufacturer(channelID string, manufacturer *model.Manufacturer) (*model.Manufacturer, error) {
	if err := s.ManufacturerRepository.Create(channelID, manufacturer); err != nil {
		return nil, err
	}

	return manufacturer, nil
}
