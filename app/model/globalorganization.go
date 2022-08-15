package model

type GlobalOrganization struct {
	ID          string      `json:"_id"`
	Rev         string      `json:"_rev,omitempty"`
	Seq         int         `json:"seq"`
	Role        string      `json:"role"`
	Name        string      `json:"name"`
	Code        string      `json:"code,omitempty"`
	MSPID       string      `json:"msp_id"`
	Domain      string      `json:"domain"`
	Location    Location    `json:"location"`
	ContactInfo ContactInfo `json:"contact_info"`
}

type GlobalOrganizationService interface {
	GetAllOrganizations(filters map[string]string) ([]*GlobalOrganization, error)
	GetOrganizationsByIDs(ID []string) ([]*GlobalOrganization, error)
	GetOrganization(ID string) (*GlobalOrganization, error)
	CheckCodeExists(code string) (bool, error)
	CheckMSPIDExists(MSPID string) (bool, error)
	CheckDomainExists(domain string) (bool, error)
	CreateOrganization(org *GlobalOrganization) (*GlobalOrganization, error)
}

type GlobalOrganizationRepository interface {
	FindAll(filters map[string]string) ([]*GlobalOrganization, error)
	FindByIDs(IDs []string) ([]*GlobalOrganization, error)
	FindByID(ID string) (*GlobalOrganization, error)
	FindByCode(code string) (*GlobalOrganization, error)
	FindByMSPID(MSPID string) (*GlobalOrganization, error)
	FindByDomain(domain string) (*GlobalOrganization, error)
	Create(org *GlobalOrganization) (*GlobalOrganization, error)
}
