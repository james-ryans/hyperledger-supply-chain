package request

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 255)),
	)
}

func (r *LoginRequest) Sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

type ChangePasswordRequest struct {
	OldPassword             string `json:"old_password"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

func (r ChangePasswordRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OldPassword, validation.Required, validation.Length(6, 255)),
		validation.Field(&r.NewPassword, validation.Required, validation.Length(6, 255)),
		validation.Field(&r.NewPasswordConfirmation, validation.Required, validation.Length(6, 255)),
	)
}

func (r *ChangePasswordRequest) Sanitize() {
	r.OldPassword = strings.TrimSpace(r.OldPassword)
	r.NewPassword = strings.TrimSpace(r.NewPassword)
	r.NewPasswordConfirmation = strings.TrimSpace(r.NewPasswordConfirmation)
}
