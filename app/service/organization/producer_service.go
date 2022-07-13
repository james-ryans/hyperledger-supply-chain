package service

import "github.com/meneketehe/hehe/app/model"

type producerService struct {
	ProducerRepository model.ProducerRepository
}

type ProducerServiceConfig struct {
	ProducerRepository model.ProducerRepository
}

func NewProducerService(c *ProducerServiceConfig) model.ProducerService {
	return &producerService{
		ProducerRepository: c.ProducerRepository,
	}
}

func (s *producerService) GetAllProducers(channelID string) ([]*model.Producer, error) {
	return s.ProducerRepository.FindAll(channelID)
}

func (s *producerService) GetProducerByID(channelID, ID string) (*model.Producer, error) {
	return s.ProducerRepository.FindByID(channelID, ID)
}

func (s *producerService) CreateProducer(channelID string, producer *model.Producer) (*model.Producer, error) {
	if err := s.ProducerRepository.Create(channelID, producer); err != nil {
		return nil, err
	}

	return producer, nil
}
