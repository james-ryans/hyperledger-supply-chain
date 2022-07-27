package service

import (
	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/model"
)

type globalOrganizationService struct {
	GlobalOrganizationRepository model.GlobalOrganizationRepository
}

type GlobalOrganizationServiceConfig struct {
	GlobalOrganizationRepository model.GlobalOrganizationRepository
}

func NewGlobalOrganizationService(c *GlobalOrganizationServiceConfig) model.GlobalOrganizationService {
	return &globalOrganizationService{
		GlobalOrganizationRepository: c.GlobalOrganizationRepository,
	}
}

func (s *globalOrganizationService) GetAllOrganizations(filters map[string]string) ([]*model.GlobalOrganization, error) {
	orgs, err := s.GlobalOrganizationRepository.FindAll(filters)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (s *globalOrganizationService) GetOrganizationsByIDs(IDs []string) ([]*model.GlobalOrganization, error) {
	orgs, err := s.GlobalOrganizationRepository.FindByIDs(IDs)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (s *globalOrganizationService) GetOrganization(ID string) (*model.GlobalOrganization, error) {
	org, err := s.GlobalOrganizationRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (s *globalOrganizationService) CheckCodeExists(code string) (bool, error) {
	org, err := s.GlobalOrganizationRepository.FindByCode(code)
	if err != nil {
		return false, err
	}

	return org != nil, nil
}

func (s *globalOrganizationService) CheckMSPIDExists(MSPID string) (bool, error) {
	org, err := s.GlobalOrganizationRepository.FindByMSPID(MSPID)
	if err != nil {
		return false, err
	}

	return org != nil, nil
}

func (s *globalOrganizationService) CheckDomainExists(domain string) (bool, error) {
	org, err := s.GlobalOrganizationRepository.FindByDomain(domain)
	if err != nil {
		return false, err
	}

	return org != nil, nil
}

func (s *globalOrganizationService) CreateOrganization(org *model.GlobalOrganization) (*model.GlobalOrganization, error) {
	org.ID = uuid.New().String()

	org, err := s.GlobalOrganizationRepository.Create(org)
	if err != nil {
		return nil, err
	}

	return org, nil
}
