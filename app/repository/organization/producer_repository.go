package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/meneketehe/hehe/app/model"
)

type producerRepository struct {
	Fabric *gateway.Gateway
}

func NewProducerRepository(fabric *gateway.Gateway) model.ProducerRepository {
	return &producerRepository{
		Fabric: fabric,
	}
}

func (r *producerRepository) FindAll(channelID string) ([]*model.Producer, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	producersJSON, err := contract.EvaluateTransaction("ProducerContract:FindAll")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var producers []*model.Producer
	err = json.Unmarshal(producersJSON, &producers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return producers, nil
}

func (r *producerRepository) FindByID(channelID, ID string) (*model.Producer, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	producerJSON, err := contract.EvaluateTransaction("ProducerContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	producer := model.Producer{}
	err = json.Unmarshal(producerJSON, &producer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &producer, nil
}

func (r *producerRepository) Create(channelID string, producer *model.Producer) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction(
		"ProducerContract:Create",
		producer.ID,
		producer.Type,
		producer.Name,
		producer.Location.Province,
		producer.Location.City,
		producer.Location.District,
		producer.Location.PostalCode,
		producer.Location.Address,
		producer.ContactInfo.Phone,
		producer.ContactInfo.Email,
		strconv.FormatFloat(float64(producer.Location.Coordinate.Longitude), 'f', -1, 32),
		strconv.FormatFloat(float64(producer.Location.Coordinate.Latitude), 'f', -1, 32),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
