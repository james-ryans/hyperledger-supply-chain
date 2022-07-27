package model

type GlobalChannel struct {
	ID              string   `json:"_id"`
	Rev             string   `json:"_rev,omitempty"`
	Name            string   `json:"name"`
	SuppliersID     []string `json:"suppliers_id"`
	ProducersID     []string `json:"producers_id"`
	ManufacturersID []string `json:"manufacturers_id"`
	DistributorsID  []string `json:"distributors_id"`
	RetailersID     []string `json:"retailers_id"`
}

type GlobalChannelService interface {
	GetAllChannels() ([]*GlobalChannel, error)
	GetChannel(ID string) (*GlobalChannel, error)
	CheckNameExists(name string) (bool, error)
	CreateChannel(ch *GlobalChannel) (*GlobalChannel, error)
}

type GlobalChannelRepository interface {
	FindAll() ([]*GlobalChannel, error)
	FindByID(ID string) (*GlobalChannel, error)
	FindByName(name string) (*GlobalChannel, error)
	Create(ch *GlobalChannel) (*GlobalChannel, error)
}
