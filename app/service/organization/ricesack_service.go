package service

import "github.com/meneketehe/hehe/app/model"

type riceSackService struct {
	RiceSackRepository model.RiceSackRepository
}

type RiceSackServiceConfig struct {
	RiceSackRepository model.RiceSackRepository
}

func NewRiceSackService(c *RiceSackServiceConfig) model.RiceSackService {
	return &riceSackService{
		RiceSackRepository: c.RiceSackRepository,
	}
}

func (s *riceSackService) GetAllRiceSack(channelID, stockpileID string) ([]*model.RiceSack, error) {
	return s.RiceSackRepository.FindAll(channelID, stockpileID)
}

func (s *riceSackService) GetAllRiceSackByRiceOrderID(channelID, riceOrderID string) ([]*model.RiceSack, error) {
	return s.RiceSackRepository.FindAllByRiceOrderID(channelID, riceOrderID)
}

func (s *riceSackService) GetRiceSack(channelID, ID string) (*model.RiceSack, error) {
	return s.RiceSackRepository.FindByID(channelID, ID)
}
