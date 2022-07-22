package usermodel

type UserAccount struct {
	ID       string `json:"_id"`
	Rev      string `json:"_rev,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAccountService interface {
	GetUserAccountByID(ID string) (*UserAccount, error)
	GetUserAccountByEmail(email string) (*UserAccount, error)
	Register(account *UserAccount) (*UserAccount, error)
	Login(email, password string) (*UserAccount, error)
	EditProfile(account *UserAccount) (*UserAccount, error)
}

type UserAccountRepository interface {
	FindByID(ID string) (*UserAccount, error)
	FindByEmail(email string) (*UserAccount, error)
	Create(account *UserAccount) (*UserAccount, error)
	Update(account *UserAccount) (*UserAccount, error)
}
