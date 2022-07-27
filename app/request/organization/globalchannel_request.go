package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateChannelRequest struct {
	Name            string   `json:"name"`
	SuppliersID     []string `json:"suppliers_id"`
	ProducersID     []string `json:"producers_id"`
	ManufacturersID []string `json:"manufacturers_id"`
	DistributorsID  []string `json:"distributors_id"`
	RetailersID     []string `json:"retailers_id"`
}

func (r CreateChannelRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, is.Alphanumeric, is.LowerCase),
		validation.Field(&r.SuppliersID, validation.Required, validation.Length(1, 0)),
		validation.Field(&r.ProducersID, validation.Required, validation.Length(1, 0)),
		validation.Field(&r.ManufacturersID, validation.Required, validation.Length(1, 0)),
		validation.Field(&r.DistributorsID, validation.Required, validation.Length(1, 0)),
		validation.Field(&r.RetailersID, validation.Required, validation.Length(1, 0)),
	)
}

func (r *CreateChannelRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
}
