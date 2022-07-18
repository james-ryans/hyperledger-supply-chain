package usermodel

import "github.com/meneketehe/hehe/app/model"

type Organization struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

func FromOrgModel(org *model.Organization) Organization {
	return Organization{
		ID:   org.ID,
		Name: org.Name,
		Location: Location{
			Province:   org.Location.Province,
			City:       org.Location.City,
			District:   org.Location.District,
			PostalCode: org.Location.PostalCode,
			Address:    org.Location.Address,
			Coordinate: Coordinate{
				Latitude:  org.Location.Coordinate.Latitude,
				Longitude: org.Location.Coordinate.Longitude,
			},
		},
	}
}
