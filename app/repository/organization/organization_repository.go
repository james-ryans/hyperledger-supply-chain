package repository

import (
	"github.com/meneketehe/hehe/app/model"
)

func GetMe(id string) (*model.Organization, error) {
	return &model.Organization{
		ID:   id,
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
