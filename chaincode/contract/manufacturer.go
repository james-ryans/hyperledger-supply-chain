package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ManufacturerContract struct {
	contractapi.Contract
}

type Manufacturer struct {
	organization
}

func (c *ManufacturerContract) CreateManufacturer(ctx contractapi.TransactionContextInterface, id string, name string, province string, city string, district string, postalCode string, address string, phone string, email string, latitude float32, longitude float32) error {
	err := c.authorizeRoleAsManufacturer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.manufacturerExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the manufacturer %s already exists", id)
	}

	manufacturer := Manufacturer{
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
	manufacturerJSON, err := json.Marshal(manufacturer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, manufacturerJSON)
}

func (c *ManufacturerContract) ReadManufacturer(ctx contractapi.TransactionContextInterface, id string) (*Manufacturer, error) {
	manufacturerJSON, err := c.getManufacturer(ctx, id)
	if err != nil {
		return nil, err
	}
	if manufacturerJSON == nil {
		return nil, fmt.Errorf("the manufacturer %s does not exist", id)
	}

	var manufacturer Manufacturer
	err = json.Unmarshal(manufacturerJSON, &manufacturer)
	if err != nil {
		return nil, err
	}

	return &manufacturer, nil
}

func (c *ManufacturerContract) authorizeRoleAsManufacturer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "manufacturer")
	if err != nil {
		return errors.New("you are not authorized to create manufacturer organization, only manufacturer allowed")
	}

	return nil
}

func (c *ManufacturerContract) getManufacturer(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	manufacturerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return manufacturerJSON, nil
}

func (c *ManufacturerContract) manufacturerExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	manufacturerJSON, err := c.getManufacturer(ctx, id)
	if err != nil {
		return false, err
	}

	return manufacturerJSON != nil, nil
}
