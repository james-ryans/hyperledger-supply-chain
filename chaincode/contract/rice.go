package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type RiceContract struct {
	contractapi.Contract
}

type RiceDoc struct {
	DocType string `json:"doc_type"`
	model.Rice
}

func (c *RiceContract) CreateRice(ctx contractapi.TransactionContextInterface, id string, manufacturerId string, code string, brandName string, weight float32, texture string, amyloseRate float32) error {
	err := c.authorizeRoleAsManufacturer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.riceExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the rice %s already exists", id)
	}

	rice := RiceDoc{
		DocType: "rice",
		Rice: model.Rice{
			ID:             id,
			ManufacturerID: manufacturerId,
			Code:           code,
			BrandName:      brandName,
			Weight:         weight,
			Texture:        texture,
			AmyloseRate:    amyloseRate,
		},
	}
	riceJSON, err := json.Marshal(rice)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceJSON)
}

func (c *RiceContract) ReadRice(ctx contractapi.TransactionContextInterface, id string) (*model.Rice, error) {
	riceJSON, err := c.getRice(ctx, id)
	if err != nil {
		return nil, err
	}
	if riceJSON == nil {
		return nil, fmt.Errorf("the rice %s does not exist", id)
	}

	var rice model.Rice
	err = json.Unmarshal(riceJSON, &rice)
	if err != nil {
		return nil, err
	}

	return &rice, nil
}

func (c *RiceContract) UpdateRice(ctx contractapi.TransactionContextInterface, id string, manufacturerId string, code string, brandName string, weight float32, texture string, amyloseRate float32) error {
	err := c.authorizeRoleAsManufacturer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.riceExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the rice %s does not exist", id)
	}

	rice := RiceDoc{
		DocType: "rice",
		Rice: model.Rice{
			ID:             id,
			ManufacturerID: manufacturerId,
			Code:           code,
			BrandName:      brandName,
			Weight:         weight,
			Texture:        texture,
			AmyloseRate:    amyloseRate,
		},
	}
	riceJSON, err := json.Marshal(rice)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceJSON)
}

func (c *RiceContract) DeleteRice(ctx contractapi.TransactionContextInterface, id string) error {
	err := c.authorizeRoleAsManufacturer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.riceExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the rice %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (c *RiceContract) QueryRices(ctx contractapi.TransactionContextInterface, query string) ([]*model.Rice, error) {
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	rices := make([]*model.Rice, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var rice model.Rice
		err = json.Unmarshal(result.Value, &rice)
		if err != nil {
			return nil, err
		}
		rices = append(rices, &rice)
	}

	return rices, nil
}

func (c *RiceContract) authorizeRoleAsManufacturer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "manufacturer")
	if err != nil {
		return errors.New("you are not authorized to create rice asset, only manufacturer allowed")
	}

	return nil
}

func (c *RiceContract) getRice(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	riceJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return riceJSON, nil
}

func (c *RiceContract) riceExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	riceJSON, err := c.getRice(ctx, id)
	if err != nil {
		return false, err
	}

	return riceJSON != nil, nil
}
