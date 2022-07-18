package usermodel

type Location struct {
	Province   string     `json:"province"`
	City       string     `json:"city"`
	District   string     `json:"district"`
	PostalCode string     `json:"postal_code"`
	Address    string     `json:"address"`
	Coordinate Coordinate `json:"coordinate"`
}
