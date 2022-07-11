package service

import (
	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/model"
)

type seedService struct {
	SeedRepository model.SeedRepository
}

type SeedServiceConfig struct {
	SeedRepository model.SeedRepository
}

func NewSeedService(c *SeedServiceConfig) model.SeedService {
	return &seedService{
		SeedRepository: c.SeedRepository,
	}
}

func (s *seedService) GetAllSeeds(channelID, ID string) (*[]model.Seed, error) {
	return s.SeedRepository.FindAll(channelID, ID)
}

func (s *seedService) GetSeedByID(channelID, ID string) (*model.Seed, error) {
	return s.SeedRepository.FindByID(channelID, ID)
}

func (s *seedService) CreateSeed(channelID string, seed *model.Seed) (*model.Seed, error) {
	seed.ID = uuid.New().String()

	if err := s.SeedRepository.Create(channelID, seed); err != nil {
		return nil, err
	}

	return seed, nil
}

func (s *seedService) UpdateSeed(channelID string, seed *model.Seed) error {
	return s.SeedRepository.Update(channelID, seed)
}

func (s *seedService) DeleteSeed(channelID, ID string) error {
	return s.SeedRepository.Delete(channelID, ID)
}
