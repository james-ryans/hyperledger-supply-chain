package model

type RiceSack struct {
	ID              string   `json:"id"`
	RiceOrderID     []string `json:"rice_order_id"`
	RiceStockpileID string   `json:"rice_stockpile_id"`
	Code            string   `json:"code"`
}

type RiceSackService interface {
	GetAllRiceSack(channelID, stockpileID string) ([]*RiceSack, error)
	GetAllRiceSackByRiceOrderID(channelID, riceOrderID string) ([]*RiceSack, error)
	GetRiceSack(channelID, ID string) (*RiceSack, error)
}

type RiceSackRepository interface {
	FindAll(channelID, stockpileID string) ([]*RiceSack, error)
	FindAllByRiceOrderID(channelID, riceOrderID string) ([]*RiceSack, error)
	FindByID(channelID, ID string) (*RiceSack, error)
}
