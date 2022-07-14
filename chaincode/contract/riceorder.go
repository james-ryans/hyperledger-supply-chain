package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type RiceOrderContract struct {
	contractapi.Contract
}

type RiceOrderDoc struct {
	DocType string `json:"doc_type"`
	model.RiceOrder
}

func NewRiceOrderDoc(riceOrder model.RiceOrder) RiceOrderDoc {
	return RiceOrderDoc{
		DocType:   "riceorder",
		RiceOrder: riceOrder,
	}
}

func (c *RiceOrderContract) FindAllOutgoing(ctx contractapi.TransactionContextInterface, ordererId string) ([]*model.RiceOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"riceorder","orderer_id":"%s"}}`, ordererId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	riceOrders := make([]*model.RiceOrder, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var riceOrder model.RiceOrder
		err = json.Unmarshal(result.Value, &riceOrder)
		if err != nil {
			return nil, err
		}
		riceOrders = append(riceOrders, &riceOrder)
	}

	return riceOrders, nil
}

func (c *RiceOrderContract) FindAllIncoming(ctx contractapi.TransactionContextInterface, sellerId string) ([]*model.RiceOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"riceorder","seller_id":"%s"}}`, sellerId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	riceOrders := make([]*model.RiceOrder, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var riceOrder model.RiceOrder
		err = json.Unmarshal(result.Value, &riceOrder)
		if err != nil {
			return nil, err
		}
		riceOrders = append(riceOrders, &riceOrder)
	}

	return riceOrders, nil
}

func (c *RiceOrderContract) FindAllAcceptedIncoming(ctx contractapi.TransactionContextInterface, sellerId string) ([]*model.RiceOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"riceorder","seller_id":"%s","status":"accepted"}}`, sellerId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	riceOrders := make([]*model.RiceOrder, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var riceOrder model.RiceOrder
		err = json.Unmarshal(result.Value, &riceOrder)
		if err != nil {
			return nil, err
		}
		riceOrders = append(riceOrders, &riceOrder)
	}

	return riceOrders, nil
}

func (c *RiceOrderContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.RiceOrder, error) {
	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if riceOrder == nil {
		return nil, fmt.Errorf("the rice order %s does not exist", id)
	}

	return riceOrder, nil
}

func (c *RiceOrderContract) Create(ctx contractapi.TransactionContextInterface, id string, ordererId string, sellerId string, riceId string, quantity int32, orderedAt time.Time) error {
	err := c.authorizeRoleAsVendor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to create rice order, %w", err)
	}

	exists, err := c.riceOrderExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the rice order %s already exists", id)
	}

	riceOrderDoc := NewRiceOrderDoc(
		model.RiceOrder{
			ID:        id,
			OrdererID: ordererId,
			SellerID:  sellerId,
			RiceID:    riceId,
			Quantity:  quantity,
			Order:     model.NewOrder(orderedAt),
		},
	)
	riceOrderDocJSON, err := json.Marshal(riceOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}

func (c *RiceOrderContract) Accept(ctx contractapi.TransactionContextInterface, id string, acceptedAt time.Time) error {
	err := c.authorizeRoleAsManufacturerOrDistributor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to accept rice order, %w", err)
	}

	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceOrder == nil {
		return fmt.Errorf("the rice order %s does not exist", id)
	}

	riceOrderDoc := NewRiceOrderDoc(*riceOrder)
	riceOrderDoc.Accept(acceptedAt)
	riceOrderDocJSON, err := json.Marshal(riceOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}

func (c *RiceOrderContract) Reject(ctx contractapi.TransactionContextInterface, id string, rejectedAt time.Time, reason string) error {
	err := c.authorizeRoleAsManufacturerOrDistributor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to reject rice order, %w", err)
	}

	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceOrder == nil {
		return fmt.Errorf("the rice order %s does not exist", id)
	}

	riceOrderDoc := NewRiceOrderDoc(*riceOrder)
	riceOrderDoc.Reject(rejectedAt, reason)
	riceOrderDocJSON, err := json.Marshal(riceOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}

func (c *RiceOrderContract) Ship(ctx contractapi.TransactionContextInterface, id string, shippedAt time.Time, grade string, millingDate time.Time, storageTemperature float32, storageHumidity float32) error {
	err := c.authorizeRoleAsManufacturerOrDistributor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to ship rice order, %w", err)
	}

	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceOrder == nil {
		return fmt.Errorf("the rice order %s does not exist", id)
	}

	riceOrderDoc := NewRiceOrderDoc(*riceOrder)
	riceOrderDoc.Ship(shippedAt, grade, millingDate, storageTemperature, storageHumidity)
	riceOrderDocJSON, err := json.Marshal(riceOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}

func (c *RiceOrderContract) Receive(ctx contractapi.TransactionContextInterface, id string, receivedAt time.Time) error {
	err := c.authorizeRoleAsVendor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to receive rice order, %w", err)
	}

	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceOrder == nil {
		return fmt.Errorf("the rice order %s does not exist", id)
	}

	riceOrderDoc := NewRiceOrderDoc(*riceOrder)
	riceOrderDoc.Receive(receivedAt)
	riceOrderDocJSON, err := json.Marshal(riceOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}

func (c *RiceOrderContract) authorizeRoleAsVendor(ctx contractapi.TransactionContextInterface) error {
	role, found, err := ctx.GetClientIdentity().GetAttributeValue("hf.Affiliation")
	if err != nil || !found || (role != "distributor" && role != "retailer") {
		return errors.New("only vendor allowed")
	}

	return nil
}

func (c *RiceOrderContract) authorizeRoleAsManufacturerOrDistributor(ctx contractapi.TransactionContextInterface) error {
	role, found, err := ctx.GetClientIdentity().GetAttributeValue("hf.Affiliation")
	if err != nil || !found || (role != "manufacturer" && role != "distributor") {
		return errors.New("only manufacturer and distributor allowed")
	}

	return nil
}

func (c *RiceOrderContract) getRiceOrder(ctx contractapi.TransactionContextInterface, id string) (*model.RiceOrder, error) {
	riceOrderJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if riceOrderJSON == nil {
		return nil, nil
	}

	var riceOrder model.RiceOrder
	err = json.Unmarshal(riceOrderJSON, &riceOrder)
	if err != nil {
		return nil, err
	}

	return &riceOrder, nil
}

func (c *RiceOrderContract) riceOrderExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return false, err
	}

	return riceOrder != nil, nil
}
