package model

type RiceStockpile struct {
	ID       string `json:"id"`
	RiceID   string `json:"rice_id"`
	VendorID string `json:"vendor_id"`
	Stock    int32  `json:"stock"`
}

func (r *RiceStockpile) AddStock(quantity int32) {
	r.Stock += quantity
}

func (r *RiceStockpile) SubtractStock(quantity int32) {
	r.Stock -= quantity
}

type RiceStockpileService interface {
	GetAllRiceStockpile(channelID, vendorID string) ([]*RiceStockpile, error)
	GetRiceStockpileByID(channelID, ID string) (*RiceStockpile, error)
	GetRiceStockpileByVendorIDAndRiceID(channelID, vendorID, riceID string) (*RiceStockpile, error)
}

type RiceStockpileRepository interface {
	FindAll(channelID, vendorID string) ([]*RiceStockpile, error)
	FindByID(channelID, ID string) (*RiceStockpile, error)
	FindByVendorIDAndRiceID(channelID, vendorID, riceID string) (*RiceStockpile, error)
}
