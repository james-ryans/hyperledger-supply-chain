package model

type RiceGrain struct {
	ID          string `json:"id"`
	ProducerID  string `json:"producer_id"`
	VarietyName string `json:"variety_name"`
	GrainShape  string `json:"grain_shape"`
	GrainColor  string `json:"grain_color"`
}
