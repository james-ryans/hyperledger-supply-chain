package response

import "github.com/meneketehe/hehe/app/model"

type retailerResponse struct {
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

func RetailersResponse(retailers []*model.Retailer) []*retailerResponse {
	res := make([]*retailerResponse, 0)
	for _, retailer := range retailers {
		res = append(res, RetailerResponse(retailer))
	}

	return res
}

func RetailerResponse(retailer *model.Retailer) *retailerResponse {
	return &retailerResponse{
		ID:         retailer.ID,
		Type:       retailer.Type,
		Name:       retailer.Name,
		Province:   retailer.Location.Province,
		City:       retailer.Location.City,
		District:   retailer.Location.District,
		PostalCode: retailer.Location.PostalCode,
		Address:    retailer.Location.Address,
		Latitude:   retailer.Location.Coordinate.Latitude,
		Longitude:  retailer.Location.Coordinate.Longitude,
		Phone:      retailer.ContactInfo.Phone,
		Email:      retailer.ContactInfo.Email,
	}
}
