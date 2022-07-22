package usermodel

type User struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Email         string        `json:"email"`
	ScanHistories []ScanHistory `json:"scan_histories"`
}

func UserFromAccount(account *UserAccount) *User {
	return &User{
		ID:            account.ID,
		Name:          account.Name,
		Email:         account.Email,
		ScanHistories: []ScanHistory{},
	}
}

type UserService interface {
	GetUserByID(ID string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(user *User) (*User, error)
}

type UserRepository interface {
	FindByID(ID string) (*User, error)
	Create(user *User) error
	Update(user *User) error
}
