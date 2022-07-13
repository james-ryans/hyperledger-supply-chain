package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type ManufacturerContract struct {
	contractapi.Contract
}

type ManufacturerDoc struct {
	DocType string `json:"doc_type"`
	model.Manufacturer
}

func (c *ManufacturerContract) FindAll(ctx contractapi.TransactionContextInterface) ([]*model.Manufacturer, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"manufacturer"}}`)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	manufacturers := make([]*model.Manufacturer, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var manufacturer model.Manufacturer
		err = json.Unmarshal(result.Value, &manufacturer)
		if err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, &manufacturer)
	}

	return manufacturers, nil
}

func (c *ManufacturerContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.Manufacturer, error) {
	manufacturerJSON, err := c.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if manufacturerJSON == nil {
		return nil, fmt.Errorf("the manufacturer %s does not exist", id)
	}

	var manufacturer model.Manufacturer
	err = json.Unmarshal(manufacturerJSON, &manufacturer)
	if err != nil {
		return nil, err
	}

	return &manufacturer, nil
}

func (c *ManufacturerContract) Create(ctx contractapi.TransactionContextInterface, id string, orgType string, name string, province string, city string, district string, postalCode string, address string, phone string, email string, latitude float32, longitude float32) error {
	err := c.authorizeRoleAsManufacturer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.exists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the manufacturer %s already exists", id)
	}

	manufacturer := ManufacturerDoc{
		DocType: "manufacturer",
		Manufacturer: model.Manufacturer{
			Organization: model.Organization{
				ID:   id,
				Type: orgType,
				Name: name,
				Location: model.Location{
					Province:   province,
					City:       city,
					District:   district,
					PostalCode: postalCode,
					Address:    address,
					Coordinate: model.Coordinate{
						Latitude:  latitude,
						Longitude: longitude,
					},
				},
				ContactInfo: model.ContactInfo{
					Phone: phone,
					Email: email,
				},
			},
		},
	}

	manufacturerJSON, err := json.Marshal(manufacturer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, manufacturerJSON)
}

func (c *ManufacturerContract) authorizeRoleAsManufacturer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "manufacturer")
	if err != nil {
		return errors.New("you are not authorized to create manufacturer organization, only manufacturer allowed")
	}

	return nil
}

func (c *ManufacturerContract) get(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	manufacturerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return manufacturerJSON, nil
}

func (c *ManufacturerContract) exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	manufacturerJSON, err := c.get(ctx, id)
	if err != nil {
		return false, err
	}

	return manufacturerJSON != nil, nil
}
