package response

import (
	"time"

	"github.com/meneketehe/hehe/app/model"
)

type riceOrdersResponse struct {
	ID                string     `json:"id"`
	Quantity          int32      `json:"quantity"`
	Status            string     `json:"status"`
	LastTransactionAt *time.Time `json:"last_transaction_at"`
}

type riceOrderResponse struct {
	ID                 string     `json:"id"`
	SellerID           string     `json:"seller_id"`
	RiceID             string     `json:"rice_id"`
	Quantity           int32      `json:"quantity"`
	Status             string     `json:"status"`
	Grade              *string    `json:"grade"`
	MillingDate        *time.Time `json:"milling_date"`
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

func RiceOrdersResponse(riceOrders []*model.RiceOrder) []*riceOrdersResponse {
	res := make([]*riceOrdersResponse, 0)
	for _, riceOrder := range riceOrders {
		res = append(res, &riceOrdersResponse{
			ID:                riceOrder.ID,
			Quantity:          riceOrder.Quantity,
			Status:            riceOrder.Status,
			LastTransactionAt: riceOrder.GetLastTransactionAt(),
		})
	}

	return res
}

func RiceOrderResponse(riceOrder *model.RiceOrder) *riceOrderResponse {
	return &riceOrderResponse{
		ID:                 riceOrder.ID,
		SellerID:           riceOrder.SellerID,
		RiceID:             riceOrder.RiceID,
		Quantity:           riceOrder.Quantity,
		Status:             riceOrder.Status,
		Grade:              riceOrder.RiceInstance.GetGrade(),
		MillingDate:        riceOrder.RiceInstance.GetMillingDate(),
		StorageTemperature: riceOrder.RiceInstance.GetStorageTemperature(),
		StorageHumidity:    riceOrder.RiceInstance.GetStorageHumidity(),
		OrderedAt:          riceOrder.GetOrderedAt(),
		AcceptedAt:         riceOrder.GetAcceptedAt(),
		RejectedAt:         riceOrder.GetRejectedAt(),
		RejectedReason:     riceOrder.RejectReason,
		ProcessingAt:       riceOrder.GetProcessingAt(),
		AvailableAt:        riceOrder.GetAvailableAt(),
		ShippedAt:          riceOrder.GetShippedAt(),
		ReceivedAt:         riceOrder.GetReceivedAt(),
	}
}
