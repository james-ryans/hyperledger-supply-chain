package service

import "github.com/meneketehe/hehe/app/model"

type supplierService struct {
	SupplierRepository model.SupplierRepository
}

type SupplierServiceConfig struct {
	SupplierRepository model.SupplierRepository
}

func NewSupplierService(c *SupplierServiceConfig) model.SupplierService {
	return &supplierService{
		SupplierRepository: c.SupplierRepository,
	}
}

func (s *supplierService) GetAllSuppliers(channelID string) ([]*model.Supplier, error) {
	return s.SupplierRepository.FindAll(channelID)
}

func (s *supplierService) GetSupplierByID(channelID, ID string) (*model.Supplier, error) {
	return s.SupplierRepository.FindByID(channelID, ID)
}

func (s *supplierService) CreateSupplier(channelID string, supplier *model.Supplier) (*model.Supplier, error) {
	if err := s.SupplierRepository.Create(channelID, supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}
