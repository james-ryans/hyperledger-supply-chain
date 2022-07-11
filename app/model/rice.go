package model

type Rice struct {
	ID             string  `json:"id"`
	ManufacturerID string  `json:"manufacturer_id"`
	BrandName      string  `json:"brand_name"`
	Weight         float32 `json:"weight"`
	Texture        string  `json:"texture"`
	AmyloseRate    float32 `json:"amylose_rate"`
}
