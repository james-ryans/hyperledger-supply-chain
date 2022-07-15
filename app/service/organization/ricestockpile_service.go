package service

import "github.com/meneketehe/hehe/app/model"

type riceStockpileService struct {
	RiceStockpileRepository model.RiceStockpileRepository
}

type RiceStockpileServiceConfig struct {
	RiceStockpileRepository model.RiceStockpileRepository
}

func NewRiceStockpileService(c *RiceStockpileServiceConfig) model.RiceStockpileService {
	return &riceStockpileService{
		RiceStockpileRepository: c.RiceStockpileRepository,
	}
}

func (s *riceStockpileService) GetAllRiceStockpile(channelID, vendorID string) ([]*model.RiceStockpile, error) {
	return s.RiceStockpileRepository.FindAll(channelID, vendorID)
}

func (s *riceStockpileService) GetRiceStockpileByID(channelID, ID string) (*model.RiceStockpile, error) {
	return s.RiceStockpileRepository.FindByID(channelID, ID)
}
