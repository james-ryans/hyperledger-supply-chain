package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/model"
	"github.com/meneketehe/hehe/app/model/enum"
)

type seedOrderService struct {
	SeedOrderRepository model.SeedOrderRepository
}

type SeedOrderServiceConfig struct {
	SeedOrderRepository model.SeedOrderRepository
}

func NewSeedOrderService(c *SeedOrderServiceConfig) model.SeedOrderService {
	return &seedOrderService{
		SeedOrderRepository: c.SeedOrderRepository,
	}
}

func (s *seedOrderService) GetAllOutgoingSeedOrder(channelID, ordererID string) ([]*model.SeedOrder, error) {
	return s.SeedOrderRepository.FindAllOutgoing(channelID, ordererID)
}

func (s *seedOrderService) GetAllIncomingSeedOrder(channelID, sellerID string) ([]*model.SeedOrder, error) {
	return s.SeedOrderRepository.FindAllIncoming(channelID, sellerID)
}

func (s *seedOrderService) GetSeedOrderByID(channelID, ID string) (*model.SeedOrder, error) {
	return s.SeedOrderRepository.FindByID(channelID, ID)
}

func (s *seedOrderService) CreateSeedOrder(channelID string, seedOrder *model.SeedOrder) (*model.SeedOrder, error) {
	seedOrder.ID = uuid.New().String()

	if err := s.SeedOrderRepository.Create(channelID, seedOrder); err != nil {
		return nil, err
	}

	return seedOrder, nil
}

func (s *seedOrderService) AcceptSeedOrder(channelID string, seedOrder *model.SeedOrder, acceptedAt time.Time) error {
	if seedOrder.Status != enum.OrderOrdered {
		return fmt.Errorf("you can only accept seed order when it's status is %s", enum.OrderOrdered)
	}

	seedOrder.Accept(acceptedAt)
	if err := s.SeedOrderRepository.Accept(channelID, seedOrder.ID, seedOrder.AcceptedAt); err != nil {
		return err
	}

	return nil
}

func (s *seedOrderService) RejectSeedOrder(channelID string, seedOrder *model.SeedOrder, rejectedAt time.Time, reason string) error {
	if seedOrder.Status != enum.OrderOrdered {
		return fmt.Errorf("you can only reject seed order when it's status is %s", enum.OrderOrdered)
	}

	seedOrder.Reject(rejectedAt, reason)
	if err := s.SeedOrderRepository.Reject(channelID, seedOrder.ID, seedOrder.RejectedAt, seedOrder.RejectReason); err != nil {
		return err
	}

	return nil
}
