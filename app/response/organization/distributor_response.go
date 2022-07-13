package response

import "github.com/meneketehe/hehe/app/model"

type distributorResponse struct {
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

func DistributorsResponse(distributors []*model.Distributor) []*distributorResponse {
	res := make([]*distributorResponse, 0)
	for _, distributor := range distributors {
		res = append(res, DistributorResponse(distributor))
	}

	return res
}

func DistributorResponse(distributor *model.Distributor) *distributorResponse {
	return &distributorResponse{
		ID:         distributor.ID,
		Type:       distributor.Type,
		Name:       distributor.Name,
		Province:   distributor.Location.Province,
		City:       distributor.Location.City,
		District:   distributor.Location.District,
		PostalCode: distributor.Location.PostalCode,
		Address:    distributor.Location.Address,
		Latitude:   distributor.Location.Coordinate.Latitude,
		Longitude:  distributor.Location.Coordinate.Longitude,
		Phone:      distributor.ContactInfo.Phone,
		Email:      distributor.ContactInfo.Email,
	}
}
