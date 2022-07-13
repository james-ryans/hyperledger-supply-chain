package service

import "github.com/meneketehe/hehe/app/model"

type distributorService struct {
	DistributorRepository model.DistributorRepository
}

type DistributorServiceConfig struct {
	DistributorRepository model.DistributorRepository
}

func NewDistributorService(c *DistributorServiceConfig) model.DistributorService {
	return &distributorService{
		DistributorRepository: c.DistributorRepository,
	}
}

func (s *distributorService) GetAllDistributors(channelID string) ([]*model.Distributor, error) {
	return s.DistributorRepository.FindAll(channelID)
}

func (s *distributorService) GetDistributorByID(channelID, ID string) (*model.Distributor, error) {
	return s.DistributorRepository.FindByID(channelID, ID)
}

func (s *distributorService) CreateDistributor(channelID string, distributor *model.Distributor) (*model.Distributor, error) {
	if err := s.DistributorRepository.Create(channelID, distributor); err != nil {
		return nil, err
	}

	return distributor, nil
}
