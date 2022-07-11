package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RiceGrainRequest struct {
	VarietyName string `json:"variety_name"`
	GrainShape  string `json:"grain_shape"`
	GrainColor  string `json:"grain_color"`
}

func (r RiceGrainRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.VarietyName, validation.Required),
		validation.Field(&r.GrainShape, validation.Required),
		validation.Field(&r.GrainColor, validation.Required),
	)
}

func (r *RiceGrainRequest) Sanitize() {
	r.VarietyName = strings.TrimSpace(r.VarietyName)
	r.GrainShape = strings.TrimSpace(r.GrainShape)
	r.GrainColor = strings.TrimSpace(r.GrainColor)
}
