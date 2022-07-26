package response

import "github.com/meneketehe/hehe/app/model"

type globalOrganizationResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	Code       *string `json:"code"`
	MSPID      string  `json:"msp_id"`
	Domain     string  `json:"domain"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	District   string  `json:"district"`
	PostalCode string  `json:"postal_code"`
	Address    string  `json:"address"`
	Longitude  float32 `json:"longitude"`
	Latitude   float32 `json:"latitude"`
}

func GlobalOrganizationsResponse(orgs []*model.GlobalOrganization) []*globalOrganizationResponse {
	res := make([]*globalOrganizationResponse, 0)
	for _, org := range orgs {
		res = append(res, GlobalOrganizationResponse(org))
	}

	return res
}

func GlobalOrganizationResponse(org *model.GlobalOrganization) *globalOrganizationResponse {
	var code *string
	if org.Code != "" {
		code = &org.Code
	}

	return &globalOrganizationResponse{
		ID:         org.ID,
		Name:       org.Name,
		Role:       org.Role,
		Code:       code,
		MSPID:      org.MSPID,
		Domain:     org.Domain,
		Phone:      org.ContactInfo.Phone,
		Email:      org.ContactInfo.Email,
		Province:   org.Location.Province,
		City:       org.Location.City,
		District:   org.Location.District,
		PostalCode: org.Location.PostalCode,
		Address:    org.Location.Address,
		Longitude:  org.Location.Coordinate.Longitude,
		Latitude:   org.Location.Coordinate.Latitude,
	}
}
