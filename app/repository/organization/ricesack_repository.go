package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/meneketehe/hehe/app/model"
)

type riceSackRepository struct {
	Fabric *gateway.Gateway
}

func NewRiceSackRepository(fabric *gateway.Gateway) model.RiceSackRepository {
	return &riceSackRepository{
		Fabric: fabric,
	}
}

func (r *riceSackRepository) FindAll(channelID, stockpileID string) ([]*model.RiceSack, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceSacksJSON, err := contract.EvaluateTransaction("RiceSackContract:FindAll", stockpileID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	var riceSacks []*model.RiceSack
	err = json.Unmarshal(riceSacksJSON, &riceSacks)
	if err != nil {
		return nil, fmt.Errorf("failed to parsed result: %w", err)
	}

	return riceSacks, nil
}

func (r *riceSackRepository) FindAllByRiceOrderID(channelID, riceOrderID string) ([]*model.RiceSack, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceSacksJSON, err := contract.EvaluateTransaction("RiceSackContract:FindAllByRiceOrderId", riceOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	var riceSacks []*model.RiceSack
	err = json.Unmarshal(riceSacksJSON, &riceSacks)
	if err != nil {
		return nil, fmt.Errorf("failed to parsed result: %w", err)
	}

	return riceSacks, nil
}

func (r *riceSackRepository) FindByID(channelID, ID string) (*model.RiceSack, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceSackJSON, err := contract.EvaluateTransaction("RiceSackContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	var riceSack *model.RiceSack
	err = json.Unmarshal(riceSackJSON, &riceSack)
	if err != nil {
		return nil, fmt.Errorf("failed to parsed result: %w", err)
	}

	return riceSack, nil
}
