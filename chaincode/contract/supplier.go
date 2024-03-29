package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type SupplierContract struct {
	contractapi.Contract
}

type SupplierDoc struct {
	DocType string `json:"doc_type"`
	model.Supplier
}

func (c *SupplierContract) FindAll(ctx contractapi.TransactionContextInterface) ([]*model.Supplier, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"supplier"}}`)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	suppliers := make([]*model.Supplier, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var supplier model.Supplier
		err = json.Unmarshal(result.Value, &supplier)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, &supplier)
	}

	return suppliers, nil
}

func (c *SupplierContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.Supplier, error) {
	supplierJSON, err := c.get(ctx, id)
	if err != nil {
		return nil, err
	}
	if supplierJSON == nil {
		return nil, fmt.Errorf("the supplier %s does not exist", id)
	}

	var supplier model.Supplier
	err = json.Unmarshal(supplierJSON, &supplier)
	if err != nil {
		return nil, err
	}

	return &supplier, nil
}

func (c *SupplierContract) Create(ctx contractapi.TransactionContextInterface, id string, orgType string, name string, province string, city string, district string, postalCode string, address string, phone string, email string, latitude float32, longitude float32) error {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return err
	}

	exists, err := c.exists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the supplier %s already exists", id)
	}

	supplier := SupplierDoc{
		DocType: "supplier",
		Supplier: model.Supplier{
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
	supplierJSON, err := json.Marshal(supplier)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, supplierJSON)
}

func (c *SupplierContract) authorizeRoleAsSupplier(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "supplier")
	if err != nil {
		return errors.New("you are not authorized to create supplier organization, only supplier allowed")
	}

	return nil
}

func (c *SupplierContract) get(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	supplierJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return supplierJSON, nil
}

func (c *SupplierContract) exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	supplierJSON, err := c.get(ctx, id)
	if err != nil {
		return false, err
	}

	return supplierJSON != nil, nil
}
