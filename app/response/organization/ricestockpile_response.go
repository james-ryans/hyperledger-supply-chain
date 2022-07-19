package response

import "github.com/meneketehe/hehe/app/model"

type riceStockpilesResponse struct {
	ID     string `json:"id"`
	RiceID string `json:"rice_id"`
	Stock  int32  `json:"stock"`
}

type riceStockpileResponse struct {
	ID        string              `json:"id"`
	RiceID    string              `json:"rice_id"`
	Stock     int32               `json:"stock"`
	RiceSacks []*riceSackResponse `json:"rice_sacks"`
}

func RiceStockpilesResponse(piles []*model.RiceStockpile) []*riceStockpilesResponse {
	res := make([]*riceStockpilesResponse, 0)
	for _, pile := range piles {
		res = append(res, &riceStockpilesResponse{
			ID:     pile.ID,
			RiceID: pile.RiceID,
			Stock:  pile.Stock,
		})
	}

	return res
}

func RiceStockpileResponse(pile *model.RiceStockpile, sacks []*model.RiceSack) *riceStockpileResponse {
	return &riceStockpileResponse{
		ID:        pile.ID,
		RiceID:    pile.RiceID,
		Stock:     pile.Stock,
		RiceSacks: RiceSacksResponse(sacks),
	}
}
