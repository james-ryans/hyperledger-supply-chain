package usermodel

import (
	"encoding/json"
	"fmt"

	"github.com/meneketehe/hehe/app/model"
)

type RiceSack struct {
	Code     string    `json:"code"`
	Comments []Comment `json:"comments"`
	Traces   []Trace   `json:"traces"`
}

type RiceSackService interface {
	GetRiceSackByCode(userID, code string) (*RiceSack, error)
	CreateRiceSack(sack *RiceSack) (*RiceSack, error)
	TraceRiceSack(channelID string, riceSack *model.RiceSack) (*RiceSack, error)
}

type RiceSackRepository interface {
	FindByCode(userID, code string) (*RiceSack, error)
	Create(sack *RiceSack) error
}

func UnmarshalRiceSack(riceSackJSON []byte) (*RiceSack, error) {
	riceSack := &RiceSack{}

	riceSackObject, err := unmarshalObject(riceSackJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json")
	}

	err = json.Unmarshal(riceSackObject["code"], &riceSack.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal code")
	}

	err = json.Unmarshal(riceSackObject["comments"], &riceSack.Comments)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal comments")
	}

	tracesArray, err := unmarshalArray(riceSackObject["traces"])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal traces")
	}

	traces := make([]Trace, 0)
	for index, traceJSON := range tracesArray {
		trace := Trace{}

		traceObject, err := unmarshalObject(traceJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal traces[%d]", index)
		}

		err = json.Unmarshal(traceObject["shipped_at"], &trace.ShippedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal traces[%d].shipped_at", index)
		}

		err = json.Unmarshal(traceObject["received_at"], &trace.ReceivedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal traces[%d].received_at", index)
		}

		err = json.Unmarshal(traceObject["organization"], &trace.Organization)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal traces[%d].organization", index)
		}

		commodityObject, err := unmarshalObject(traceObject["commodity"])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal traces[%d].commodity", index)
		}

		var commodityType string
		err = json.Unmarshal(commodityObject["type"], &commodityType)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal traces[%d].commodity.type", index)
		}

		switch commodityType {
		case "seed":
			var seed Seed
			err = json.Unmarshal(traceObject["commodity"], &seed)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal traces[%d].commodity", index)
			}
			trace.Commodity = &seed
		case "rice_grain":
			var riceGrain RiceGrain
			err = json.Unmarshal(traceObject["commodity"], &riceGrain)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal traces[%d].commodity", index)
			}
			trace.Commodity = &riceGrain
		case "rice":
			var rice Rice
			err = json.Unmarshal(traceObject["commodity"], &rice)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal traces[%d].commodity", index)
			}
			trace.Commodity = &rice
		}

		traces = append(traces, trace)
	}
	riceSack.Traces = traces

	return riceSack, nil
}

func unmarshalObject(objectJSON []byte) (map[string]json.RawMessage, error) {
	var obj map[string]json.RawMessage

	err := json.Unmarshal(objectJSON, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func unmarshalArray(arrayJSON []byte) ([]json.RawMessage, error) {
	var arr []json.RawMessage

	err := json.Unmarshal(arrayJSON, &arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}
