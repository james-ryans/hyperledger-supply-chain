package response

import (
	"time"

	"github.com/meneketehe/hehe/app/model"
)

type riceGrainOrdersResponse struct {
	ID                string     `json:"id"`
	Weight            float32    `json:"weight"`
	Status            string     `json:"status"`
	LastTransactionAt *time.Time `json:"last_transaction_at"`
}

type riceGrainOrderDetailedResponse struct {
	Order    *riceGrainOrderResponse `json:"order"`
	Previous *seedOrderResponse      `json:"previous"`
	Next     *riceOrderResponse      `json:"next"`
}

type riceGrainOrderResponse struct {
	ID                 string     `json:"id"`
	OrdererID          string     `json:"orderer_id"`
	SellerID           string     `json:"seller_id"`
	RiceGrainID        string     `json:"rice_grain_id"`
	RiceOrderID        string     `json:"rice_order_id"`
	Weight             float32    `json:"weight"`
	Status             string     `json:"status"`
	PlowMethod         *string    `json:"plow_method"`
	SowMethod          *string    `json:"sow_method"`
	Irrigation         *string    `json:"irrigation"`
	Fertilization      *string    `json:"fertilization"`
	PlantDate          *time.Time `json:"plant_date"`
	HarvestDate        *time.Time `json:"harvest_date"`
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

func RiceGrainOrdersResponse(riceGrainOrders []*model.RiceGrainOrder) []*riceGrainOrdersResponse {
	res := make([]*riceGrainOrdersResponse, 0)
	for _, riceGrainOrder := range riceGrainOrders {
		res = append(res, &riceGrainOrdersResponse{
			ID:                riceGrainOrder.ID,
			Weight:            riceGrainOrder.Weight,
			Status:            riceGrainOrder.Status,
			LastTransactionAt: riceGrainOrder.GetLastTransactionAt(),
		})
	}

	return res
}

func RiceGrainOrderDetailedResponse(riceGrainOrder *model.RiceGrainOrder, riceOrder *model.RiceOrder, seedOrder *model.SeedOrder) *riceGrainOrderDetailedResponse {
	res := &riceGrainOrderDetailedResponse{
		Order: RiceGrainOrderResponse(riceGrainOrder),
		Next:  RiceOrderResponse(riceOrder),
	}

	if seedOrder != nil {
		res.Previous = SeedOrderResponse(seedOrder)
	}

	return res
}

func RiceGrainOrderResponse(riceGrainOrder *model.RiceGrainOrder) *riceGrainOrderResponse {
	return &riceGrainOrderResponse{
		ID:                 riceGrainOrder.ID,
		OrdererID:          riceGrainOrder.OrdererID,
		SellerID:           riceGrainOrder.SellerID,
		RiceGrainID:        riceGrainOrder.RiceGrainID,
		RiceOrderID:        riceGrainOrder.RiceOrderID,
		Weight:             riceGrainOrder.Weight,
		Status:             riceGrainOrder.Status,
		PlowMethod:         riceGrainOrder.RiceGrainInstance.GetPlowMethod(),
		SowMethod:          riceGrainOrder.RiceGrainInstance.GetSowMethod(),
		Irrigation:         riceGrainOrder.RiceGrainInstance.GetIrrigation(),
		Fertilization:      riceGrainOrder.RiceGrainInstance.GetFertilization(),
		PlantDate:          riceGrainOrder.RiceGrainInstance.GetPlantDate(),
		HarvestDate:        riceGrainOrder.RiceGrainInstance.GetHarvestDate(),
		StorageTemperature: riceGrainOrder.RiceGrainInstance.GetStorageTemperature(),
		StorageHumidity:    riceGrainOrder.RiceGrainInstance.GetStorageHumidity(),
		OrderedAt:          riceGrainOrder.GetOrderedAt(),
		AcceptedAt:         riceGrainOrder.GetAcceptedAt(),
		RejectedAt:         riceGrainOrder.GetRejectedAt(),
		RejectedReason:     riceGrainOrder.RejectReason,
		ProcessingAt:       riceGrainOrder.GetProcessingAt(),
		AvailableAt:        riceGrainOrder.GetAvailableAt(),
		ShippedAt:          riceGrainOrder.GetShippedAt(),
		ReceivedAt:         riceGrainOrder.GetReceivedAt(),
	}
}
