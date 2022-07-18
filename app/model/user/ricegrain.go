package usermodel

import (
	"time"

	"github.com/meneketehe/hehe/app/model"
)

type RiceGrain struct {
	Type               string     `json:"type"`
	VarietyName        string     `json:"variety_name"`
	GrainShape         string     `json:"grain_shape"`
	GrainColor         string     `json:"grain_color"`
	PlowMethod         *string    `json:"plow_method"`
	SowMethod          *string    `json:"sow_method"`
	Irrigation         *string    `json:"irrigation"`
	Fertilization      *string    `json:"fertilization"`
	PlantDate          *time.Time `json:"plant_date"`
	HarvestDate        *time.Time `json:"harvest_date"`
	StorageTemperature *float32   `json:"storage_temperature"`
	StorageHumidity    *float32   `json:"storage_humidity"`
}

func (r *RiceGrain) SetType() {
	r.Type = "rice_grain"
}

func (r *RiceGrain) GetType() string {
	return "rice_grain"
}

func FromRiceGrainModel(riceGrain *model.RiceGrain, riceGrainInstance *model.RiceGrainInstance) *RiceGrain {
	return &RiceGrain{
		Type:               "rice_grain",
		VarietyName:        riceGrain.VarietyName,
		GrainShape:         riceGrain.GrainShape,
		GrainColor:         riceGrain.GrainColor,
		PlowMethod:         riceGrainInstance.GetPlowMethod(),
		SowMethod:          riceGrainInstance.GetSowMethod(),
		Irrigation:         riceGrainInstance.GetIrrigation(),
		Fertilization:      riceGrainInstance.GetFertilization(),
		PlantDate:          riceGrainInstance.GetPlantDate(),
		HarvestDate:        riceGrainInstance.GetHarvestDate(),
		StorageTemperature: riceGrainInstance.GetStorageTemperature(),
		StorageHumidity:    riceGrainInstance.GetStorageHumidity(),
	}
}
