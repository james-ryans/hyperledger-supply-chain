package service

import (
	"os"

	"github.com/meneketehe/hehe/app/model"
)

type organizationService struct {
}

func NewOrganizationService() model.OrganizationService {
	return &organizationService{}
}

func (s *organizationService) GetMe() (*model.Organization, error) {
	id := os.Getenv("ORG_ID")
	name := os.Getenv("ORG_NAME")
	orgType := os.Getenv("ORG_TYPE")

	return &model.Organization{
		ID:   id,
		Type: orgType,
		Name: name,
	}, nil
}
