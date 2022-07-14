package model

import (
	"time"
)

type RiceOrder struct {
	ID           string        `json:"id"`
	OrdererID    string        `json:"orderer_id"`
	SellerID     string        `json:"seller_id"`
	RiceID       string        `json:"rice_id"`
	Quantity     int32         `json:"quantity"`
	RiceInstance *RiceInstance `json:"rice_instance"`
	Order
}

type RiceInstance struct {
	Grade              string    `json:"grade"`
	MillingDate        time.Time `json:"milling_date"`
	StorageTemperature float32   `json:"storage_temperature"`
	StorageHumidity    float32   `json:"storage_humidity"`
}

func (r *RiceOrder) Ship(shipAt time.Time, grade string, millingDate time.Time, storageTemperature float32, storageHumidity float32) {
	r.Order.Ship(shipAt)
	r.RiceInstance = &RiceInstance{
		Grade:              grade,
		MillingDate:        millingDate,
		StorageTemperature: storageTemperature,
		StorageHumidity:    storageHumidity,
	}
}

func (r *RiceInstance) GetGrade() *string {
	if r == nil {
		return nil
	}
	return &r.Grade
}

func (r *RiceInstance) GetMillingDate() *time.Time {
	if r == nil {
		return nil
	}
	return &r.MillingDate
}

func (r *RiceInstance) GetStorageTemperature() *float32 {
	if r == nil {
		return nil
	}
	return &r.StorageTemperature
}

func (r *RiceInstance) GetStorageHumidity() *float32 {
	if r == nil {
		return nil
	}
	return &r.StorageHumidity
}

type RiceOrderService interface {
	GetAllOutgoingRiceOrder(channelID, ordererID string) ([]*RiceOrder, error)
	GetAllIncomingRiceOrder(channelID, sellerID string) ([]*RiceOrder, error)
	GetAllAcceptedIncomingRiceOrder(channelID, sellerID string) ([]*RiceOrder, error)
	GetRiceOrderByID(channelID, ID string) (*RiceOrder, error)
	CreateRiceOrder(channelID string, riceOrder *RiceOrder) (*RiceOrder, error)
	AcceptRiceOrder(channelID string, riceOrder *RiceOrder, acceptedAt time.Time) error
	RejectRiceOrder(channelID string, riceOrder *RiceOrder, rejectedAt time.Time, reason string) error
	ShipRiceOrder(channelID string, riceOrder *RiceOrder, shippedAt time.Time, grade string, millingDate time.Time, storageTemperature float32, storageHumidity float32) error
	ReceiveRiceOrder(channelID string, riceOrder *RiceOrder, receivedAt time.Time) error
}

type RiceOrderRepository interface {
	FindAllOutgoing(channelID, ordererID string) ([]*RiceOrder, error)
	FindAllIncoming(channelID, sellerID string) ([]*RiceOrder, error)
	FindAllAcceptedIncoming(channelID, sellerID string) ([]*RiceOrder, error)
	FindByID(channelID, ID string) (*RiceOrder, error)
	Create(channelID string, riceOrder *RiceOrder) error
	Accept(channelID, ID string, acceptedAt time.Time) error
	Reject(channelID, ID string, rejectedAt time.Time, reason string) error
	Ship(channelID, ID string, shippedAt time.Time, grade string, millingDate time.Time, storageTemperature float32, storageHumidity float32) error
	Receive(channelID, ID string, receivedAt time.Time) error
}
