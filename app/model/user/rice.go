package usermodel

import (
	"time"

	"github.com/meneketehe/hehe/app/model"
)

type Rice struct {
	Type               string     `json:"type"`
	BrandName          string     `json:"brand_name"`
	Weight             float32    `json:"weight"`
	Texture            string     `json:"texture"`
	AmyloseRate        float32    `json:"amylose_rate"`
	Grade              *string    `json:"grade"`
	MillingDate        *time.Time `json:"milling_date"`
	StorageTemperature *float32   `json:"storage_temperature"`
	StorageHumidity    *float32   `json:"storage_humidity"`
}

func (r *Rice) SetType() {
	r.Type = "rice"
}

func (r *Rice) GetType() string {
	return "rice"
}

func FromRiceModel(rice *model.Rice, riceInstance *model.RiceInstance) *Rice {
	return &Rice{
		Type:               "rice",
		BrandName:          rice.BrandName,
		Weight:             rice.Weight,
		Texture:            rice.Texture,
		AmyloseRate:        rice.AmyloseRate,
		Grade:              riceInstance.GetGrade(),
		MillingDate:        riceInstance.GetMillingDate(),
		StorageTemperature: riceInstance.GetStorageTemperature(),
		StorageHumidity:    riceInstance.GetStorageHumidity(),
	}
}
