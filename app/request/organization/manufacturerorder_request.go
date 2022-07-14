package request

import (
	"strings"
	"time"

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

type ShipRiceOrderRequest struct {
	Grade              string    `json:"grade"`
	MillingDate        time.Time `json:"milling_date"`
	StorageTemperature float32   `json:"storage_temperature"`
	StorageHumidity    float32   `json:"storage_humidity"`
}

func (r ShipRiceOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Grade, validation.Required),
		validation.Field(&r.MillingDate, validation.Required),
		validation.Field(&r.StorageTemperature, validation.Required, validation.Min(float64(0)), validation.Max(float64(100))),
		validation.Field(&r.StorageHumidity, validation.Required, validation.Min(float64(0)), validation.Max(float64(100))),
	)
}

func (r ShipRiceOrderRequest) Sanitize() {
	r.Grade = strings.TrimSpace(r.Grade)
}
