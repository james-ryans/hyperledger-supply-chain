package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SeedRequest struct {
	VarietyName string  `json:"variety_name"`
	PlantAge    float32 `json:"plant_age"`
	PlantShape  string  `json:"plant_shape"`
	PlantHeight float32 `json:"plant_height"`
	LeafShape   string  `json:"leaf_shape"`
}

func (r SeedRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.VarietyName, validation.Required),
		validation.Field(&r.PlantAge, validation.Required),
		validation.Field(&r.PlantShape, validation.Required),
		validation.Field(&r.PlantHeight, validation.Required),
		validation.Field(&r.LeafShape, validation.Required),
	)
}

func (r *SeedRequest) Sanitize() {
	r.VarietyName = strings.TrimSpace(r.VarietyName)
	r.PlantShape = strings.TrimSpace(r.PlantShape)
	r.LeafShape = strings.TrimSpace(r.LeafShape)
}
