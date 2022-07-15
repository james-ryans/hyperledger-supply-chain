package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type riceStockpileRepository struct {
	Fabric *client.Gateway
}

func NewRiceStockpileRepository(fabric *client.Gateway) model.RiceStockpileRepository {
	return &riceStockpileRepository{
		Fabric: fabric,
	}
}

func (r *riceStockpileRepository) FindAll(channelID, vendorID string) ([]*model.RiceStockpile, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceStockpilesJSON, err := contract.EvaluateTransaction("RiceStockpileContract:FindAll", vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	var riceStockpiles []*model.RiceStockpile
	err = json.Unmarshal(riceStockpilesJSON, &riceStockpiles)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceStockpiles, nil
}

func (r *riceStockpileRepository) FindByID(channelID, ID string) (*model.RiceStockpile, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceStockpileJSON, err := contract.EvaluateTransaction("RiceStockpileContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	var riceStockpile *model.RiceStockpile
	err = json.Unmarshal(riceStockpileJSON, &riceStockpile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceStockpile, nil
}

func (r *riceStockpileRepository) FindByVendorIDAndRiceID(channelID, vendorID, riceID string) (*model.RiceStockpile, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceStockpileJSON, err := contract.EvaluateTransaction("RiceStockpileContract:FindByVendorIDAndRiceID", vendorID, riceID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	var riceStockpile *model.RiceStockpile
	err = json.Unmarshal(riceStockpileJSON, &riceStockpile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceStockpile, nil
}
