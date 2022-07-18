package model

import "time"

type SeedOrder struct {
	ID               string        `json:"id"`
	OrdererID        string        `json:"orderer_id"`
	SellerID         string        `json:"seller_id"`
	SeedID           string        `json:"seed_id"`
	RiceGrainOrderID string        `json:"rice_grain_order_id"`
	Weight           float32       `json:"weight"`
	SeedInstance     *SeedInstance `json:"seed_instance"`
	Order
}

type SeedInstance struct {
	StorageTemperature float32 `json:"storage_temperature"`
	StorageHumidity    float32 `json:"storage_humidity"`
}

func (s *SeedOrder) Ship(shippedAt time.Time, storageTemperature, storageHumidity float32) {
	s.Order.Ship(shippedAt)
	s.SeedInstance = &SeedInstance{
		StorageTemperature: storageTemperature,
		StorageHumidity:    storageHumidity,
	}
}

func (s *SeedInstance) GetStorageTemperature() *float32 {
	if s == nil {
		return nil
	}
	return &s.StorageTemperature
}

func (s *SeedInstance) GetStorageHumidity() *float32 {
	if s == nil {
		return nil
	}
	return &s.StorageHumidity
}

type SeedOrderService interface {
	GetAllOutgoingSeedOrder(channelID, ordererID string) ([]*SeedOrder, error)
	GetAllIncomingSeedOrder(channelID, sellerID string) ([]*SeedOrder, error)
	GetSeedOrderByID(channelID, ID string) (*SeedOrder, error)
	CreateSeedOrder(channelID string, seedOrder *SeedOrder) (*SeedOrder, error)
	AcceptSeedOrder(channelID string, seedOrder *SeedOrder, acceptedAt time.Time) error
	RejectSeedOrder(channelID string, seedOrder *SeedOrder, rejectedAt time.Time, reason string) error
	ShipSeedOrder(channelID string, seedOrder *SeedOrder, shippedAt time.Time, storageTemperature, storageHumidity float32) error
	ReceiveSeedOrder(channelID string, seedOrder *SeedOrder, receivedAt time.Time) error
}

type SeedOrderRepository interface {
	FindAllOutgoing(channelID, ordererID string) ([]*SeedOrder, error)
	FindAllIncoming(channelID, sellerID string) ([]*SeedOrder, error)
	FindByID(channelID, ID string) (*SeedOrder, error)
	FindByRiceGrainOrderID(channelID, riceGrainOrderId string) (*SeedOrder, error)
	Create(channelID string, riceOrder *SeedOrder) error
	Accept(channelID, ID string, acceptedAt time.Time) error
	Reject(channelID, ID string, rejectedAt time.Time, reason string) error
	Ship(channelID, ID string, shippedAt time.Time, storageTemperature, storageHumidity float32) error
	Receive(channelID, ID string, receivedAt time.Time) error
}
