package userrepository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type scanHistoryRepository struct {
	Fabric *client.Gateway
}

func NewScanHistoryRepository(fabric *client.Gateway) usermodel.ScanHistoryRepository {
	return &scanHistoryRepository{
		Fabric: fabric,
	}
}

func (r *scanHistoryRepository) FindAll(userID string) ([]*usermodel.ScanHistory, error) {
	network := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	contract := network.GetContract(os.Getenv("FABRIC_GLOBALCHAINCODE_NAME"))

	scanHistoriesJSON, err := contract.EvaluateTransaction("ScanHistoryContract:FindAll", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var scanHistories []*usermodel.ScanHistory
	err = json.Unmarshal(scanHistoriesJSON, &scanHistories)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return scanHistories, nil
}
