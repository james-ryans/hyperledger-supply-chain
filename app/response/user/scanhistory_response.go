package userresponse

import (
	"time"

	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type scanHistoryResponse struct {
	ID           string    `json:"id"`
	RiceSackCode string    `json:"rice_sack_code"`
	Name         string    `json:"name"`
	ScanAt       time.Time `json:"scan_at"`
}

func ScanHistoriesResponse(scanHistories []*usermodel.ScanHistory) []*scanHistoryResponse {
	res := make([]*scanHistoryResponse, 0)
	for _, scanHistory := range scanHistories {
		res = append(res, ScanHistoryResponse(scanHistory))
	}

	return res
}

func ScanHistoryResponse(scanHistory *usermodel.ScanHistory) *scanHistoryResponse {
	return &scanHistoryResponse{
		ID:           scanHistory.ID,
		RiceSackCode: scanHistory.RiceSackCode,
		Name:         scanHistory.Name,
		ScanAt:       scanHistory.ScanAt,
	}
}
