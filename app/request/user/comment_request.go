package userrequest

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type WriteCommentRequest struct {
	Text string `json:"text"`
}

func (r WriteCommentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Text, validation.Required),
	)
}

func (r *WriteCommentRequest) Sanitize() {
	r.Text = strings.TrimSpace(r.Text)
}
