package service

import "github.com/meneketehe/hehe/app/model"

type retailerService struct {
	RetailerRepository model.RetailerRepository
}

type RetailerServiceConfig struct {
	RetailerRepository model.RetailerRepository
}

func NewRetailerService(c *RetailerServiceConfig) model.RetailerService {
	return &retailerService{
		RetailerRepository: c.RetailerRepository,
	}
}

func (s *retailerService) GetAllRetailers(channelID string) ([]*model.Retailer, error) {
	return s.RetailerRepository.FindAll(channelID)
}

func (s *retailerService) GetRetailerByID(channelID, ID string) (*model.Retailer, error) {
	return s.RetailerRepository.FindByID(channelID, ID)
}

func (s *retailerService) CreateRetailer(channelID string, retailer *model.Retailer) (*model.Retailer, error) {
	if err := s.RetailerRepository.Create(channelID, retailer); err != nil {
		return nil, err
	}

	return retailer, nil
}
