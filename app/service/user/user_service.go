package userservice

import (
	usermodel "github.com/meneketehe/hehe/app/model/user"
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

func (s *userService) UpdateUser(user *usermodel.User) (*usermodel.User, error) {
	if err := s.UserRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
