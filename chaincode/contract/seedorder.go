package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type SeedOrderContract struct {
	contractapi.Contract
}

type SeedOrderDoc struct {
	DocType string `json:"doc_type"`
	model.SeedOrder
}

func NewSeedOrderDoc(seedOrder model.SeedOrder) SeedOrderDoc {
	return SeedOrderDoc{
		DocType:   "seedorder",
		SeedOrder: seedOrder,
	}
}

func (c *SeedOrderContract) FindAllOutgoing(ctx contractapi.TransactionContextInterface, ordererId string) ([]*model.SeedOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"seedorder","orderer_id":"%s"}}`, ordererId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	seedOrders := make([]*model.SeedOrder, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var seedOrder model.SeedOrder
		err = json.Unmarshal(result.Value, &seedOrder)
		if err != nil {
			return nil, err
		}
		seedOrders = append(seedOrders, &seedOrder)
	}

	return seedOrders, nil
}

func (c *SeedOrderContract) FindAllIncoming(ctx contractapi.TransactionContextInterface, sellerID string) ([]*model.SeedOrder, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"seedorder","seller_id":"%s"}}`, sellerID)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	seedOrders := make([]*model.SeedOrder, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var seedOrder model.SeedOrder
		err = json.Unmarshal(result.Value, &seedOrder)
		if err != nil {
			return nil, err
		}
		seedOrders = append(seedOrders, &seedOrder)
	}

	return seedOrders, nil
}

func (c *SeedOrderContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.SeedOrder, error) {
	seedOrder, err := c.getSeedOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if seedOrder == nil {
		return nil, fmt.Errorf("the rice grain order %s does not exist", id)
	}

	return seedOrder, nil
}

func (c *SeedOrderContract) Create(ctx contractapi.TransactionContextInterface, id string, ordererId string, sellerId string, seedId string, riceGrainOrderId string, weight float32, orderedAt time.Time) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to create seed order, %w", err)
	}

	exists, err := c.seedExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the seed order %s already exists", id)
	}

	if err := c.processRiceGrainOrder(ctx, riceGrainOrderId, orderedAt); err != nil {
		return err
	}

	seedOrder := NewSeedOrderDoc(
		model.SeedOrder{
			ID:               id,
			OrdererID:        ordererId,
			SellerID:         sellerId,
			SeedID:           seedId,
			RiceGrainOrderID: riceGrainOrderId,
			Weight:           weight,
			Order:            model.NewOrder(orderedAt),
		},
	)
	seedOrderJSON, err := json.Marshal(seedOrder)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, seedOrderJSON)
}

func (c *SeedOrderContract) Accept(ctx contractapi.TransactionContextInterface, id string, acceptedAt time.Time) error {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to accept seed order, %w", err)
	}

	seedOrder, err := c.getSeedOrder(ctx, id)
	if err != nil {
		return err
	}
	if seedOrder == nil {
		return fmt.Errorf("the seed order %s does not exist", id)
	}

	seedOrderDoc := NewSeedOrderDoc(*seedOrder)
	seedOrderDoc.Accept(acceptedAt)
	seedOrderDoc.Process(acceptedAt)
	seedOrderDoc.Available(acceptedAt)
	seedOrderDocJSON, err := json.Marshal(seedOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, seedOrderDocJSON)
}

func (c *SeedOrderContract) Reject(ctx contractapi.TransactionContextInterface, id string, rejectedAt time.Time, reason string) error {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to reject seed order, %w", err)
	}

	seedOrder, err := c.getSeedOrder(ctx, id)
	if err != nil {
		return err
	}
	if seedOrder == nil {
		return fmt.Errorf("the seed order %s does not exist", id)
	}

	seedOrderDoc := NewSeedOrderDoc(*seedOrder)
	seedOrderDoc.Reject(rejectedAt, reason)
	seedOrderDocJSON, err := json.Marshal(seedOrderDoc)
	if err != nil {
		return err
	}

	if err := c.rejectRiceGrainOrder(ctx, seedOrder.RiceGrainOrderID, rejectedAt, reason); err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, seedOrderDocJSON)
}

func (c *SeedOrderContract) Ship(ctx contractapi.TransactionContextInterface, id string, shippedAt time.Time) error {
	err := c.authorizeRoleAsSupplier(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to ship seed order, %w", err)
	}

	seedOrder, err := c.getSeedOrder(ctx, id)
	if err != nil {
		return err
	}
	if seedOrder == nil {
		return fmt.Errorf("the seed order %s does not exist", id)
	}

	seedOrderDoc := NewSeedOrderDoc(*seedOrder)
	seedOrderDoc.Ship(shippedAt)
	seedOrderDocJSON, err := json.Marshal(seedOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, seedOrderDocJSON)
}

func (c *SeedOrderContract) Receive(ctx contractapi.TransactionContextInterface, id string, receivedAt time.Time) error {
	err := c.authorizeRoleAsProducer(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to receive seed order, %w", err)
	}

	seedOrder, err := c.getSeedOrder(ctx, id)
	if err != nil {
		return err
	}
	if seedOrder == nil {
		return fmt.Errorf("the seed order %s does not exist", id)
	}

	seedOrderDoc := NewSeedOrderDoc(*seedOrder)
	seedOrderDoc.Receive(receivedAt)
	seedOrderDocJSON, err := json.Marshal(seedOrderDoc)
	if err != nil {
		return err
	}

	if err := c.availableRiceGrainOrder(ctx, seedOrderDoc.RiceGrainOrderID, seedOrderDoc.ReceivedAt); err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, seedOrderDocJSON)
}

func (c *SeedOrderContract) authorizeRoleAsProducer(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "producer")
	if err != nil {
		return errors.New("only producer allowed")
	}

	return nil
}

func (c *SeedOrderContract) authorizeRoleAsSupplier(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetClientIdentity().AssertAttributeValue("hf.Affiliation", "supplier")
	if err != nil {
		return errors.New("only supplier allowed")
	}

	return nil
}

func (c *SeedOrderContract) getSeedOrder(ctx contractapi.TransactionContextInterface, id string) (*model.SeedOrder, error) {
	seedOrderJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if seedOrderJSON == nil {
		return nil, nil
	}

	var seedOrder model.SeedOrder
	err = json.Unmarshal(seedOrderJSON, &seedOrder)
	if err != nil {
		return nil, err
	}

	return &seedOrder, nil
}

func (c *SeedOrderContract) getRiceGrainOrder(ctx contractapi.TransactionContextInterface, id string) (*model.RiceGrainOrder, error) {
	riceGrainOrderJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
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

func (c *SeedOrderContract) getRiceOrder(ctx contractapi.TransactionContextInterface, id string) (*model.RiceOrder, error) {
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

func (c *SeedOrderContract) seedExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	riceGrainOrderJSON, err := c.getSeedOrder(ctx, id)
	if err != nil {
		return false, err
	}

	return riceGrainOrderJSON != nil, nil
}

func (c *SeedOrderContract) processRiceGrainOrder(ctx contractapi.TransactionContextInterface, id string, processingAt time.Time) error {
	riceGrainOrder, err := c.getRiceGrainOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceGrainOrder == nil {
		return fmt.Errorf("the rice grain order %s does not exist", id)
	}

	riceGrainOrderDoc := NewRiceGrainOrderDoc(*riceGrainOrder)
	riceGrainOrderDoc.Process(processingAt)
	riceGrainOrderDocJSON, err := json.Marshal(riceGrainOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainOrderDocJSON)
}

func (c *SeedOrderContract) rejectRiceGrainOrder(ctx contractapi.TransactionContextInterface, id string, rejectedAt time.Time, reason string) error {
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

	if err := c.rejectRiceOrder(ctx, riceGrainOrder.RiceOrderID, rejectedAt, reason); err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainOrderDocJSON)
}

func (c *SeedOrderContract) rejectRiceOrder(ctx contractapi.TransactionContextInterface, id string, rejectedAt time.Time, reason string) error {
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

func (c *SeedOrderContract) availableRiceGrainOrder(ctx contractapi.TransactionContextInterface, id string, availableAt time.Time) error {
	riceGrainOrder, err := c.getRiceGrainOrder(ctx, id)
	if err != nil {
		return err
	}
	if riceGrainOrder == nil {
		return fmt.Errorf("the rice grain order %s does not exist", id)
	}

	riceGrainOrderDoc := NewRiceGrainOrderDoc(*riceGrainOrder)
	riceGrainOrderDoc.Available(availableAt)
	riceGrainOrderDocJSON, err := json.Marshal(riceGrainOrderDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, riceGrainOrderDocJSON)
}
