package usermodel

type User struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Email         string        `json:"email"`
	Password      string        `json:"password"`
	ScanHistories []ScanHistory `json:"scan_histories"`
}

type UserService interface {
	GetUserByID(ID string) (*User, error)
	CreateUser(user *User) (*User, error)
	HashPassword(password string) (string, error)
}

type UserRepository interface {
	FindByID(ID string) (*User, error)
	Create(user *User) error
}
