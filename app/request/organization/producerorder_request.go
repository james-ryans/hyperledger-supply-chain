package request

import (
	"strings"
	"time"

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

type ShipRiceGrainOrderRequest struct {
	PlowMethod         string    `json:"plow_method"`
	SowMethod          string    `json:"sow_method"`
	Irrigation         string    `json:"irrigation"`
	Fertilization      string    `json:"fertilization"`
	PlantDate          time.Time `json:"plant_date"`
	HarvestDate        time.Time `json:"harvest_date"`
	StorageTemperature float32   `json:"storage_temperature"`
	StorageHumidity    float32   `json:"storage_humidity"`
}

func (r ShipRiceGrainOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.PlowMethod, validation.Required),
		validation.Field(&r.SowMethod, validation.Required),
		validation.Field(&r.Irrigation, validation.Required),
		validation.Field(&r.Fertilization, validation.Required),
		validation.Field(&r.PlantDate, validation.Required),
		validation.Field(&r.HarvestDate, validation.Required),
		validation.Field(&r.StorageTemperature, validation.Required, validation.Min(float64(0)), validation.Max(float64(100))),
		validation.Field(&r.StorageHumidity, validation.Required, validation.Min(float64(0)), validation.Max(float64(100))),
	)
}

func (r ShipRiceGrainOrderRequest) Sanitize() {
	r.PlowMethod = strings.TrimSpace(r.PlowMethod)
	r.SowMethod = strings.TrimSpace(r.SowMethod)
	r.Irrigation = strings.TrimSpace(r.Irrigation)
	r.Fertilization = strings.TrimSpace(r.Fertilization)
}
