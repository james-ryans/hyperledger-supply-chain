package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type RiceGrainOrderContract struct {
	contractapi.Contract
}

type RiceGrainOrderDoc struct {
	DocType string `json:"doc_type"`
	model.RiceGrainOrder
}

func NewRiceGrainOrderDoc(riceGrainOrder model.RiceGrainOrder) RiceGrainOrderDoc {
	return RiceGrainOrderDoc{
		DocType:        "ricegrainorder",
		RiceGrainOrder: riceGrainOrder,
	}
}

func (c *RiceGrainOrderContract) FindAllOutgoing(ctx contractapi.TransactionContextInterface, ordererId string) ([]*model.RiceGrainOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"ricegrainorder","orderer_id":"%s"}}`, ordererId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	riceGrainOrders := make([]*model.RiceGrainOrder, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var riceGrainOrder model.RiceGrainOrder
		err = json.Unmarshal(result.Value, &riceGrainOrder)
		if err != nil {
			return nil, err
		}
		riceGrainOrders = append(riceGrainOrders, &riceGrainOrder)
	}

	return riceGrainOrders, nil
}

func (c *RiceGrainOrderContract) FindAllIncoming(ctx contractapi.TransactionContextInterface, sellerId string) ([]*model.RiceGrainOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"ricegrainorder","seller_id":"%s"}}`, sellerId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	riceGrainOrders := make([]*model.RiceGrainOrder, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var riceGrainOrder model.RiceGrainOrder
		err = json.Unmarshal(result.Value, &riceGrainOrder)
		if err != nil {
			return nil, err
		}
		riceGrainOrders = append(riceGrainOrders, &riceGrainOrder)
	}

	return riceGrainOrders, nil
}

func (c *RiceGrainOrderContract) FindAllAcceptedIncoming(ctx contractapi.TransactionContextInterface, sellerId string) ([]*model.RiceGrainOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"ricegrainorder","seller_id":"%s","status":"accepted"}}`, sellerId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	riceGrainOrders := make([]*model.RiceGrainOrder, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var riceGrainOrder model.RiceGrainOrder
		err = json.Unmarshal(result.Value, &riceGrainOrder)
		if err != nil {
			return nil, err
		}
		riceGrainOrders = append(riceGrainOrders, &riceGrainOrder)
	}

	return riceGrainOrders, nil
}

func (c *RiceGrainOrderContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.RiceGrainOrder, error) {
	riceGrainOrder, err := c.getRiceGrainOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if riceGrainOrder == nil {
		return nil, fmt.Errorf("the rice grain order %s does not exist", id)
	}

	return riceGrainOrder, nil
}

func (c *RiceGrainOrderContract) FindByRiceOrderID(ctx contractapi.TransactionContextInterface, id string) (*model.RiceGrainOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"ricegrainorder","rice_order_id":"%s"}}`, id)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	if !resultIterator.HasNext() {
		return nil, fmt.Errorf("the rice grain order does not exists")
	}

	riceGrainOrderJSON, err := resultIterator.Next()
	if err != nil {
		return nil, err
	}

	var riceGrainOrder *model.RiceGrainOrder
	err = json.Unmarshal(riceGrainOrderJSON.Value, &riceGrainOrder)
	if err != nil {
		return nil, err
	}

	return riceGrainOrder, nil
}

func (c *RiceGrainOrderContract) Create(ctx contractapi.TransactionContextInterface, id string, ordererId string, sellerId string, riceGrainId string, riceOrderId string, weight float32, orderedAt time.Time) error {
	err := c.authorizeRoleAsManufacturer(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to create rice grain order, %w", err)
	}

	exists, err := c.riceGrainOrderExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the rice grain order %s already exists", id)
	}

	if err := c.processRiceOrder(ctx, riceOrderId, orderedAt); err != nil {
		return err
	}

	riceGrainOrder := NewRiceGrainOrderDoc(
		model.RiceGrainOrder{
			ID:          id,
			OrdererID:   ordererId,
			SellerID:    sellerId,
			RiceGrainID: riceGrainId,
			RiceOrderID: riceOrderId,
			Weight:      weight,
			Order:       model.NewOrder(orderedAt),
		},
	)
	riceGrainOrderJSON, err := json.Marshal(riceGrainOrder)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainOrderJSON)
}

func (c *RiceGrainOrderContract) Accept(ctx contractapi.TransactionContextInterface, id string, acceptedAt time.Time) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to accept rice grain order, %w", err)
	}

	riceGrainOrder, err := c.getRiceGrainOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceGrainOrder == nil {
		return fmt.Errorf("the rice grain order %s does not exist", id)
	}

	riceGrainOrderDoc := NewRiceGrainOrderDoc(*riceGrainOrder)
	riceGrainOrderDoc.Accept(acceptedAt)
	riceGrainOrderDocJSON, err := json.Marshal(riceGrainOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainOrderDocJSON)
}

func (c *RiceGrainOrderContract) Reject(ctx contractapi.TransactionContextInterface, id string, rejectedAt time.Time, reason string) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to reject rice grain order, %w", err)
	}

	riceGrainOrder, err := c.getRiceGrainOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceGrainOrder == nil {
		return fmt.Errorf("the rice grain order %s does not exist", id)
	}

	riceGrainOrderDoc := NewRiceGrainOrderDoc(*riceGrainOrder)
	riceGrainOrderDoc.Reject(rejectedAt, reason)
	riceGrainOrderDocJSON, err := json.Marshal(riceGrainOrderDoc)
	if err != nil {
		return err
	}

	if err := c.orderRiceOrder(ctx, riceGrainOrder.RiceOrderID); err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainOrderDocJSON)
}

func (c *RiceGrainOrderContract) Ship(ctx contractapi.TransactionContextInterface, id string, shippedAt time.Time, plowMethod, sowMethod, irrigation, fertilization string, plantDate, harvestDate time.Time, storageTemperature, storageHumidity float32) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to ship rice grain order, %w", err)
	}

	riceGrainOrder, err := c.getRiceGrainOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceGrainOrder == nil {
		return fmt.Errorf("the rice grain order %s does not exist", id)
	}

	riceGrainOrderDoc := NewRiceGrainOrderDoc(*riceGrainOrder)
	riceGrainOrderDoc.Ship(shippedAt, plowMethod, sowMethod, irrigation, fertilization, plantDate, harvestDate, storageTemperature, storageHumidity)
	riceGrainOrderDocJSON, err := json.Marshal(riceGrainOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainOrderDocJSON)
}

func (c *RiceGrainOrderContract) Receive(ctx contractapi.TransactionContextInterface, id string, receivedAt time.Time) error {
	err := c.authorizeRoleAsManufacturer(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to receive rice grain order, %w", err)
	}

	riceGrainOrder, err := c.getRiceGrainOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceGrainOrder == nil {
		return fmt.Errorf("the rice grain order %s does not exist", id)
	}

	riceGrainOrderDoc := NewRiceGrainOrderDoc(*riceGrainOrder)
	riceGrainOrderDoc.Receive(receivedAt)
	riceGrainOrderDocJSON, err := json.Marshal(riceGrainOrderDoc)
	if err != nil {
		return err
	}

	if err := c.availableRiceOrder(ctx, riceGrainOrderDoc.RiceOrderID, riceGrainOrderDoc.ReceivedAt); err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainOrderDocJSON)
}

func (c *RiceGrainOrderContract) authorizeRoleAsManufacturer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "manufacturer")
	if err != nil {
		return errors.New("only manufacturer allowed")
	}

	return nil
}

func (c *RiceGrainOrderContract) authorizeRoleAsProducer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "producer")
	if err != nil {
		return errors.New("only manufacturer allowed")
	}

	return nil
}

func (c *RiceGrainOrderContract) getRiceGrainOrder(ctx contractapi.TransactionContextInterface, id string) (*model.RiceGrainOrder, error) {
	riceGrainOrderJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if riceGrainOrderJSON == nil {
		return nil, nil
	}

	var riceGrainOrder model.RiceGrainOrder
	err = json.Unmarshal(riceGrainOrderJSON, &riceGrainOrder)
	if err != nil {
		return nil, err
	}

	return &riceGrainOrder, nil
}

func (c *RiceGrainOrderContract) getRiceOrder(ctx contractapi.TransactionContextInterface, id string) (*model.RiceOrder, error) {
	riceOrderJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
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

func (c *RiceGrainOrderContract) riceGrainOrderExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	riceGrainOrder, err := c.getRiceGrainOrder(ctx, id)
	if err != nil {
		return false, err
	}

	return riceGrainOrder != nil, nil
}

func (c *RiceGrainOrderContract) orderRiceOrder(ctx contractapi.TransactionContextInterface, id string) error {
	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceOrder == nil {
		return fmt.Errorf("the rice order %s does not exist", id)
	}

	riceOrderDoc := NewRiceOrderDoc(*riceOrder)
	riceOrderDoc.Order = model.NewOrder(riceOrderDoc.OrderedAt)
	riceOrderDocJSON, err := json.Marshal(riceOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}

func (c *RiceGrainOrderContract) processRiceOrder(ctx contractapi.TransactionContextInterface, id string, processingAt time.Time) error {
	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceOrder == nil {
		return fmt.Errorf("the rice order %s does not exist", id)
	}

	riceOrderDoc := NewRiceOrderDoc(*riceOrder)
	riceOrderDoc.Process(processingAt)
	riceOrderDocJSON, err := json.Marshal(riceOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}

func (c *RiceGrainOrderContract) availableRiceOrder(ctx contractapi.TransactionContextInterface, id string, availableAt time.Time) error {
	riceOrder, err := c.getRiceOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceOrder == nil {
		return fmt.Errorf("the rice order %s does not exist", id)
	}

	riceOrderDoc := NewRiceOrderDoc(*riceOrder)
	riceOrderDoc.Available(availableAt)
	riceOrderDocJSON, err := json.Marshal(riceOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}
