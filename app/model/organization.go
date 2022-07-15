package model

type Organization struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
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

type OrganizationService interface {
	GetMe() (*Organization, error)
}
