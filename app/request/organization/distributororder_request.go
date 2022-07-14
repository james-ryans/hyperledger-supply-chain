package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateRiceOrderRequest struct {
	ManufacturerID string `json:"manufacturer_id"`
	RiceID         string `json:"rice_id"`
	Quantity       int32  `json:"quantity"`
}

func (r CreateRiceOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ManufacturerID, validation.Required, is.UUID),
		validation.Field(&r.RiceID, validation.Required, is.UUID),
		validation.Field(&r.Quantity, validation.Required),
	)
}
