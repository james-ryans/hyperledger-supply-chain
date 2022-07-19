package userservice

import (
	usermodel "github.com/meneketehe/hehe/app/model/user"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	UserRepository usermodel.UserRepository
}

type UserServiceConfig struct {
	UserRepository usermodel.UserRepository
}

func NewUserService(c *UserServiceConfig) usermodel.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

func (s *userService) GetUserByID(ID string) (*usermodel.User, error) {
	return s.UserRepository.FindByID(ID)
}

func (s *userService) CreateUser(user *usermodel.User) (*usermodel.User, error) {
	if err := s.UserRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
