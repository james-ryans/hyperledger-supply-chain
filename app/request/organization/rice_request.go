package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RiceRequest struct {
	BrandName   string  `json:"brand_name"`
	Weight      float32 `json:"weight"`
	Texture     string  `json:"texture"`
	AmyloseRate float32 `json:"amylose_rate"`
}

func (r RiceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BrandName, validation.Required),
		validation.Field(&r.Weight, validation.Required),
		validation.Field(&r.Texture, validation.Required),
		validation.Field(&r.AmyloseRate, validation.Required),
	)
}

func (r *RiceRequest) Sanitize() {
	r.BrandName = strings.TrimSpace(r.BrandName)
	r.Texture = strings.TrimSpace(r.Texture)
}
