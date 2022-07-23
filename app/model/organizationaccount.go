package model

type OrganizationAccount struct {
	ID             string `json:"_id"`
	Rev            string `json:"_rev,omitempty"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
}

type OrganizationAccountService interface {
	GetOrganizationAccountByID(ID string) (*OrganizationAccount, error)
	GetOrganizationAccountByEmail(email string) (*OrganizationAccount, error)
	Login(email, password string) (*OrganizationAccount, error)
	ChangePassword(account *OrganizationAccount) (*OrganizationAccount, error)
}

type OrganizationAccountRepository interface {
	FindByID(ID string) (*OrganizationAccount, error)
	FindByEmail(email string) (*OrganizationAccount, error)
	Update(account *OrganizationAccount) (*OrganizationAccount, error)
}
