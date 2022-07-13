package response

import "github.com/meneketehe/hehe/app/model"

type supplierResponse struct {
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

func SuppliersResponse(suppliers []*model.Supplier) []*supplierResponse {
	res := make([]*supplierResponse, 0)
	for _, supplier := range suppliers {
		res = append(res, SupplierResponse(supplier))
	}

	return res
}

func SupplierResponse(supplier *model.Supplier) *supplierResponse {
	return &supplierResponse{
		ID:         supplier.ID,
		Type:       supplier.Type,
		Name:       supplier.Name,
		Province:   supplier.Location.Province,
		City:       supplier.Location.City,
		District:   supplier.Location.District,
		PostalCode: supplier.Location.PostalCode,
		Address:    supplier.Location.Address,
		Latitude:   supplier.Location.Coordinate.Latitude,
		Longitude:  supplier.Location.Coordinate.Longitude,
		Phone:      supplier.ContactInfo.Phone,
		Email:      supplier.ContactInfo.Email,
	}
}
