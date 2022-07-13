package model

import (
	"time"

	"github.com/meneketehe/hehe/app/model/enum"
)

type Order struct {
	OrderedAt    time.Time `json:"ordered_at"`
	AcceptedAt   time.Time `json:"accepted_at"`
	RejectedAt   time.Time `json:"rejected_at"`
	RejectReason string    `json:"reject_reason"`
	ProcessingAt time.Time `json:"processing_at"`
	AvailableAt  time.Time `json:"available_at"`
	ShippedAt    time.Time `json:"shipped_at"`
	ReceivedAt   time.Time `json:"received_at"`
	Status       string    `json:"status"`
}

func NewOrder(at time.Time) Order {
	return Order{
		OrderedAt: at,
		Status:    enum.OrderOrdered,
	}
}

func (order *Order) GetOrderedAt() *time.Time {
	if order.OrderedAt.IsZero() {
		return nil
	}
	return &order.OrderedAt
}

func (order *Order) GetAcceptedAt() *time.Time {
	if order.AcceptedAt.IsZero() {
		return nil
	}
	return &order.AcceptedAt
}

func (order *Order) GetRejectedAt() *time.Time {
	if order.RejectedAt.IsZero() {
		return nil
	}
	return &order.RejectedAt
}

func (order *Order) GetProcessingAt() *time.Time {
	if order.ProcessingAt.IsZero() {
		return nil
	}
	return &order.ProcessingAt
}

func (order *Order) GetAvailableAt() *time.Time {
	if order.AvailableAt.IsZero() {
		return nil
	}
	return &order.AvailableAt
}

func (order *Order) GetShippedAt() *time.Time {
	if order.ShippedAt.IsZero() {
		return nil
	}
	return &order.ShippedAt
}

func (order *Order) GetReceivedAt() *time.Time {
	if order.ReceivedAt.IsZero() {
		return nil
	}
	return &order.ReceivedAt
}

func (order *Order) GetLastTransactionAt() *time.Time {
	switch order.Status {
	case enum.OrderOrdered:
		return order.GetOrderedAt()
	case enum.OrderAccepted:
		return order.GetAcceptedAt()
	case enum.OrderRejected:
		return order.GetRejectedAt()
	case enum.OrderProcessing:
		return order.GetProcessingAt()
	case enum.OrderAvailable:
		return order.GetAvailableAt()
	case enum.OrderShipped:
		return order.GetShippedAt()
	case enum.OrderReceived:
		return order.GetReceivedAt()
	default:
		return nil
	}
}

func (order *Order) Accept(at time.Time) {
	order.AcceptedAt = at
	order.Status = enum.OrderAccepted
}

func (order *Order) Reject(at time.Time, reason string) {
	order.RejectedAt = at
	order.RejectReason = reason
	order.Status = enum.OrderRejected
}

func (order *Order) Process(at time.Time) {
	order.ProcessingAt = at
	order.Status = enum.OrderProcessing
}

func (order *Order) Available(at time.Time) {
	order.AvailableAt = at
	order.Status = enum.OrderAvailable
}

func (order *Order) Ship(at time.Time) {
	order.ShippedAt = at
	order.Status = enum.OrderShipped
}

func (order *Order) Receive(at time.Time) {
	order.ReceivedAt = at
	order.Status = enum.OrderReceived
}
