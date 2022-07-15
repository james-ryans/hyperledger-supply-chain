package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type riceSackRepository struct {
	Fabric *client.Gateway
}

func NewRiceSackRepository(fabric *client.Gateway) model.RiceSackRepository {
	return &riceSackRepository{
		Fabric: fabric,
	}
}

func (r *riceSackRepository) FindAll(channelID, stockpileID string) ([]*model.RiceSack, error) {
	network := r.Fabric.GetNetwork(channelID)
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

func (r *riceSackRepository) FindByID(channelID, ID string) (*model.RiceSack, error) {
	network := r.Fabric.GetNetwork(channelID)
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
