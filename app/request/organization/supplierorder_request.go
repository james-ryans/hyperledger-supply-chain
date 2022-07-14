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
