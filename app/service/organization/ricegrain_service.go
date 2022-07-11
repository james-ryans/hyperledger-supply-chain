package service

import (
	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/model"
)

type riceGrainService struct {
	RiceGrainRepository model.RiceGrainRepository
}

type RiceGrainServiceConfig struct {
	RiceGrainRepository model.RiceGrainRepository
}

func NewRiceGrainService(c *RiceGrainServiceConfig) model.RiceGrainService {
	return &riceGrainService{
		RiceGrainRepository: c.RiceGrainRepository,
	}
}

func (s *riceGrainService) GetAllRiceGrains(channelID, ID string) (*[]model.RiceGrain, error) {
	return s.RiceGrainRepository.FindAll(channelID, ID)
}

func (s *riceGrainService) GetRiceGrainByID(channelID, ID string) (*model.RiceGrain, error) {
	return s.RiceGrainRepository.FindByID(channelID, ID)
}

func (s *riceGrainService) CreateRiceGrain(channelID string, riceGrain *model.RiceGrain) (*model.RiceGrain, error) {
	riceGrain.ID = uuid.New().String()

	if err := s.RiceGrainRepository.Create(channelID, riceGrain); err != nil {
		return nil, err
	}

	return riceGrain, nil
}

func (s *riceGrainService) UpdateRiceGrain(channelID string, riceGrain *model.RiceGrain) error {
	return s.RiceGrainRepository.Update(channelID, riceGrain)
}

func (s *riceGrainService) DeleteRiceGrain(channelID, ID string) error {
	return s.RiceGrainRepository.Delete(channelID, ID)
}
