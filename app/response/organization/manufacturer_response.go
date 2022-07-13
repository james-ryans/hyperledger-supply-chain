package response

import "github.com/meneketehe/hehe/app/model"

type manufacturerResponse struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	District   string  `json:"district"`
	PostalCode string  `json:"postal_code"`
	Address    string  `json:"address"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
}

func ManufacturersResponse(manufacturers []*model.Manufacturer) []*manufacturerResponse {
	res := make([]*manufacturerResponse, 0)
	for _, manufacturer := range manufacturers {
		res = append(res, ManufacturerResponse(manufacturer))
	}

	return res
}

func ManufacturerResponse(manufacturer *model.Manufacturer) *manufacturerResponse {
	return &manufacturerResponse{
		ID:         manufacturer.ID,
		Type:       manufacturer.Type,
		Name:       manufacturer.Name,
		Province:   manufacturer.Location.Province,
		City:       manufacturer.Location.City,
		District:   manufacturer.Location.District,
		PostalCode: manufacturer.Location.PostalCode,
		Address:    manufacturer.Location.Address,
		Latitude:   manufacturer.Location.Coordinate.Latitude,
		Longitude:  manufacturer.Location.Coordinate.Longitude,
		Phone:      manufacturer.ContactInfo.Phone,
		Email:      manufacturer.ContactInfo.Email,
	}
}
