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

func (s *seedOrderService) GetSeedOrderByRiceGrainID(channelID, ID string) (*model.SeedOrder, error) {
	return s.SeedOrderRepository.FindByRiceGrainOrderID(channelID, ID)
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
	seedOrder.Process(acceptedAt)
	seedOrder.Available(acceptedAt)
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

func (s *seedOrderService) ShipSeedOrder(channelID string, seedOrder *model.SeedOrder, shippedAt time.Time, storageTemperature, storageHumidity float32) error {
	if seedOrder.Status != enum.OrderAvailable {
		return fmt.Errorf("you can only ship seed order when it is %s", enum.OrderAvailable)
	}

	seedOrder.Ship(shippedAt, storageTemperature, storageHumidity)
	if err := s.SeedOrderRepository.Ship(channelID, seedOrder.ID, seedOrder.ShippedAt, seedOrder.SeedInstance.StorageTemperature, seedOrder.SeedInstance.StorageHumidity); err != nil {
		return err
	}

	return nil
}

func (s *seedOrderService) ReceiveSeedOrder(channelID string, seedOrder *model.SeedOrder, receivedAt time.Time) error {
	if seedOrder.Status != enum.OrderShipped {
		return fmt.Errorf("you can only receive seed order when it is %s", enum.OrderShipped)
	}

	seedOrder.Receive(receivedAt)
	if err := s.SeedOrderRepository.Receive(channelID, seedOrder.ID, seedOrder.ReceivedAt); err != nil {
		return err
	}

	return nil
}
