package usermodel

import "github.com/meneketehe/hehe/app/model"

type Seed struct {
	Type               string   `json:"type"`
	VarietyName        string   `json:"variety_name"`
	PlantAge           float32  `json:"plant_age"`
	PlantShape         string   `json:"plant_shape"`
	PlantHeight        float32  `json:"plant_height"`
	LeafShape          string   `json:"leaf_shape"`
	StorageTemperature *float32 `json:"storage_temperature"`
	StorageHumidity    *float32 `json:"storage_humidity"`
}

func (s *Seed) SetType() {
	s.Type = "seed"
}

func (s *Seed) GetType() string {
	return "seed"
}

func FromSeedModel(seed *model.Seed, seedInstance *model.SeedInstance) *Seed {
	return &Seed{
		Type:               "seed",
		VarietyName:        seed.VarietyName,
		PlantAge:           seed.PlantAge,
		PlantShape:         seed.PlantShape,
		PlantHeight:        seed.PlantHeight,
		LeafShape:          seed.LeafShape,
		StorageTemperature: seedInstance.GetStorageTemperature(),
		StorageHumidity:    seedInstance.GetStorageHumidity(),
	}
}
