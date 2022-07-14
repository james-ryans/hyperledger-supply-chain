package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RejectSeedOrderRequest struct {
	Reason string `json:"reason"`
}

func (r RejectSeedOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Reason, validation.Required),
	)
}

func (r RejectSeedOrderRequest) Sanitize() {
	r.Reason = strings.TrimSpace(r.Reason)
}

type ShipSeedOrderRequest struct {
	StorageTemperature float32 `json:"storage_temperature"`
	StorageHumidity    float32 `json:"storage_humidity"`
}

func (r ShipSeedOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.StorageTemperature, validation.Required, validation.Min(float64(0)), validation.Min(float64(100))),
		validation.Field(&r.StorageHumidity, validation.Required, validation.Min(float64(0)), validation.Min(float64(100))),
	)
}
