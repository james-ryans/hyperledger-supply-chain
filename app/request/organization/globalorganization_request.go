package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateOrganizationRequest struct {
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	Code       string  `json:"code"`
	MSPID      string  `json:"msp_id"`
	Domain     string  `json:"domain"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	District   string  `json:"district"`
	PostalCode string  `json:"postal_code"`
	Address    string  `json:"address"`
	Longitude  float32 `json:"longitude"`
	Latitude   float32 `json:"latitude"`
}

func (r CreateOrganizationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Role, validation.Required, validation.In("supplier", "producer", "manufacturer", "distributor", "retailer")),
		validation.Field(&r.Code, validation.Required.When(r.Role == "manufacturer").Error("required when role is manufacturer")),
		validation.Field(&r.MSPID, validation.Required, is.Alphanumeric),
		validation.Field(&r.Domain, validation.Required, is.Domain),
		validation.Field(&r.Phone, validation.Required, is.Digit),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Province, validation.Required),
		validation.Field(&r.City, validation.Required),
		validation.Field(&r.District, validation.Required),
		validation.Field(&r.PostalCode, validation.Required, is.Digit),
		validation.Field(&r.Address, validation.Required),
		validation.Field(&r.Longitude, validation.Required),
		validation.Field(&r.Latitude, validation.Required),
	)
}

func (r *CreateOrganizationRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Province = strings.TrimSpace(r.Province)
	r.City = strings.TrimSpace(r.City)
	r.District = strings.TrimSpace(r.District)
	r.Address = strings.TrimSpace(r.Address)
}
