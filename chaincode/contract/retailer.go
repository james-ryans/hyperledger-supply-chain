package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type RetailerContract struct {
	contractapi.Contract
}

type Retailer struct {
	organization
}

func (c *RetailerContract) CreateRetailer(ctx contractapi.TransactionContextInterface, id string, name string, province string, city string, district string, postalCode string, address string, phone string, email string, latitude float32, longitude float32) error {
	err := c.authorizeRoleAsRetailer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.retailerExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the retailer %s already exists", id)
	}

	retailer := Retailer{
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
	retailerJSON, err := json.Marshal(retailer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, retailerJSON)
}

func (c *RetailerContract) ReadRetailer(ctx contractapi.TransactionContextInterface, id string) (*Retailer, error) {
	retailerJSON, err := c.getRetailer(ctx, id)
	if err != nil {
		return nil, err
	}
	if retailerJSON == nil {
		return nil, fmt.Errorf("the retailer %s does not exist", id)
	}

	var retailer Retailer
	err = json.Unmarshal(retailerJSON, &retailer)
	if err != nil {
		return nil, err
	}

	return &retailer, nil
}

func (c *RetailerContract) authorizeRoleAsRetailer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "retailer")
	if err != nil {
		return errors.New("you are not authorized to create retailer organization, only retailer allowed")
	}

	return nil
}

func (c *RetailerContract) getRetailer(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	retailerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return retailerJSON, nil
}

func (c *RetailerContract) retailerExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	retailerJSON, err := c.getRetailer(ctx, id)
	if err != nil {
		return false, err
	}

	return retailerJSON != nil, nil
}
