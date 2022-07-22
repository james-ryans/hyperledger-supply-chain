package userservice

import (
	"fmt"

	"github.com/google/uuid"
	usermodel "github.com/meneketehe/hehe/app/model/user"
	"golang.org/x/crypto/bcrypt"
)

type userAccountService struct {
	UserAccountRepository usermodel.UserAccountRepository
}

type UserAccountServiceConfig struct {
	UserAccountRepository usermodel.UserAccountRepository
}

func NewUserAccountService(c *UserAccountServiceConfig) usermodel.UserAccountService {
	return &userAccountService{
		UserAccountRepository: c.UserAccountRepository,
	}
}

func (s userAccountService) GetUserAccountByID(ID string) (*usermodel.UserAccount, error) {
	account, err := s.UserAccountRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s userAccountService) GetUserAccountByEmail(email string) (*usermodel.UserAccount, error) {
	account, err := s.UserAccountRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s userAccountService) Register(account *usermodel.UserAccount) (*usermodel.UserAccount, error) {
	hashedPassword, err := HashPassword(account.Password)
	if err != nil {
		return nil, err
	}

	account.ID = uuid.New().String()
	account.Password = hashedPassword

	return s.UserAccountRepository.Create(account)
}

func (s userAccountService) Login(email, password string) (*usermodel.UserAccount, error) {
	account, err := s.UserAccountRepository.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email and password combination")
	}

	if match := CheckPasswordHash(password, account.Password); !match {
		return nil, fmt.Errorf("invalid email and password combination")
	}

	return account, nil
}

func (s userAccountService) EditProfile(account *usermodel.UserAccount) (*usermodel.UserAccount, error) {
	hashedPassword, err := HashPassword(account.Password)
	if err != nil {
		return nil, err
	}
	account.Password = hashedPassword

	return s.UserAccountRepository.Update(account)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
