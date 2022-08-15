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
	GetChannelNameByFile(filename string) (string, error)
	GetAllChannels() ([]*GlobalChannel, error)
	GetChannel(ID string) (*GlobalChannel, error)
	GetChannelByName(name string) (*GlobalChannel, error)
	CheckNameExists(name string) (bool, error)
	CreateChannel(ch *GlobalChannel) (*GlobalChannel, error)
	CreateChannelBlock(ch *GlobalChannel, orgs []*GlobalOrganization) error
	GetJoinedChannels() ([]string, error)
	JoinChannel(name, blockPath string, orgs []*GlobalOrganization) error
}

type GlobalChannelRepository interface {
	FindAll() ([]*GlobalChannel, error)
	FindByID(ID string) (*GlobalChannel, error)
	FindByName(name string) (*GlobalChannel, error)
	Create(ch *GlobalChannel) (*GlobalChannel, error)
}

func (ch *GlobalChannel) OrgsID() []string {
	orgsID := append(ch.SuppliersID, ch.ProducersID...)
	orgsID = append(orgsID, ch.ManufacturersID...)
	orgsID = append(orgsID, ch.DistributorsID...)
	orgsID = append(orgsID, ch.RetailersID...)

	return orgsID
}
