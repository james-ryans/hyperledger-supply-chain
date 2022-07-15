package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type RiceStockpileContract struct {
	contractapi.Contract
}

type RiceStockpileDoc struct {
	DocType string `json:"doc_type"`
	model.RiceStockpile
}

func NewRiceStockpileDoc(riceStockpile model.RiceStockpile) RiceStockpileDoc {
	return RiceStockpileDoc{
		DocType:       "ricestockpile",
		RiceStockpile: riceStockpile,
	}
}

func (c *RiceStockpileContract) FindAll(ctx contractapi.TransactionContextInterface, vendorId string) ([]*model.RiceStockpile, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"ricestockpile","vendor_id":"%s"}}`, vendorId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	stockpiles := make([]*model.RiceStockpile, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var stockpile model.RiceStockpile
		err = json.Unmarshal(result.Value, &stockpile)
		if err != nil {
			return nil, err
		}
		stockpiles = append(stockpiles, &stockpile)
	}

	return stockpiles, nil
}

func (c *RiceStockpileContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.RiceStockpile, error) {
	stockpile, err := c.getRiceStockpile(ctx, id)
	if err != nil {
		return nil, err
	}
	if stockpile == nil {
		return nil, fmt.Errorf("the rice stockpile %s does not exist", id)
	}

	return stockpile, nil
}

func (c *RiceStockpileContract) getRiceStockpile(ctx contractapi.TransactionContextInterface, id string) (*model.RiceStockpile, error) {
	stockpileJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if stockpileJSON == nil {
		return nil, nil
	}

	var stockpile model.RiceStockpile
	err = json.Unmarshal(stockpileJSON, &stockpile)
	if err != nil {
		return nil, err
	}

	return &stockpile, nil
}
