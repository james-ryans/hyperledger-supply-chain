package service

import (
	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/model"
)

type riceService struct {
	RiceRepository model.RiceRepository
}

type RiceServiceConfig struct {
	RiceRepository model.RiceRepository
}

func NewRiceService(c *RiceServiceConfig) model.RiceService {
	return &riceService{
		RiceRepository: c.RiceRepository,
	}
}

func (s *riceService) GetAllRices(channelID, ID string) (*[]model.Rice, error) {
	return s.RiceRepository.FindAll(channelID, ID)
}

func (s *riceService) GetRiceByID(channelID, ID string) (*model.Rice, error) {
	return s.RiceRepository.FindByID(channelID, ID)
}

func (s *riceService) CreateRice(channelID string, rice *model.Rice) (*model.Rice, error) {
	rice.ID = uuid.New().String()

	if err := s.RiceRepository.Create(channelID, rice); err != nil {
		return nil, err
	}

	return rice, nil
}

func (s *riceService) UpdateRice(channelID string, rice *model.Rice) error {
	return s.RiceRepository.Update(channelID, rice)
}

func (s *riceService) DeleteRice(channelID, ID string) error {
	return s.RiceRepository.Delete(channelID, ID)
}
