package repository

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type organizationRepository struct {
	Fabric *client.Gateway
}

func NewOrganizationRepository(fabric *client.Gateway) model.OrganizationRepository {
	return &organizationRepository{
		Fabric: fabric,
	}
}

func (r *organizationRepository) FindByID(ID string) (*model.Organization, error) {
	return &model.Organization{
		ID:   ID,
		Type: "supplier",
		Name: "Supplier 0",
		Location: model.Location{
			Province:   "North Sumatra",
			City:       "Medan",
			District:   "Medan Kota",
			PostalCode: "20212",
			Address:    "Jl. Thamrin",
			Coordinate: model.Coordinate{
				Latitude:  123.1,
				Longitude: 321.1,
			},
		},
		ContactInfo: model.ContactInfo{
			Phone: "081234567890",
			Email: "supplier0@hehe.com",
		},
	}, nil
}
