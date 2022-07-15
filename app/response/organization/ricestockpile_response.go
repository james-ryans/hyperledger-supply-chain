package response

import "github.com/meneketehe/hehe/app/model"

type riceStockpileResponse struct {
	ID     string `json:"id"`
	RiceID string `json:"rice_id"`
	Stock  int32  `json:"stock"`
}

func RiceStockpilesResponse(piles []*model.RiceStockpile) []*riceStockpileResponse {
	res := make([]*riceStockpileResponse, 0)
	for _, pile := range piles {
		res = append(res, RiceStockpileResponse(pile))
	}

	return res
}

func RiceStockpileResponse(pile *model.RiceStockpile) *riceStockpileResponse {
	return &riceStockpileResponse{
		ID:     pile.ID,
		RiceID: pile.RiceID,
		Stock:  pile.Stock,
	}
}
