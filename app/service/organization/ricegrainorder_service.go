package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/model"
	"github.com/meneketehe/hehe/app/model/enum"
)

type riceGrainOrderService struct {
	RiceGrainOrderRepository model.RiceGrainOrderRepository
}

type RiceGrainOrderServiceConfig struct {
	RiceGrainOrderRepository model.RiceGrainOrderRepository
}

func NewRiceGrainOrderService(c *RiceGrainOrderServiceConfig) model.RiceGrainOrderService {
	return &riceGrainOrderService{
		RiceGrainOrderRepository: c.RiceGrainOrderRepository,
	}
}

func (s *riceGrainOrderService) GetAllOutgoingRiceGrainOrder(channelID, ordererID string) ([]*model.RiceGrainOrder, error) {
	return s.RiceGrainOrderRepository.FindAllOutgoing(channelID, ordererID)
}

func (s *riceGrainOrderService) GetAllIncomingRiceGrainOrder(channelID, sellerID string) ([]*model.RiceGrainOrder, error) {
	return s.RiceGrainOrderRepository.FindAllIncoming(channelID, sellerID)
}

func (s *riceGrainOrderService) GetAllAcceptedIncomingRiceGrainOrder(channelID, sellerID string) ([]*model.RiceGrainOrder, error) {
	return s.RiceGrainOrderRepository.FindAllAcceptedIncoming(channelID, sellerID)
}

func (s *riceGrainOrderService) GetRiceGrainOrderByID(channelID, ID string) (*model.RiceGrainOrder, error) {
	return s.RiceGrainOrderRepository.FindByID(channelID, ID)
}

func (s *riceGrainOrderService) GetRiceGrainOrderByRiceOrderID(channelID, ID string) (*model.RiceGrainOrder, error) {
	return s.RiceGrainOrderRepository.FindByRiceOrderID(channelID, ID)
}

func (s *riceGrainOrderService) CreateRiceGrainOrder(channelID string, riceGrainOrder *model.RiceGrainOrder) (*model.RiceGrainOrder, error) {
	riceGrainOrder.ID = uuid.New().String()

	if err := s.RiceGrainOrderRepository.Create(channelID, riceGrainOrder); err != nil {
		return nil, err
	}

	return riceGrainOrder, nil
}

func (s *riceGrainOrderService) AcceptRiceGrainOrder(channelID string, riceGrainOrder *model.RiceGrainOrder, acceptedAt time.Time) error {
	if riceGrainOrder.Status != enum.OrderOrdered {
		return fmt.Errorf("you can only accept rice grain order when it's status is %s", enum.OrderOrdered)
	}

	riceGrainOrder.Accept(acceptedAt)
	if err := s.RiceGrainOrderRepository.Accept(channelID, riceGrainOrder.ID, riceGrainOrder.AcceptedAt); err != nil {
		return err
	}

	return nil
}

func (s *riceGrainOrderService) RejectRiceGrainOrder(channelID string, riceGrainOrder *model.RiceGrainOrder, rejectedAt time.Time, reason string) error {
	if riceGrainOrder.Status != enum.OrderOrdered {
		return fmt.Errorf("you can only reject rice grain order when it's status is %s", enum.OrderOrdered)
	}

	riceGrainOrder.Reject(rejectedAt, reason)
	if err := s.RiceGrainOrderRepository.Reject(channelID, riceGrainOrder.ID, riceGrainOrder.RejectedAt, riceGrainOrder.RejectReason); err != nil {
		return err
	}

	return nil
}

func (s *riceGrainOrderService) ShipRiceGrainOrder(channelID string, riceGrainOrder *model.RiceGrainOrder, shippedAt time.Time, plowMethod, sowMethod, irrigation, fertilization string, plantDate, harvestDate time.Time, storageTemperature, storageHumidity float32) error {
	if riceGrainOrder.Status != enum.OrderAvailable {
		return fmt.Errorf("you can only ship rice grain order when it is %s", enum.OrderAvailable)
	}

	riceGrainOrder.Ship(shippedAt, plowMethod, sowMethod, irrigation, fertilization, plantDate, harvestDate, storageTemperature, storageHumidity)
	if err := s.RiceGrainOrderRepository.Ship(
		channelID,
		riceGrainOrder.ID,
		riceGrainOrder.ShippedAt,
		riceGrainOrder.RiceGrainInstance.PlowMethod,
		riceGrainOrder.RiceGrainInstance.SowMethod,
		riceGrainOrder.RiceGrainInstance.Irrigation,
		riceGrainOrder.RiceGrainInstance.Fertilization,
		riceGrainOrder.RiceGrainInstance.PlantDate,
		riceGrainOrder.RiceGrainInstance.HarvestDate,
		riceGrainOrder.RiceGrainInstance.StorageTemperature,
		riceGrainOrder.RiceGrainInstance.StorageHumidity,
	); err != nil {
		return err
	}

	return nil
}

func (s *riceGrainOrderService) ReceiveRiceGrainOrder(channelID string, riceGrainOrder *model.RiceGrainOrder, receivedAt time.Time) error {
	if riceGrainOrder.Status != enum.OrderShipped {
		return fmt.Errorf("you can only receive rice grain order when it is %s", enum.OrderShipped)
	}

	riceGrainOrder.Receive(receivedAt)
	if err := s.RiceGrainOrderRepository.Receive(channelID, riceGrainOrder.ID, riceGrainOrder.ReceivedAt); err != nil {
		return err
	}

	return nil
}
