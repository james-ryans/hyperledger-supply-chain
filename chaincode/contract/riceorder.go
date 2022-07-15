package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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

type ProductCounterDoc struct {
	DocType        string `json:"doc_type"`
	ID             string `json:"id"`
	ManufacturerID string `json:"manufacturer_id"`
	Count          int32  `json:"count"`
}

func NewRiceOrderDoc(riceOrder model.RiceOrder) RiceOrderDoc {
	return RiceOrderDoc{
		DocType:   "riceorder",
		RiceOrder: riceOrder,
	}
}

func NewProductCounterDoc(id, manufacturerId string, count int32) ProductCounterDoc {
	return ProductCounterDoc{
		DocType:        "productcounter",
		ID:             id,
		ManufacturerID: manufacturerId,
		Count:          count,
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
	riceOrder, err := getRiceOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if riceOrder == nil {
		return nil, fmt.Errorf("the rice order %s does not exist", id)
	}

	return riceOrder, nil
}

func (c *RiceOrderContract) Create(ctx contractapi.TransactionContextInterface, id string, ordererId string, sellerId string, riceId string, quantity int32, orderedAt time.Time) error {
	err := authorizeRoleAsVendor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to create rice order, %w", err)
	}

	exists, err := riceOrderExists(ctx, id)
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
	err := authorizeRoleAsManufacturerOrDistributor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to accept rice order, %w", err)
	}

	riceOrder, err := getRiceOrder(ctx, id)
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
	err := authorizeRoleAsManufacturerOrDistributor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to reject rice order, %w", err)
	}

	riceOrder, err := getRiceOrder(ctx, id)
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
	err := authorizeRoleAsManufacturerOrDistributor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to ship rice order, %w", err)
	}

	riceOrder, err := getRiceOrder(ctx, id)
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
	err := authorizeRoleAsVendor(ctx)
	if err != nil {
		return fmt.Errorf("you are not authorized to receive rice order, %w", err)
	}

	riceOrder, err := getRiceOrder(ctx, id)
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

	if role, found, _ := getRole(ctx); found && role == "distributor" {
		err = addRiceStock(ctx, riceOrder)
		if err != nil {
			return err
		}
	}

	if role, found, _ := getRole(ctx); found && role == "retailer" {
		err = moveRiceStock(ctx, riceOrder)
		if err != nil {
			return err
		}
	}

	return ctx.GetStub().PutState(id, riceOrderDocJSON)
}

func newDeterministicUuid(input string) (string, error) {
	id, err := uuid.NewRandomFromReader(strings.NewReader(fmt.Sprintf("%16s", input)))
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func getRole(ctx contractapi.TransactionContextInterface) (string, bool, error) {
	return ctx.GetClientIdentity().GetAttributeValue("hf.Affiliation")
}

func authorizeRoleAsVendor(ctx contractapi.TransactionContextInterface) error {
	role, found, err := ctx.GetClientIdentity().GetAttributeValue("hf.Affiliation")
	if err != nil || !found || (role != "distributor" && role != "retailer") {
		return errors.New("only vendor allowed")
	}

	return nil
}

func authorizeRoleAsManufacturerOrDistributor(ctx contractapi.TransactionContextInterface) error {
	role, found, err := ctx.GetClientIdentity().GetAttributeValue("hf.Affiliation")
	if err != nil || !found || (role != "manufacturer" && role != "distributor") {
		return errors.New("only manufacturer and distributor allowed")
	}

	return nil
}

func getManufacturer(ctx contractapi.TransactionContextInterface, id string) (*model.Manufacturer, error) {
	manufacturerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if manufacturerJSON == nil {
		return nil, nil
	}

	var manufacturer model.Manufacturer
	err = json.Unmarshal(manufacturerJSON, &manufacturer)
	if err != nil {
		return nil, err
	}

	return &manufacturer, nil
}

func getRiceOrder(ctx contractapi.TransactionContextInterface, id string) (*model.RiceOrder, error) {
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

func riceOrderExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	riceOrder, err := getRiceOrder(ctx, id)
	if err != nil {
		return false, err
	}

	return riceOrder != nil, nil
}

func getRiceStockpileByVendorIdAndRiceId(ctx contractapi.TransactionContextInterface, vendorId, riceId string) (*model.RiceStockpile, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"ricestockpile","vendor_id":"%s","rice_id":"%s"},"limit":1}`, vendorId, riceId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	if !resultIterator.HasNext() {
		return nil, nil
	}

	result, err := resultIterator.Next()
	if err != nil {
		return nil, err
	}

	var riceStockpile model.RiceStockpile
	err = json.Unmarshal(result.Value, &riceStockpile)
	if err != nil {
		return nil, err
	}

	return &riceStockpile, nil
}

func getManufacturerProductCounter(ctx contractapi.TransactionContextInterface, id string) (*ProductCounterDoc, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"productcounter","manufacturer_id":"%s"},"limit":1}`, id)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	if !resultIterator.HasNext() {
		newUuid, err := newDeterministicUuid(fmt.Sprintf("productcounter:%s", ctx.GetStub().GetTxID()))
		if err != nil {
			return nil, err
		}

		productCounterDoc := NewProductCounterDoc(newUuid, id, 0)
		return &productCounterDoc, nil
	}

	result, err := resultIterator.Next()
	if err != nil {
		return nil, err
	}

	var productCounterDoc ProductCounterDoc
	err = json.Unmarshal(result.Value, &productCounterDoc)
	if err != nil {
		return nil, err
	}

	return &productCounterDoc, nil
}

func putManufacturerProductCounter(ctx contractapi.TransactionContextInterface, doc *ProductCounterDoc) error {
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(doc.ID, docJSON)
}

func addRiceStock(ctx contractapi.TransactionContextInterface, riceOrder *model.RiceOrder) error {
	manufacturer, err := getManufacturer(ctx, riceOrder.SellerID)
	if err != nil {
		return err
	}

	riceStockpile, err := getRiceStockpileByVendorIdAndRiceId(ctx, riceOrder.OrdererID, riceOrder.RiceID)
	if err != nil {
		return err
	}

	if riceStockpile == nil {
		newUuid, err := newDeterministicUuid(fmt.Sprintf("addricestock:%s", ctx.GetStub().GetTxID()))
		if err != nil {
			return err
		}

		riceStockpile = &model.RiceStockpile{
			ID:       newUuid,
			RiceID:   riceOrder.RiceID,
			VendorID: riceOrder.OrdererID,
			Stock:    riceOrder.Quantity,
		}
	} else {
		riceStockpile.AddStock(riceOrder.Quantity)
	}

	riceStockpileDoc := NewRiceStockpileDoc(*riceStockpile)
	riceStockpileDocJSON, err := json.Marshal(riceStockpileDoc)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(riceStockpileDoc.ID, riceStockpileDocJSON)
	if err != nil {
		return err
	}

	err = createRiceSack(ctx, manufacturer, riceOrder.ID, riceStockpile.ID, riceOrder.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func moveRiceStock(ctx contractapi.TransactionContextInterface, riceOrder *model.RiceOrder) error {
	sourceStock, err := getRiceStockpileByVendorIdAndRiceId(ctx, riceOrder.SellerID, riceOrder.RiceID)
	if err != nil {
		return err
	}
	if sourceStock == nil {
		return errors.New("out of stock")
	}

	targetStock, err := getRiceStockpileByVendorIdAndRiceId(ctx, riceOrder.OrdererID, riceOrder.RiceID)
	if err != nil {
		return err
	}

	if targetStock == nil {
		newUuid, err := newDeterministicUuid(fmt.Sprintf("movericestock:%s", ctx.GetStub().GetTxID()))
		if err != nil {
			return err
		}

		targetStock = &model.RiceStockpile{
			ID:       newUuid,
			RiceID:   riceOrder.RiceID,
			VendorID: riceOrder.OrdererID,
			Stock:    riceOrder.Quantity,
		}
	} else {
		targetStock.AddStock(riceOrder.Quantity)
	}
	sourceStock.SubtractStock(riceOrder.Quantity)

	sourceStockDoc := NewRiceStockpileDoc(*sourceStock)
	sourceStockDocJSON, err := json.Marshal(sourceStockDoc)
	if err != nil {
		return err
	}

	targetStockDoc := NewRiceStockpileDoc(*targetStock)
	targetStockDocJSON, err := json.Marshal(targetStockDoc)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(sourceStockDoc.ID, sourceStockDocJSON)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(targetStockDoc.ID, targetStockDocJSON)
	if err != nil {
		return err
	}

	err = moveRiceSack(ctx, riceOrder.ID, sourceStockDoc.ID, targetStockDoc.ID, riceOrder.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func createRiceSack(ctx contractapi.TransactionContextInterface, manufacturer *model.Manufacturer, riceOrderId, riceStockpileId string, qty int32) error {
	orgCode := manufacturer.Code
	productCounterDoc, err := getManufacturerProductCounter(ctx, manufacturer.ID)
	if err != nil {
		return err
	}

	for i := 0; i < int(qty); i++ {
		newUuid, err := newDeterministicUuid(fmt.Sprintf("createricesack:%d:%s", i, ctx.GetStub().GetTxID()))
		if err != nil {
			return err
		}

		productCounterDoc.Count += 1
		riceSackDoc := NewRiceSackDoc(model.RiceSack{
			ID:              newUuid,
			RiceOrderID:     riceOrderId,
			RiceStockpileID: riceStockpileId,
			Code:            fmt.Sprintf("(01)%s%06d", orgCode, productCounterDoc.Count),
		})

		riceSackDocJSON, err := json.Marshal(riceSackDoc)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(riceSackDoc.ID, riceSackDocJSON)
		if err != nil {
			return err
		}
	}

	err = putManufacturerProductCounter(ctx, productCounterDoc)
	if err != nil {
		return err
	}

	return nil
}

func moveRiceSack(ctx contractapi.TransactionContextInterface, riceOrderId, sourceStockId, targetStockId string, qty int32) error {
	query := fmt.Sprintf(`{"selector":{"doc_type":"ricesack","rice_stockpile_id":"%s"},"limit":%d}`, sourceStockId, qty)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return err
	}
	defer resultIterator.Close()

	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return err
		}

		var riceSack model.RiceSack
		err = json.Unmarshal(result.Value, &riceSack)
		if err != nil {
			return err
		}

		riceSack.RiceOrderID = riceOrderId
		riceSack.RiceStockpileID = targetStockId
		riceSackDoc := NewRiceSackDoc(riceSack)
		riceSackDocJSON, err := json.Marshal(riceSackDoc)
		if err != nil {
			return nil
		}

		err = ctx.GetStub().PutState(riceSackDoc.ID, riceSackDocJSON)
		if err != nil {
			return err
		}
	}

	return nil
}
