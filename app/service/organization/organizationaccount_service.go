package service

import (
	"fmt"

	"github.com/meneketehe/hehe/app/model"
	"golang.org/x/crypto/bcrypt"
)

type organizationAccountService struct {
	OrganizationAccountRepository model.OrganizationAccountRepository
}

type OrganizationAccountServiceConfig struct {
	OrganizationAccountRepository model.OrganizationAccountRepository
}

func NewOrganizationAccountService(c *OrganizationAccountServiceConfig) model.OrganizationAccountService {
	return &organizationAccountService{
		OrganizationAccountRepository: c.OrganizationAccountRepository,
	}
}

func (s *organizationAccountService) GetOrganizationAccountByID(ID string) (*model.OrganizationAccount, error) {
	account, err := s.OrganizationAccountRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *organizationAccountService) GetOrganizationAccountByEmail(email string) (*model.OrganizationAccount, error) {
	account, err := s.OrganizationAccountRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *organizationAccountService) Login(email, password string) (*model.OrganizationAccount, error) {
	account, err := s.OrganizationAccountRepository.FindByEmail(email)
	if err != nil || account == nil {
		return nil, fmt.Errorf("invalid email and password combination")
	}

	if match := CheckPasswordHash(password, account.Password); !match {
		return nil, fmt.Errorf("invalid email and password combination")
	}

	return account, nil
}

func (s *organizationAccountService) ChangePassword(account *model.OrganizationAccount) (*model.OrganizationAccount, error) {
	hashedPassword, err := HashPassword(account.Password)
	if err != nil {
		return nil, err
	}
	account.Password = hashedPassword

	return s.OrganizationAccountRepository.Update(account)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
