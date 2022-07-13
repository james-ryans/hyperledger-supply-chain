package model

type Supplier struct {
	Organization
}

type SupplierService interface {
	GetAllSuppliers(channelID string) ([]*Supplier, error)
	GetSupplierByID(channelID, ID string) (*Supplier, error)
	CreateSupplier(channelID string, supplier *Supplier) (*Supplier, error)
}

type SupplierRepository interface {
	FindAll(channelID string) ([]*Supplier, error)
	FindByID(channelID, ID string) (*Supplier, error)
	Create(channelID string, supplier *Supplier) error
}
