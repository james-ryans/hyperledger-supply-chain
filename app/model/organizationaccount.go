package model

import "time"

type OrganizationAccount struct {
	ID             string    `json:"_id"`
	Rev            string    `json:"_rev,omitempty"`
	Role           string    `json:"role"`
	Type           string    `json:"type"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Password       string    `json:"password"`
	RegisteredAt   time.Time `json:"registered_at"`
}

type OrganizationAccountService interface {
	GetAllOrganizationUserAccounts() ([]*OrganizationAccount, error)
	GetOrganizationAccountByID(ID string) (*OrganizationAccount, error)
	GetOrganizationAccountByEmail(email string) (*OrganizationAccount, error)
	Login(email, password string) (*OrganizationAccount, error)
	ChangePassword(account *OrganizationAccount) (*OrganizationAccount, error)
	CreateUser(account *OrganizationAccount) (*OrganizationAccount, error)
	UpdateUser(account *OrganizationAccount, changePassword bool) (*OrganizationAccount, error)
	DeleteUser(account *OrganizationAccount) error
}

type OrganizationAccountRepository interface {
	FindAllUser() ([]*OrganizationAccount, error)
	FindByID(ID string) (*OrganizationAccount, error)
	FindByEmail(email string) (*OrganizationAccount, error)
	Create(account *OrganizationAccount) (*OrganizationAccount, error)
	Update(account *OrganizationAccount) (*OrganizationAccount, error)
	Delete(account *OrganizationAccount) error
}
