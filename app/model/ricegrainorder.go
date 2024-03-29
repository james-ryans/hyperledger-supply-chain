package model

import (
	"time"
)

type RiceGrainOrder struct {
	ID                string             `json:"id"`
	OrdererID         string             `json:"orderer_id"`
	SellerID          string             `json:"seller_id"`
	RiceGrainID       string             `json:"rice_grain_id"`
	RiceOrderID       string             `json:"rice_order_id"`
	Weight            float32            `json:"weight"`
	RiceGrainInstance *RiceGrainInstance `json:"rice_grain_instance"`
	Order
}

type RiceGrainInstance struct {
	PlowMethod         string    `json:"plow_method"`
	SowMethod          string    `json:"sow_method"`
	Irrigation         string    `json:"irrigation"`
	Fertilization      string    `json:"fertilization"`
	PlantDate          time.Time `json:"plant_date"`
	HarvestDate        time.Time `json:"harvest_date"`
	StorageTemperature float32   `json:"storage_temperature"`
	StorageHumidity    float32   `json:"storage_humidity"`
}

func (r *RiceGrainOrder) Ship(shippedAt time.Time, plowMethod, sowMethod, irrigation, fertilization string, plantDate, harvestDate time.Time, storageTemperature, storageHumidity float32) {
	r.Order.Ship(shippedAt)
	r.RiceGrainInstance = &RiceGrainInstance{
		PlowMethod:         plowMethod,
		SowMethod:          sowMethod,
		Irrigation:         irrigation,
		Fertilization:      fertilization,
		PlantDate:          plantDate,
		HarvestDate:        harvestDate,
		StorageTemperature: storageTemperature,
		StorageHumidity:    storageHumidity,
	}
}

func (r *RiceGrainInstance) GetPlowMethod() *string {
	if r == nil {
		return nil
	}
	return &r.PlowMethod
}

func (r *RiceGrainInstance) GetSowMethod() *string {
	if r == nil {
		return nil
	}
	return &r.SowMethod
}

func (r *RiceGrainInstance) GetIrrigation() *string {
	if r == nil {
		return nil
	}
	return &r.Irrigation
}

func (r *RiceGrainInstance) GetFertilization() *string {
	if r == nil {
		return nil
	}
	return &r.Fertilization
}

func (r *RiceGrainInstance) GetPlantDate() *time.Time {
	if r == nil {
		return nil
	}
	return &r.PlantDate
}

func (r *RiceGrainInstance) GetHarvestDate() *time.Time {
	if r == nil {
		return nil
	}
	return &r.HarvestDate
}

func (r *RiceGrainInstance) GetStorageTemperature() *float32 {
	if r == nil {
		return nil
	}
	return &r.StorageTemperature
}

func (r *RiceGrainInstance) GetStorageHumidity() *float32 {
	if r == nil {
		return nil
	}
	return &r.StorageHumidity
}

type RiceGrainOrderService interface {
	GetAllOutgoingRiceGrainOrder(channelID, ordererID string) ([]*RiceGrainOrder, error)
	GetAllIncomingRiceGrainOrder(channelID, sellerID string) ([]*RiceGrainOrder, error)
	GetAllAcceptedIncomingRiceGrainOrder(channelID, sellerID string) ([]*RiceGrainOrder, error)
	GetRiceGrainOrderByID(channelID, ID string) (*RiceGrainOrder, error)
	GetRiceGrainOrderByRiceOrderID(channelID, ID string) (*RiceGrainOrder, error)
	CreateRiceGrainOrder(channelID string, riceOrder *RiceGrainOrder) (*RiceGrainOrder, error)
	AcceptRiceGrainOrder(channelID string, riceOrder *RiceGrainOrder, acceptedAt time.Time) error
	RejectRiceGrainOrder(channelID string, riceOrder *RiceGrainOrder, rejectedAt time.Time, reason string) error
	ShipRiceGrainOrder(channelID string, riceOrder *RiceGrainOrder, shippedAt time.Time, plowMethod, sowMethod, irrigation, fertilization string, plantDate, harvestDate time.Time, storageTemperature, storageHumidity float32) error
	ReceiveRiceGrainOrder(channelID string, riceOrder *RiceGrainOrder, receivedAt time.Time) error
}

type RiceGrainOrderRepository interface {
	FindAllOutgoing(channelID, ordererID string) ([]*RiceGrainOrder, error)
	FindAllIncoming(channelID, sellerID string) ([]*RiceGrainOrder, error)
	FindAllAcceptedIncoming(channelID, sellerID string) ([]*RiceGrainOrder, error)
	FindByID(channelID, ID string) (*RiceGrainOrder, error)
	FindByRiceOrderID(channelID, riceOrderID string) (*RiceGrainOrder, error)
	Create(channelID string, riceOrder *RiceGrainOrder) error
	Accept(channelID, ID string, acceptedAt time.Time) error
	Reject(channelID, ID string, rejectedAt time.Time, reason string) error
	Ship(channelID, ID string, shippedAt time.Time, plowMethod, sowMethod, irrigation, fertilization string, plantDate, harvestDate time.Time, storageTemperature, storageHumidity float32) error
	Receive(channelID, ID string, receivedAt time.Time) error
}
