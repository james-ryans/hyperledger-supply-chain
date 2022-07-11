package contract

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type ProducerContract struct {
	contractapi.Contract
}

type ProducerDoc struct {
	DocType string `json:"doc_type"`
	model.Producer
}

func (c *ProducerContract) CreateProducer(ctx contractapi.TransactionContextInterface, id string, orgType string, name string, province string, city string, district string, postalCode string, address string, phone string, email string, latitude float32, longitude float32) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return err
	}

	exists, err := c.producerExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the producer %s already exists", id)
	}

	producer := ProducerDoc{
		DocType: "producer",
		Producer: model.Producer{
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
	producerJSON, err := json.Marshal(producer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, producerJSON)
}

func (c *ProducerContract) ReadProducer(ctx contractapi.TransactionContextInterface, id string) (*model.Producer, error) {
	producerJSON, err := c.getProducer(ctx, id)
	if err != nil {
		return nil, err
	}
	if producerJSON == nil {
		return nil, fmt.Errorf("the producer %s does not exist", id)
	}

	var producer model.Producer
	err = json.Unmarshal(producerJSON, &producer)
	if err != nil {
		return nil, err
	}

	return &producer, nil
}

func (c *ProducerContract) authorizeRoleAsProducer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "producer")
	if err != nil {
		return errors.New("you are not authorized to create producer organization, only producer allowed")
	}

	return nil
}

func (c *ProducerContract) getProducer(ctx contractapi.TransactionContextInterface, id string) ([]byte, error) {
	producerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	return producerJSON, nil
}

func (c *ProducerContract) producerExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	producerJSON, err := c.getProducer(ctx, id)
	if err != nil {
		return false, err
	}

	return producerJSON != nil, nil
}
