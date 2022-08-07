package userrepository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type scanHistoryRepository struct {
	Fabric *gateway.Gateway
}

func NewScanHistoryRepository(fabric *gateway.Gateway) usermodel.ScanHistoryRepository {
	return &scanHistoryRepository{
		Fabric: fabric,
	}
}

func (r *scanHistoryRepository) FindAll(userID string) ([]*usermodel.ScanHistory, error) {
	network, err := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", os.Getenv("FABRIC_GLOBALCHANNEL_NAME"), err)
	}
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
