package contract

type organization struct {
	ID          string      `json:"_id"`
	Name        string      `json:"name"`
	Location    location    `json:"location"`
	ContactInfo contactInfo `json:"contact_info"`
}

type location struct {
	Province   string     `json:"province"`
	City       string     `json:"city"`
	District   string     `json:"district"`
	PostalCode string     `json:"postal_code"`
	Address    string     `json:"address"`
	Coordinate coordinate `json:"coordinate"`
}

type contactInfo struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type coordinate struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
