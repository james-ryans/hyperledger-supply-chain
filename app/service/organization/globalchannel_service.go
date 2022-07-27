package service

import (
	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/model"
)

type globalChannelService struct {
	GlobalChannelRepository model.GlobalChannelRepository
}

type GlobalChannelServiceConfig struct {
	GlobalChannelRepository model.GlobalChannelRepository
}

func NewGlobalChannelService(c *GlobalChannelServiceConfig) model.GlobalChannelService {
	return &globalChannelService{
		GlobalChannelRepository: c.GlobalChannelRepository,
	}
}

func (s *globalChannelService) GetAllChannels() ([]*model.GlobalChannel, error) {
	chs, err := s.GlobalChannelRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return chs, nil
}

func (s *globalChannelService) GetChannel(ID string) (*model.GlobalChannel, error) {
	ch, err := s.GlobalChannelRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (s *globalChannelService) CheckNameExists(name string) (bool, error) {
	ch, err := s.GlobalChannelRepository.FindByName(name)
	if err != nil {
		return false, err
	}

	return ch != nil, nil
}

func (s *globalChannelService) CreateChannel(ch *model.GlobalChannel) (*model.GlobalChannel, error) {
	ch.ID = uuid.New().String()

	ch, err := s.GlobalChannelRepository.Create(ch)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
