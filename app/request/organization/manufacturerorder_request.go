package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateRiceGrainOrderRequest struct {
	ProducerID  string  `json:"producer_id"`
	RiceGrainID string  `json:"rice_grain_id"`
	RiceOrderID string  `json:"rice_order_id"`
	Weight      float32 `json:"weight"`
}

func (r CreateRiceGrainOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProducerID, validation.Required, is.UUID),
		validation.Field(&r.RiceGrainID, validation.Required, is.UUID),
		validation.Field(&r.RiceOrderID, validation.Required, is.UUID),
		validation.Field(&r.Weight, validation.Required),
	)
}

type RejectRiceOrderRequest struct {
	Reason string `json:"reason"`
}

func (r RejectRiceOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Reason, validation.Required),
	)
}

func (r RejectRiceOrderRequest) Sanitize() {
	r.Reason = strings.TrimSpace(r.Reason)
}
