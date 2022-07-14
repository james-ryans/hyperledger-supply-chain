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

type seedOrderRepository struct {
	Fabric *client.Gateway
}

func NewSeedOrderRepository(fabric *client.Gateway) model.SeedOrderRepository {
	return &seedOrderRepository{
		Fabric: fabric,
	}
}

func (r *seedOrderRepository) FindAllOutgoing(channelID, ordererID string) ([]*model.SeedOrder, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	seedOrdersJSON, err := contract.EvaluateTransaction("SeedOrderContract:FindAllOutgoing", ordererID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var seedOrders []*model.SeedOrder
	err = json.Unmarshal(seedOrdersJSON, &seedOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return seedOrders, nil
}

func (r *seedOrderRepository) FindAllIncoming(channelID, sellerID string) ([]*model.SeedOrder, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	seedOrdersJSON, err := contract.EvaluateTransaction("SeedOrderContract:FindAllIncoming", sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var seedOrders []*model.SeedOrder
	err = json.Unmarshal(seedOrdersJSON, &seedOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return seedOrders, nil
}

func (r *seedOrderRepository) FindByID(channelID, ID string) (*model.SeedOrder, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	seedOrderJSON, err := contract.EvaluateTransaction("SeedOrderContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	seedOrder := model.SeedOrder{}
	err = json.Unmarshal(seedOrderJSON, &seedOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &seedOrder, nil
}

func (r *seedOrderRepository) Create(channelID string, seedOrder *model.SeedOrder) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction(
		"SeedOrderContract:Create",
		seedOrder.ID,
		seedOrder.OrdererID,
		seedOrder.SellerID,
		seedOrder.SeedID,
		seedOrder.RiceGrainOrderID,
		strconv.FormatFloat(float64(seedOrder.Weight), 'f', -1, 32),
		seedOrder.OrderedAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *seedOrderRepository) Accept(channelID string, ID string, acceptedAt time.Time) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("SeedOrderContract:Accept", ID, acceptedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *seedOrderRepository) Reject(channelID string, ID string, rejectedAt time.Time, reason string) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("SeedOrderContract:Reject", ID, rejectedAt.Format(time.RFC3339), reason)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *seedOrderRepository) Ship(channelID string, ID string, shippedAt time.Time) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("SeedOrderContract:Ship", ID, shippedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *seedOrderRepository) Receive(channelID string, ID string, receivedAt time.Time) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("SeedOrderContract:Receive", ID, receivedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
