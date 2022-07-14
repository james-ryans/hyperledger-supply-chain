package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateSeedOrderRequest struct {
	SupplierID       string  `json:"supplier_id"`
	SeedID           string  `json:"seed_id"`
	RiceGrainOrderID string  `json:"rice_grain_order_id"`
	Weight           float32 `json:"weight"`
}

func (r CreateSeedOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SupplierID, validation.Required, is.UUID),
		validation.Field(&r.SeedID, validation.Required, is.UUID),
		validation.Field(&r.RiceGrainOrderID, validation.Required, is.UUID),
		validation.Field(&r.Weight, validation.Required),
	)
}

type RejectRiceGrainOrderRequest struct {
	Reason string `json:"reason"`
}

func (r RejectRiceGrainOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Reason, validation.Required),
	)
}

func (r RejectRiceGrainOrderRequest) Sanitize() {
	r.Reason = strings.TrimSpace(r.Reason)
}
