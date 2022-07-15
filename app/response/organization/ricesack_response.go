package response

import (
	"encoding/base64"

	"github.com/meneketehe/hehe/app/model"
)

type riceSackResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type riceSackQRCodeResponse struct {
	QRCode string `json:"qr_code"`
}

func RiceSacksResponse(sacks []*model.RiceSack) []*riceSackResponse {
	res := make([]*riceSackResponse, 0)
	for _, sack := range sacks {
		res = append(res, &riceSackResponse{
			ID:   sack.ID,
			Code: sack.Code,
		})
	}

	return res
}

func RiceSackQrCodeResponse(qrCode []byte) *riceSackQRCodeResponse {
	return &riceSackQRCodeResponse{
		QRCode: base64.StdEncoding.EncodeToString(qrCode),
	}
}
