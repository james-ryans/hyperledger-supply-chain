package userrequest

import (
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Name, validation.Required, validation.Length(4, 31)),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 255)),
		validation.Field(&r.PasswordConfirmation, validation.Required, validation.Length(6, 255), validation.Match(regexp.MustCompile(fmt.Sprintf("^%s$", *&r.Password))).Error("must be equal to password")),
	)
}

func (r *RegisterRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.PasswordConfirmation = strings.TrimSpace(r.PasswordConfirmation)
}

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

type EditProfileRequest struct {
	Name                    string `json:"name"`
	OldPassword             string `json:"old_password"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

func (r EditProfileRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.OldPassword, validation.Required, validation.Length(6, 255)),
		validation.Field(&r.NewPassword, validation.Required, validation.Length(6, 255)),
		validation.Field(&r.NewPasswordConfirmation, validation.Required, validation.Length(6, 255), validation.Match(regexp.MustCompile(fmt.Sprintf("^%s$", *&r.NewPassword))).Error("must be equal to new password")),
	)
}

func (r *EditProfileRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
	r.OldPassword = strings.TrimSpace(r.OldPassword)
	r.NewPassword = strings.TrimSpace(r.NewPassword)
	r.NewPasswordConfirmation = strings.TrimSpace(r.NewPasswordConfirmation)
}
