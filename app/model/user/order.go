package usermodel

import (
	"time"

	"github.com/meneketehe/hehe/app/model"
)

type Order struct {
	ReceivedAt   *time.Time   `json:"received_at"`
	ShippedAt    *time.Time   `json:"shipped_at"`
	Organization Organization `json:"organization"`
	Commodity    Commodity    `json:"commodity"`
}

func OrderFromRiceOrgModel(shippedAt *time.Time, receivedAt *time.Time, org *model.Organization, rice *model.Rice, riceInstance *model.RiceInstance) Order {
	return Order{
		ReceivedAt:   receivedAt,
		ShippedAt:    shippedAt,
		Organization: FromOrgModel(org),
		Commodity:    FromRiceModel(rice, riceInstance),
	}
}

func OrderFromRiceGrainOrgModel(shippedAt *time.Time, receivedAt *time.Time, org *model.Organization, riceGrain *model.RiceGrain, riceGrainInstance *model.RiceGrainInstance) Order {
	return Order{
		ReceivedAt:   receivedAt,
		ShippedAt:    shippedAt,
		Organization: FromOrgModel(org),
		Commodity:    FromRiceGrainModel(riceGrain, riceGrainInstance),
	}
}

func OrderFromSeedOrgModel(shippedAt *time.Time, receivedAt *time.Time, org *model.Organization, seed *model.Seed, seedInstance *model.SeedInstance) Order {
	return Order{
		ReceivedAt:   receivedAt,
		ShippedAt:    shippedAt,
		Organization: FromOrgModel(org),
		Commodity:    FromSeedModel(seed, seedInstance),
	}
}
