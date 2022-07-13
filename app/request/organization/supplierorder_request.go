package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RejectProducerOrderRequest struct {
	Reason string `json:"reason"`
}

func (r RejectProducerOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Reason, validation.Required),
	)
}

func (r RejectProducerOrderRequest) Sanitize() {
	r.Reason = strings.TrimSpace(r.Reason)
}
