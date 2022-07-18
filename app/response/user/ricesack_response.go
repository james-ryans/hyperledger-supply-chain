package userresponse

import usermodel "github.com/meneketehe/hehe/app/model/user"

type riceSackResponse struct {
	usermodel.RiceSack
}

func RiceSackResponse(riceSack *usermodel.RiceSack) *riceSackResponse {
	return &riceSackResponse{
		RiceSack: *riceSack,
	}
}
