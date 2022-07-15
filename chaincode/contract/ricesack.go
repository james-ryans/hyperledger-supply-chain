package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/app/model"
)

type RiceSackContract struct {
	contractapi.Contract
}

type RiceSackDoc struct {
	DocType string `json:"doc_type"`
	model.RiceSack
}

func NewRiceSackDoc(riceSack model.RiceSack) RiceSackDoc {
	return RiceSackDoc{
		DocType:  "ricesack",
		RiceSack: riceSack,
	}
}

func (c *RiceSackContract) FindAll(ctx contractapi.TransactionContextInterface, stockpileId string) ([]*model.RiceSack, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"ricesack","rice_stockpile_id":"%s"}}`, stockpileId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	sacks := make([]*model.RiceSack, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var sack model.RiceSack
		err = json.Unmarshal(result.Value, &sack)
		if err != nil {
			return nil, err
		}
		sacks = append(sacks, &sack)
	}

	return sacks, nil
}

func (c *RiceSackContract) FindByID(ctx contractapi.TransactionContextInterface, id string) (*model.RiceSack, error) {
	sack, err := getRiceSack(ctx, id)
	if err != nil {
		return nil, err
	}
	if sack == nil {
		return nil, fmt.Errorf("the rice sack %s does not exist", id)
	}

	return sack, nil
}

func getRiceSack(ctx contractapi.TransactionContextInterface, id string) (*model.RiceSack, error) {
	riceSackJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if riceSackJSON == nil {
		return nil, nil
	}

	var riceSack model.RiceSack
	err = json.Unmarshal(riceSackJSON, &riceSack)
	if err != nil {
		return nil, err
	}

	return &riceSack, nil
}
