package response

import "github.com/meneketehe/hehe/app/model"

type producerResponse struct {
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

func ProducersResponse(producers []*model.Producer) []*producerResponse {
	res := make([]*producerResponse, 0)
	for _, producer := range producers {
		res = append(res, ProducerResponse(producer))
	}

	return res
}

func ProducerResponse(producer *model.Producer) *producerResponse {
	return &producerResponse{
		ID:         producer.ID,
		Type:       producer.Type,
		Name:       producer.Name,
		Province:   producer.Location.Province,
		City:       producer.Location.City,
		District:   producer.Location.District,
		PostalCode: producer.Location.PostalCode,
		Address:    producer.Location.Address,
		Latitude:   producer.Location.Coordinate.Latitude,
		Longitude:  producer.Location.Coordinate.Longitude,
		Phone:      producer.ContactInfo.Phone,
		Email:      producer.ContactInfo.Email,
	}
}
