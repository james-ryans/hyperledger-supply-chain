package model

type Manufacturer struct {
	Organization
}

type ManufacturerService interface {
	GetAllManufacturers(channelID string) ([]*Manufacturer, error)
	GetManufacturerByID(channelID, ID string) (*Manufacturer, error)
	CreateManufacturer(channelID string, manufacturer *Manufacturer) (*Manufacturer, error)
}

type ManufacturerRepository interface {
	FindAll(channelID string) ([]*Manufacturer, error)
	FindByID(channelID, ID string) (*Manufacturer, error)
	Create(channelID string, manufacturer *Manufacturer) error
}
