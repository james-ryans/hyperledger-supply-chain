package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateDistributedRiceOrderRequest struct {
	DistributorID   string `json:"distributor_id"`
	RiceStockpileID string `json:"rice_stockpile_id"`
	Quantity        int32  `json:"quantity"`
}

func (r CreateDistributedRiceOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.DistributorID, validation.Required, is.UUID),
		validation.Field(&r.RiceStockpileID, validation.Required, is.UUID),
		validation.Field(&r.Quantity, validation.Required),
	)
}
