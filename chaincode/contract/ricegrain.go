package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type RiceGrainContract struct {
	contractapi.Contract
}

type RiceGrainDoc struct {
	DocType string `json:"doc_type"`
	model.RiceGrain
}

func (c *RiceGrainContract) CreateRiceGrain(ctx contractapi.TransactionContextInterface, id string, producerId string, varietyName string, grainShape string, grainColor string) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.riceGrainExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the rice grain %s already exists", id)
	}

	riceGrain := RiceGrainDoc{
		DocType: "ricegrain",
		RiceGrain: model.RiceGrain{
			ID:          id,
			ProducerID:  producerId,
			VarietyName: varietyName,
			GrainShape:  grainShape,
			GrainColor:  grainColor,
		},
	}
	riceGrainJSON, err := json.Marshal(riceGrain)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainJSON)
}

func (c *RiceGrainContract) ReadRiceGrain(ctx contractapi.TransactionContextInterface, id string) (*model.RiceGrain, error) {
	riceGrainJSON, err := c.getRiceGrain(ctx, id)
	if err != nil {
		return nil, err
	}
	if riceGrainJSON == nil {
		return nil, fmt.Errorf("the rice grain %s does not exist", id)
	}

	var riceGrain model.RiceGrain
	err = json.Unmarshal(riceGrainJSON, &riceGrain)
	if err != nil {
		return nil, err
	}

	return &riceGrain, nil
}

func (c *RiceGrainContract) UpdateRiceGrain(ctx contractapi.TransactionContextInterface, id string, producerId string, varietyName string, grainShape string, grainColor string) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.riceGrainExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the rice grain %s does not exist", id)
	}

	riceGrain := RiceGrainDoc{
		DocType: "ricegrain",
		RiceGrain: model.RiceGrain{
			ID:          id,
			ProducerID:  producerId,
			VarietyName: varietyName,
			GrainShape:  grainShape,
			GrainColor:  grainColor,
		},
	}
	riceGrainJSON, err := json.Marshal(riceGrain)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainJSON)
}

func (c *RiceGrainContract) DeleteRiceGrain(ctx contractapi.TransactionContextInterface, id string) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.riceGrainExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the rice grain %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (c *RiceGrainContract) QueryRiceGrains(ctx contractapi.TransactionContextInterface, query string) ([]*model.RiceGrain, error) {
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	riceGrains := make([]*model.RiceGrain, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var riceGrain model.RiceGrain
		err = json.Unmarshal(result.Value, &riceGrain)
		if err != nil {
			return nil, err
		}
		riceGrains = append(riceGrains, &riceGrain)
	}

	return riceGrains, nil
}

func (c *RiceGrainContract) authorizeRoleAsProducer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "producer")
	if err != nil {
		return errors.New("you are not authorized to create rice grain asset, only producer allowed")
	}

	return nil
}

func (c *RiceGrainContract) getRiceGrain(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	riceGrainJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return riceGrainJSON, nil
}

func (c *RiceGrainContract) riceGrainExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	riceGrainJSON, err := c.getRiceGrain(ctx, id)
	if err != nil {
		return false, err
	}

	return riceGrainJSON != nil, nil
}
