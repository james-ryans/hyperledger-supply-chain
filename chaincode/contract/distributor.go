package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type DistributorContract struct {
	contractapi.Contract
}

type DistributorDoc struct {
	DocType string `json:"doc_type"`
	model.Distributor
}

func (c *DistributorContract) FindAll(ctx contractapi.TransactionContextInterface) ([]*model.Distributor, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"distributor"}}`)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	distributors := make([]*model.Distributor, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var distributor model.Distributor
		err = json.Unmarshal(result.Value, &distributor)
		if err != nil {
			return nil, err
		}
		distributors = append(distributors, &distributor)
	}

	return distributors, nil
}

func (c *DistributorContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.Distributor, error) {
	distributorJSON, err := c.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if distributorJSON == nil {
		return nil, fmt.Errorf("the distributor %s does not exist", id)
	}

	var distributor model.Distributor
	err = json.Unmarshal(distributorJSON, &distributor)
	if err != nil {
		return nil, err
	}

	return &distributor, nil
}

func (c *DistributorContract) Create(ctx contractapi.TransactionContextInterface, id string, orgType string, name string, province string, city string, district string, postalCode string, address string, phone string, email string, latitude float32, longitude float32) error {
	err := c.authorizeRoleAsDistributor(ctx)
	if err != nil {
		return err
	}

	exists, err := c.exists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the distributor %s already exists", id)
	}

	distributor := DistributorDoc{
		DocType: "distributor",
		Distributor: model.Distributor{
			Vendor: model.Vendor{
				Organization: model.Organization{
					ID:   id,
					Name: name,
					Type: orgType,
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
	distributorJSON, err := json.Marshal(distributor)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, distributorJSON)
}

func (c *DistributorContract) authorizeRoleAsDistributor(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "distributor")
	if err != nil {
		return errors.New("you are not authorized to create distributor organization, only distributor allowed")
	}

	return nil
}

func (c *DistributorContract) get(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	distributorJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return distributorJSON, nil
}

func (c *DistributorContract) exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	distributorJSON, err := c.get(ctx, id)
	if err != nil {
		return false, err
	}

	return distributorJSON != nil, nil
}
