package response

import (
	"time"

	"github.com/meneketehe/hehe/app/model"
)

type seedOrdersResponse struct {
	ID                string     `json:"id"`
	Weight            float32    `json:"weight"`
	Status            string     `json:"status"`
	LastTransactionAt *time.Time `json:"last_transaction_at"`
}

type seedOrderResponse struct {
	ID                 string     `json:"id"`
	OrdererID          string     `json:"orderer_id"`
	SellerID           string     `json:"seller_id"`
	SeedID             string     `json:"seed_id"`
	RiceGrainOrderID   string     `json:"rice_grain_order_id"`
	Weight             float32    `json:"weight"`
	Status             string     `json:"status"`
	StorageTemperature *float32   `json:"storage_temperature"`
	StorageHumidity    *float32   `json:"storage_humidity"`
	OrderedAt          *time.Time `json:"ordered_at"`
	AcceptedAt         *time.Time `json:"accepted_at"`
	RejectedAt         *time.Time `json:"rejected_at"`
	RejectedReason     string     `json:"rejected_reason"`
	ProcessingAt       *time.Time `json:"processing_at"`
	AvailableAt        *time.Time `json:"available_at"`
	ShippedAt          *time.Time `json:"shipped_at"`
	ReceivedAt         *time.Time `json:"received_at"`
}

func SeedOrdersResponse(seedOrders []*model.SeedOrder) []*seedOrdersResponse {
	res := make([]*seedOrdersResponse, 0)
	for _, seedOrder := range seedOrders {
		res = append(res, &seedOrdersResponse{
			ID:                seedOrder.ID,
			Weight:            seedOrder.Weight,
			Status:            seedOrder.Status,
			LastTransactionAt: seedOrder.GetLastTransactionAt(),
		})
	}

	return res
}

func SeedOrderResponse(seedOrder *model.SeedOrder) *seedOrderResponse {
	return &seedOrderResponse{
		ID:                 seedOrder.ID,
		OrdererID:          seedOrder.OrdererID,
		SellerID:           seedOrder.SellerID,
		SeedID:             seedOrder.SeedID,
		RiceGrainOrderID:   seedOrder.RiceGrainOrderID,
		Weight:             seedOrder.Weight,
		Status:             seedOrder.Status,
		StorageTemperature: seedOrder.SeedInstance.GetStorageTemperature(),
		StorageHumidity:    seedOrder.SeedInstance.GetStorageHumidity(),
		OrderedAt:          seedOrder.GetOrderedAt(),
		AcceptedAt:         seedOrder.GetAcceptedAt(),
		RejectedAt:         seedOrder.GetRejectedAt(),
		RejectedReason:     seedOrder.RejectReason,
		ProcessingAt:       seedOrder.GetProcessingAt(),
		AvailableAt:        seedOrder.GetAvailableAt(),
		ShippedAt:          seedOrder.GetShippedAt(),
		ReceivedAt:         seedOrder.GetReceivedAt(),
	}
}
