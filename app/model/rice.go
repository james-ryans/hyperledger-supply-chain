package model

type Rice struct {
	ID             string  `json:"id"`
	ManufacturerID string  `json:"manufacturer_id"`
	Code           string  `json:"code"`
	BrandName      string  `json:"brand_name"`
	Weight         float32 `json:"weight"`
	Texture        string  `json:"texture"`
	AmyloseRate    float32 `json:"amylose_rate"`
}

type RiceResponse struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	BrandName   string  `json:"brand_name"`
	Weight      float32 `json:"weight"`
	Texture     string  `json:"texture"`
	AmyloseRate float32 `json:"amylose_rate"`
}

func ToRiceResponse(rice *Rice) RiceResponse {
	return RiceResponse{
		ID:          rice.ID,
		Code:        rice.Code,
		BrandName:   rice.BrandName,
		Weight:      rice.Weight,
		Texture:     rice.Texture,
		AmyloseRate: rice.AmyloseRate,
	}
}

func ToRicesResponse(rices *[]Rice) []RiceResponse {
	res := make([]RiceResponse, 0)
	for _, rice := range *rices {
		res = append(res, ToRiceResponse(&rice))
	}

	return res
}

type RiceService interface {
	GetAllRices(channelID, ID string) (*[]Rice, error)
	GetRiceByID(channelID, ID string) (*Rice, error)
	CreateRice(channelID string, rice *Rice) (*Rice, error)
	UpdateRice(channelID string, rice *Rice) error
	DeleteRice(channelID, ID string) error
}

type RiceRepository interface {
	FindAll(channelID, ID string) (*[]Rice, error)
	FindByID(channelID, ID string) (*Rice, error)
	Create(channelID string, rice *Rice) error
	Update(channelID string, rice *Rice) error
	Delete(channelID, ID string) error
}
