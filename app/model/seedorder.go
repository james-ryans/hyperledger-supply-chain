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

func (r *SeedInstance) GetStorageTemperature() *float32 {
	if r == nil {
		return nil
	}
	return &r.StorageTemperature
}

func (r *SeedInstance) GetStorageHumidity() *float32 {
	if r == nil {
		return nil
	}
	return &r.StorageHumidity
}

type SeedOrderService interface {
	GetAllOutgoingSeedOrder(channelID, ordererID string) ([]*SeedOrder, error)
	GetAllIncomingSeedOrder(channelID, sellerID string) ([]*SeedOrder, error)
	GetSeedOrderByID(channelID, ID string) (*SeedOrder, error)
	CreateSeedOrder(channelID string, seedOrder *SeedOrder) (*SeedOrder, error)
	AcceptSeedOrder(channelID string, seedOrder *SeedOrder, acceptedAt time.Time) error
	RejectSeedOrder(channelID string, seedOrder *SeedOrder, rejectedAt time.Time, reason string) error
	// ShipSeedOrder(channelID string, seedOrder *SeedOrder, shippedAt time.Time) error
	// ReceiveSeedOrder(channelID string, seedOrder *SeedOrder, receivedAt time.Time) error
}

type SeedOrderRepository interface {
	FindAllOutgoing(channelID, ordererID string) ([]*SeedOrder, error)
	FindAllIncoming(channelID, sellerID string) ([]*SeedOrder, error)
	FindByID(channelID, ID string) (*SeedOrder, error)
	Create(channelID string, riceOrder *SeedOrder) error
	Accept(channelID, ID string, acceptedAt time.Time) error
	Reject(channelID, ID string, rejectedAt time.Time, reason string) error
}
