package model

type Organization struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Location    Location    `json:"location"`
	ContactInfo ContactInfo `json:"contact_info"`
}

type Location struct {
	Province   string     `json:"province"`
	City       string     `json:"city"`
	District   string     `json:"district"`
	PostalCode string     `json:"postal_code"`
	Address    string     `json:"address"`
	Coordinate Coordinate `json:"coordinate"`
}

type ContactInfo struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type Coordinate struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type OrganizationResponse struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	District   string  `json:"district"`
	PostalCode string  `json:"postal_code"`
	Address    string  `json:"address"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
}

func ToOrganizationResponse(organization *Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:         organization.ID,
		Type:       organization.Type,
		Name:       organization.Name,
		Phone:      organization.ContactInfo.Phone,
		Email:      organization.ContactInfo.Email,
		Province:   organization.Location.Province,
		City:       organization.Location.City,
		District:   organization.Location.District,
		PostalCode: organization.Location.PostalCode,
		Address:    organization.Location.Address,
		Latitude:   organization.Location.Coordinate.Latitude,
		Longitude:  organization.Location.Coordinate.Longitude,
	}
}

type OrganizationService interface {
	GetMe(ID string) (*Organization, error)
}

type OrganizationRepository interface {
	FindByID(ID string) (*Organization, error)
}
