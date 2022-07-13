package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type RetailerContract struct {
	contractapi.Contract
}

type RetailerDoc struct {
	DocType string `json:"doc_type"`
	model.Retailer
}

func (c *RetailerContract) FindAll(ctx contractapi.TransactionContextInterface) ([]*model.Retailer, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"retailer"}}`)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	retailers := make([]*model.Retailer, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var retailer model.Retailer
		err = json.Unmarshal(result.Value, &retailer)
		if err != nil {
			return nil, err
		}
		retailers = append(retailers, &retailer)
	}

	return retailers, nil
}

func (c *RetailerContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.Retailer, error) {
	retailerJSON, err := c.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if retailerJSON == nil {
		return nil, fmt.Errorf("the retailer %s does not exist", id)
	}

	var retailer model.Retailer
	err = json.Unmarshal(retailerJSON, &retailer)
	if err != nil {
		return nil, err
	}

	return &retailer, nil
}

func (c *RetailerContract) Create(ctx contractapi.TransactionContextInterface, id string, orgType string, name string, province string, city string, district string, postalCode string, address string, phone string, email string, latitude float32, longitude float32) error {
	err := c.authorizeRoleAsRetailer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.exists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the retailer %s already exists", id)
	}

	retailer := RetailerDoc{
		DocType: "retailer",
		Retailer: model.Retailer{
			Vendor: model.Vendor{
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
		},
	}
	retailerJSON, err := json.Marshal(retailer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, retailerJSON)
}

func (c *RetailerContract) authorizeRoleAsRetailer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "retailer")
	if err != nil {
		return errors.New("you are not authorized to create retailer organization, only retailer allowed")
	}

	return nil
}

func (c *RetailerContract) get(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	retailerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return retailerJSON, nil
}

func (c *RetailerContract) exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	retailerJSON, err := c.get(ctx, id)
	if err != nil {
		return false, err
	}

	return retailerJSON != nil, nil
}
