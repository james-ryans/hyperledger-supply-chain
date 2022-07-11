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

func (s *RiceGrainContract) CreateRiceGrain(ctx contractapi.TransactionContextInterface, id string, producerId string, varietyName string, grainShape string, grainColor string) error {
	err := s.authorizeRoleAsProducer(ctx)
	if err != nil {
		return err
	}

	exists, err := s.riceGrainExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the rice grain %s already exists", id)
	}

	riceGrain := model.RiceGrain{
		DocType:     "ricegrain",
		ID:          id,
		ProducerID:  producerId,
		VarietyName: varietyName,
		GrainShape:  grainShape,
		GrainColor:  grainColor,
	}
	riceGrainJSON, err := json.Marshal(riceGrain)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainJSON)
}

func (s *RiceGrainContract) ReadRiceGrain(ctx contractapi.TransactionContextInterface, id string) (*model.RiceGrain, error) {
	err := s.authorizeRoleAsProducer(ctx)
	if err != nil {
		return nil, err
	}

	riceGrainJSON, err := s.getRiceGrain(ctx, id)
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

func (s *RiceGrainContract) UpdateRiceGrain(ctx contractapi.TransactionContextInterface, id string, producerId string, varietyName string, grainShape string, grainColor string) error {
	err := s.authorizeRoleAsProducer(ctx)
	if err != nil {
		return err
	}

	exists, err := s.riceGrainExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the rice grain %s does not exist", id)
	}

	riceGrain := model.RiceGrain{
		DocType:     "ricegrain",
		ID:          id,
		ProducerID:  producerId,
		VarietyName: varietyName,
		GrainShape:  grainShape,
		GrainColor:  grainColor,
	}
	riceGrainJSON, err := json.Marshal(riceGrain)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainJSON)
}

func (s *RiceGrainContract) DeleteRiceGrain(ctx contractapi.TransactionContextInterface, id string) error {
	err := s.authorizeRoleAsProducer(ctx)
	if err != nil {
		return err
	}

	exists, err := s.riceGrainExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the rice grain %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *RiceGrainContract) authorizeRoleAsProducer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "producer")
	if err != nil {
		return errors.New("you are not authorized to create rice grain asset, only producer allowed")
	}

	return nil
}

func (s *RiceGrainContract) getRiceGrain(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	riceGrainJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return riceGrainJSON, nil
}

func (s *RiceGrainContract) riceGrainExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	riceGrainJSON, err := s.getRiceGrain(ctx, id)
	if err != nil {
		return false, err
	}

	return riceGrainJSON != nil, nil
}
