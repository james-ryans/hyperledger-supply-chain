package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type DistributorContract struct {
	contractapi.Contract
}

type Distributor struct {
	organization
}

func (c *DistributorContract) CreateDistributor(ctx contractapi.TransactionContextInterface, id string, name string, province string, city string, district string, postalCode string, address string, phone string, email string, latitude float32, longitude float32) error {
	err := c.authorizeRoleAsDistributor(ctx)
	if err != nil {
		return err
	}

	exists, err := c.distributorExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the distributor %s already exists", id)
	}

	distributor := Distributor{
		organization: organization{
			ID:   id,
			Name: name,
			Location: location{
				Province:   province,
				City:       city,
				District:   district,
				PostalCode: postalCode,
				Address:    address,
				Coordinate: coordinate{
					Latitude:  latitude,
					Longitude: longitude,
				},
			},
			ContactInfo: contactInfo{
				Phone: phone,
				Email: email,
			},
		},
	}
	distributorJSON, err := json.Marshal(distributor)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, distributorJSON)
}

func (c *DistributorContract) ReadDistributor(ctx contractapi.TransactionContextInterface, id string) (*Distributor, error) {
	distributorJSON, err := c.getDistributor(ctx, id)
	if err != nil {
		return nil, err
	}
	if distributorJSON == nil {
		return nil, fmt.Errorf("the distributor %s does not exist", id)
	}

	var distributor Distributor
	err = json.Unmarshal(distributorJSON, &distributor)
	if err != nil {
		return nil, err
	}

	return &distributor, nil
}

func (c *DistributorContract) authorizeRoleAsDistributor(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "distributor")
	if err != nil {
		return errors.New("you are not authorized to create distributor organization, only distributor allowed")
	}

	return nil
}

func (c *DistributorContract) getDistributor(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	distributorJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return distributorJSON, nil
}

func (c *DistributorContract) distributorExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	distributorJSON, err := c.getDistributor(ctx, id)
	if err != nil {
		return false, err
	}

	return distributorJSON != nil, nil
}
