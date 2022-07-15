package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type riceOrderRepository struct {
	Fabric *client.Gateway
}

func NewRiceOrderRepository(fabric *client.Gateway) model.RiceOrderRepository {
	return &riceOrderRepository{
		Fabric: fabric,
	}
}

func (r *riceOrderRepository) FindAllOutgoing(channelID, ordererID string) ([]*model.RiceOrder, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceOrdersJSON, err := contract.EvaluateTransaction("RiceOrderContract:FindAllOutgoing", ordererID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var riceOrders []*model.RiceOrder
	err = json.Unmarshal(riceOrdersJSON, &riceOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceOrders, nil
}

func (r *riceOrderRepository) FindAllIncoming(channelID, sellerID string) ([]*model.RiceOrder, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceOrdersJSON, err := contract.EvaluateTransaction("RiceOrderContract:FindAllIncoming", sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var riceOrders []*model.RiceOrder
	err = json.Unmarshal(riceOrdersJSON, &riceOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceOrders, nil
}

func (r *riceOrderRepository) FindAllAcceptedIncoming(channelID, sellerID string) ([]*model.RiceOrder, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceOrdersJSON, err := contract.EvaluateTransaction("RiceOrderContract:FindAllAcceptedIncoming", sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var riceOrders []*model.RiceOrder
	err = json.Unmarshal(riceOrdersJSON, &riceOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceOrders, nil
}

func (r *riceOrderRepository) FindByID(channelID, ID string) (*model.RiceOrder, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceOrderJSON, err := contract.EvaluateTransaction("RiceOrderContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	riceOrder := model.RiceOrder{}
	err = json.Unmarshal(riceOrderJSON, &riceOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &riceOrder, nil
}

func (r *riceOrderRepository) Create(channelID string, riceOrder *model.RiceOrder) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction(
		"RiceOrderContract:Create",
		riceOrder.ID,
		riceOrder.OrdererID,
		riceOrder.SellerID,
		riceOrder.RiceID,
		strconv.FormatInt(int64(riceOrder.Quantity), 10),
		riceOrder.OrderedAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceOrderRepository) Accept(channelID string, ID string, acceptedAt time.Time) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("RiceOrderContract:Accept", ID, acceptedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceOrderRepository) AcceptDistribution(channelID string, ID string, acceptedAt time.Time) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("RiceOrderContract:AcceptDistribution", ID, acceptedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceOrderRepository) Reject(channelID string, ID string, rejectedAt time.Time, reason string) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("RiceOrderContract:Reject", ID, rejectedAt.Format(time.RFC3339), reason)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceOrderRepository) Ship(channelID string, ID string, shippedAt time.Time, grade string, millingDate time.Time, storageTemperature float32, storageHumidity float32) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction(
		"RiceOrderContract:Ship",
		ID,
		shippedAt.Format(time.RFC3339),
		grade,
		millingDate.Format(time.RFC3339),
		strconv.FormatFloat(float64(storageTemperature), 'f', -1, 32),
		strconv.FormatFloat(float64(storageHumidity), 'f', -1, 32),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceOrderRepository) ShipDistribution(channelID string, ID string, shippedAt time.Time, grade string, millingDate time.Time, storageTemperature float32, storageHumidity float32) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction(
		"RiceOrderContract:ShipDistribution",
		ID,
		shippedAt.Format(time.RFC3339),
		grade,
		millingDate.Format(time.RFC3339),
		strconv.FormatFloat(float64(storageTemperature), 'f', -1, 32),
		strconv.FormatFloat(float64(storageHumidity), 'f', -1, 32),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceOrderRepository) Receive(channelID string, ID string, receivedAt time.Time) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("RiceOrderContract:Receive", ID, receivedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceOrderRepository) ReceiveDistribution(channelID string, ID string, receivedAt time.Time) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("RiceOrderContract:ReceiveDistribution", ID, receivedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
