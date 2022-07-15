package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/model"
	"github.com/meneketehe/hehe/app/model/enum"
)

type riceOrderService struct {
	RiceOrderRepository model.RiceOrderRepository
}

type RiceOrderServiceConfig struct {
	RiceOrderRepository model.RiceOrderRepository
}

func NewRiceOrderService(c *RiceOrderServiceConfig) model.RiceOrderService {
	return &riceOrderService{
		RiceOrderRepository: c.RiceOrderRepository,
	}
}

func (s *riceOrderService) GetAllOutgoingRiceOrder(channelID, ordererID string) ([]*model.RiceOrder, error) {
	return s.RiceOrderRepository.FindAllOutgoing(channelID, ordererID)
}

func (s *riceOrderService) GetAllIncomingRiceOrder(channelID, sellerID string) ([]*model.RiceOrder, error) {
	return s.RiceOrderRepository.FindAllIncoming(channelID, sellerID)
}

func (s *riceOrderService) GetAllAcceptedIncomingRiceOrder(channelID, sellerID string) ([]*model.RiceOrder, error) {
	return s.RiceOrderRepository.FindAllAcceptedIncoming(channelID, sellerID)
}

func (s *riceOrderService) GetRiceOrderByID(channelID, ID string) (*model.RiceOrder, error) {
	return s.RiceOrderRepository.FindByID(channelID, ID)
}

func (s *riceOrderService) CreateRiceOrder(channelID string, riceOrder *model.RiceOrder) (*model.RiceOrder, error) {
	riceOrder.ID = uuid.New().String()

	if err := s.RiceOrderRepository.Create(channelID, riceOrder); err != nil {
		return nil, err
	}

	return riceOrder, nil
}

func (s *riceOrderService) AcceptRiceOrder(channelID string, riceOrder *model.RiceOrder, acceptedAt time.Time) error {
	if riceOrder.Status != enum.OrderOrdered {
		return fmt.Errorf("you can only accept rice order when it's status is %s", enum.OrderOrdered)
	}

	riceOrder.Accept(acceptedAt)
	if err := s.RiceOrderRepository.Accept(channelID, riceOrder.ID, riceOrder.AcceptedAt); err != nil {
		return err
	}

	return nil
}

func (s *riceOrderService) AcceptDistributionRiceOrder(channelID string, riceOrder *model.RiceOrder, acceptedAt time.Time) error {
	if riceOrder.Status != enum.OrderOrdered {
		return fmt.Errorf("you can only accept rice order when it's status is %s", enum.OrderOrdered)
	}

	riceOrder.Accept(acceptedAt)
	if err := s.RiceOrderRepository.AcceptDistribution(channelID, riceOrder.ID, riceOrder.AcceptedAt); err != nil {
		return err
	}

	return nil
}

func (s *riceOrderService) RejectRiceOrder(channelID string, riceOrder *model.RiceOrder, rejectedAt time.Time, reason string) error {
	if riceOrder.Status != enum.OrderOrdered {
		return fmt.Errorf("you can only reject rice order when it's status is %s", enum.OrderOrdered)
	}

	riceOrder.Reject(rejectedAt, reason)
	if err := s.RiceOrderRepository.Reject(channelID, riceOrder.ID, riceOrder.RejectedAt, riceOrder.RejectReason); err != nil {
		return err
	}

	return nil
}

func (s *riceOrderService) ShipRiceOrder(channelID string, riceOrder *model.RiceOrder, shippedAt time.Time, grade string, millingDate time.Time, storageTemperature float32, storageHumidity float32) error {
	if riceOrder.Status != enum.OrderAvailable {
		return fmt.Errorf("you can only ship rice order when it is %s", enum.OrderAvailable)
	}

	riceOrder.Ship(shippedAt, grade, millingDate, storageTemperature, storageHumidity)
	if err := s.RiceOrderRepository.Ship(
		channelID,
		riceOrder.ID,
		riceOrder.ShippedAt,
		riceOrder.RiceInstance.Grade,
		riceOrder.RiceInstance.MillingDate,
		riceOrder.RiceInstance.StorageTemperature,
		riceOrder.RiceInstance.StorageHumidity,
	); err != nil {
		return err
	}

	return nil
}

func (s *riceOrderService) ShipDistributionRiceOrder(channelID string, riceOrder *model.RiceOrder, shippedAt time.Time, grade string, millingDate time.Time, storageTemperature float32, storageHumidity float32) error {
	if riceOrder.Status != enum.OrderAvailable {
		return fmt.Errorf("you can only ship rice order when it is %s", enum.OrderAvailable)
	}

	riceOrder.Ship(shippedAt, grade, millingDate, storageTemperature, storageHumidity)
	if err := s.RiceOrderRepository.ShipDistribution(
		channelID,
		riceOrder.ID,
		riceOrder.ShippedAt,
		riceOrder.RiceInstance.Grade,
		riceOrder.RiceInstance.MillingDate,
		riceOrder.RiceInstance.StorageTemperature,
		riceOrder.RiceInstance.StorageHumidity,
	); err != nil {
		return err
	}

	return nil
}

func (s *riceOrderService) ReceiveRiceOrder(channelID string, riceOrder *model.RiceOrder, receivedAt time.Time) error {
	if riceOrder.Status != enum.OrderShipped {
		return fmt.Errorf("you can only ship rice order when it is %s", enum.OrderShipped)
	}

	riceOrder.Receive(receivedAt)
	if err := s.RiceOrderRepository.Receive(channelID, riceOrder.ID, riceOrder.ReceivedAt); err != nil {
		return err
	}

	return nil
}

func (s *riceOrderService) ReceiveDistributionRiceOrder(channelID string, riceOrder *model.RiceOrder, receivedAt time.Time) error {
	if riceOrder.Status != enum.OrderShipped {
		return fmt.Errorf("you can only ship rice order when it is %s", enum.OrderShipped)
	}

	riceOrder.Receive(receivedAt)
	if err := s.RiceOrderRepository.ReceiveDistribution(channelID, riceOrder.ID, riceOrder.ReceivedAt); err != nil {
		return err
	}

	return nil
}
