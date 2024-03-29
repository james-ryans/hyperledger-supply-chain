package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/meneketehe/hehe/app/model"
)

type riceGrainOrderRepository struct {
	Fabric *gateway.Gateway
}

func NewRiceGrainOrderRepository(fabric *gateway.Gateway) model.RiceGrainOrderRepository {
	return &riceGrainOrderRepository{
		Fabric: fabric,
	}
}

func (r *riceGrainOrderRepository) FindAllOutgoing(channelID, ordererID string) ([]*model.RiceGrainOrder, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceGrainOrdersJSON, err := contract.EvaluateTransaction("RiceGrainOrderContract:FindAllOutgoing", ordererID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var riceGrainOrders []*model.RiceGrainOrder
	err = json.Unmarshal(riceGrainOrdersJSON, &riceGrainOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceGrainOrders, nil
}

func (r *riceGrainOrderRepository) FindAllIncoming(channelID, sellerID string) ([]*model.RiceGrainOrder, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceGrainOrdersJSON, err := contract.EvaluateTransaction("RiceGrainOrderContract:FindAllIncoming", sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var riceGrainOrders []*model.RiceGrainOrder
	err = json.Unmarshal(riceGrainOrdersJSON, &riceGrainOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceGrainOrders, nil
}

func (r *riceGrainOrderRepository) FindAllAcceptedIncoming(channelID, sellerID string) ([]*model.RiceGrainOrder, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceGrainOrdersJSON, err := contract.EvaluateTransaction("RiceGrainOrderContract:FindAllAcceptedIncoming", sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var riceGrainOrders []*model.RiceGrainOrder
	err = json.Unmarshal(riceGrainOrdersJSON, &riceGrainOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceGrainOrders, nil
}

func (r *riceGrainOrderRepository) FindByID(channelID, ID string) (*model.RiceGrainOrder, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceGrainOrderJSON, err := contract.EvaluateTransaction("RiceGrainOrderContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	riceGrainOrder := model.RiceGrainOrder{}
	err = json.Unmarshal(riceGrainOrderJSON, &riceGrainOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &riceGrainOrder, nil
}

func (r *riceGrainOrderRepository) FindByRiceOrderID(channelID, riceOrderID string) (*model.RiceGrainOrder, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceGrainOrderJSON, err := contract.EvaluateTransaction("RiceGrainOrderContract:FindByRiceOrderID", riceOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	riceGrainOrder := model.RiceGrainOrder{}
	err = json.Unmarshal(riceGrainOrderJSON, &riceGrainOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &riceGrainOrder, err
}

func (r *riceGrainOrderRepository) Create(channelID string, riceGrainOrder *model.RiceGrainOrder) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction(
		"RiceGrainOrderContract:Create",
		riceGrainOrder.ID,
		riceGrainOrder.OrdererID,
		riceGrainOrder.SellerID,
		riceGrainOrder.RiceGrainID,
		riceGrainOrder.RiceOrderID,
		strconv.FormatFloat(float64(riceGrainOrder.Weight), 'f', -1, 32),
		riceGrainOrder.OrderedAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceGrainOrderRepository) Accept(channelID, ID string, acceptedAt time.Time) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction("RiceGrainOrderContract:Accept", ID, acceptedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceGrainOrderRepository) Reject(channelID, ID string, rejectedAt time.Time, reason string) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction("RiceGrainOrderContract:Reject", ID, rejectedAt.Format(time.RFC3339), reason)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceGrainOrderRepository) Ship(channelID string, ID string, shippedAt time.Time, plowMethod, sowMethod, irrigation, fertilization string, plantDate, harvestDate time.Time, storageTemperature, storageHumidity float32) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction(
		"RiceGrainOrderContract:Ship",
		ID,
		shippedAt.Format(time.RFC3339),
		plowMethod,
		sowMethod,
		irrigation,
		fertilization,
		plantDate.Format(time.RFC3339),
		harvestDate.Format(time.RFC3339),
		strconv.FormatFloat(float64(storageTemperature), 'f', -1, 32),
		strconv.FormatFloat(float64(storageHumidity), 'f', -1, 32),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceGrainOrderRepository) Receive(channelID string, ID string, receivedAt time.Time) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction("RiceGrainOrderContract:Receive", ID, receivedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
