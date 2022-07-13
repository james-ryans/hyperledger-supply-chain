package model

type Distributor struct {
	Vendor
}

type DistributorService interface {
	GetAllDistributors(channelID string) ([]*Distributor, error)
	GetDistributorByID(channelID, ID string) (*Distributor, error)
	CreateDistributor(channelID string, distributor *Distributor) (*Distributor, error)
}

type DistributorRepository interface {
	FindAll(channelID string) ([]*Distributor, error)
	FindByID(channelID, ID string) (*Distributor, error)
	Create(channelID string, distributor *Distributor) error
}
