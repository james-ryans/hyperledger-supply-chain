package service

import (
	"github.com/meneketehe/hehe/app/model"
)

type organizationService struct {
	OrganizationRepository model.OrganizationRepository
}

type OrganizationServiceConfig struct {
	OrganizationRepository model.OrganizationRepository
}

func NewOrganizationService(c *OrganizationServiceConfig) model.OrganizationService {
	return &organizationService{
		OrganizationRepository: c.OrganizationRepository,
	}
}

func (s *organizationService) GetMe(ID string) (*model.Organization, error) {
	return s.OrganizationRepository.FindByID(ID)
}
