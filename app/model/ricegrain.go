package model

type RiceGrain struct {
	ID          string `json:"id"`
	ProducerID  string `json:"producer_id"`
	VarietyName string `json:"variety_name"`
	GrainShape  string `json:"grain_shape"`
	GrainColor  string `json:"grain_color"`
}

type RiceGrainResponse struct {
	ID          string `json:"id"`
	VarietyName string `json:"variety_name"`
	GrainShape  string `json:"grain_shape"`
	GrainColor  string `json:"grain_color"`
}

func ToRiceGrainResponse(riceGrain *RiceGrain) RiceGrainResponse {
	return RiceGrainResponse{
		ID:          riceGrain.ID,
		VarietyName: riceGrain.VarietyName,
		GrainShape:  riceGrain.GrainShape,
		GrainColor:  riceGrain.GrainColor,
	}
}

func ToRiceGrainsResponse(riceGrains *[]RiceGrain) []RiceGrainResponse {
	res := make([]RiceGrainResponse, 0)
	for _, riceGrain := range *riceGrains {
		res = append(res, ToRiceGrainResponse(&riceGrain))
	}

	return res
}

type RiceGrainService interface {
	GetAllRiceGrains(channelID, ID string) (*[]RiceGrain, error)
	GetRiceGrainByID(channelID, ID string) (*RiceGrain, error)
	CreateRiceGrain(channelID string, riceGrain *RiceGrain) (*RiceGrain, error)
	UpdateRiceGrain(channelID string, riceGrain *RiceGrain) error
	DeleteRiceGrain(channelID, ID string) error
}

type RiceGrainRepository interface {
	FindAll(channelID, ID string) (*[]RiceGrain, error)
	FindByID(channelID, ID string) (*RiceGrain, error)
	Create(channelID string, riceGrain *RiceGrain) error
	Update(channelID string, riceGrain *RiceGrain) error
	Delete(channelID, ID string) error
}
