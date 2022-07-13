package model

type Retailer struct {
	Vendor
}

type RetailerService interface {
	GetAllRetailers(channelID string) ([]*Retailer, error)
	GetRetailerByID(channelID, ID string) (*Retailer, error)
	CreateRetailer(channelID string, retailer *Retailer) (*Retailer, error)
}

type RetailerRepository interface {
	FindAll(channelID string) ([]*Retailer, error)
	FindByID(channelID, ID string) (*Retailer, error)
	Create(channelID string, retailer *Retailer) error
}
