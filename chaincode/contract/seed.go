package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type SeedContract struct {
	contractapi.Contract
}

type SeedDoc struct {
	DocType string `json:"doc_type"`
	model.Seed
}

func (c *SeedContract) CreateSeed(ctx contractapi.TransactionContextInterface, id string, supplierId string, varietyName string, plantAge float32, plantShape string, plantHeight float32, leafShape string) error {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return err
	}

	exists, err := c.seedExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the seed %s already exists", id)
	}

	seed := SeedDoc{
		DocType: "seed",
		Seed: model.Seed{
			ID:          id,
			SupplierID:  supplierId,
			VarietyName: varietyName,
			PlantAge:    plantAge,
			PlantShape:  plantShape,
			PlantHeight: plantHeight,
			LeafShape:   leafShape,
		},
	}
	seedJSON, err := json.Marshal(seed)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, seedJSON)
}

func (c *SeedContract) ReadSeed(ctx contractapi.TransactionContextInterface, id string) (*model.Seed, error) {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return nil, err
	}

	seedJSON, err := c.getSeed(ctx, id)
	if err != nil {
		return nil, err
	}
	if seedJSON == nil {
		return nil, fmt.Errorf("the seed %s does not exist", id)
	}

	var seed model.Seed
	err = json.Unmarshal(seedJSON, &seed)
	if err != nil {
		return nil, err
	}

	return &seed, nil
}

func (c *SeedContract) UpdateSeed(ctx contractapi.TransactionContextInterface, id string, supplierId string, varietyName string, plantAge float32, plantShape string, plantHeight float32, leafShape string) error {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return err
	}

	exists, err := c.seedExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the seed %s does not exist", id)
	}

	seed := SeedDoc{
		DocType: "seed",
		Seed: model.Seed{
			ID:          id,
			SupplierID:  supplierId,
			VarietyName: varietyName,
			PlantAge:    plantAge,
			PlantShape:  plantShape,
			PlantHeight: plantHeight,
			LeafShape:   leafShape,
		},
	}
	seedJSON, err := json.Marshal(seed)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, seedJSON)
}

func (c *SeedContract) DeleteSeed(ctx contractapi.TransactionContextInterface, id string) error {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return err
	}

	exists, err := c.seedExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the seed %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (c *SeedContract) QuerySeeds(ctx contractapi.TransactionContextInterface, query string) ([]*model.Seed, error) {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return nil, err
	}

	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	seeds := make([]*model.Seed, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var seed model.Seed
		err = json.Unmarshal(result.Value, &seed)
		if err != nil {
			return nil, err
		}
		seeds = append(seeds, &seed)
	}

	return seeds, nil
}

func (c *SeedContract) authorizeRoleAsSupplier(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "supplier")
	if err != nil {
		return errors.New("you are not authorized to create seed asset, only supplier allowed")
	}

	return nil
}

func (c *SeedContract) getSeed(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	seedJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return seedJSON, nil
}

func (c *SeedContract) seedExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	seedJSON, err := c.getSeed(ctx, id)
	if err != nil {
		return false, err
	}

	return seedJSON != nil, nil
}
