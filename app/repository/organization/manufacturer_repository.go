package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type manufacturerRepository struct {
	Fabric *client.Gateway
}

func NewManufacturerRepository(fabric *client.Gateway) model.ManufacturerRepository {
	return &manufacturerRepository{
		Fabric: fabric,
	}
}

func (r *manufacturerRepository) FindAll(channelID string) ([]*model.Manufacturer, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	manufacturersJSON, err := contract.EvaluateTransaction("ManufacturerContract:FindAll")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var manufacturers []*model.Manufacturer
	err = json.Unmarshal(manufacturersJSON, &manufacturers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return manufacturers, nil
}

func (r *manufacturerRepository) FindByID(channelID, ID string) (*model.Manufacturer, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	manufacturerJSON, err := contract.EvaluateTransaction("ManufacturerContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	manufacturer := model.Manufacturer{}
	err = json.Unmarshal(manufacturerJSON, &manufacturer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &manufacturer, nil
}

func (r *manufacturerRepository) Create(channelID string, manufacturer *model.Manufacturer) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction(
		"ManufacturerContract:Create",
		manufacturer.ID,
		manufacturer.Type,
		manufacturer.Name,
		manufacturer.Code,
		manufacturer.Location.Province,
		manufacturer.Location.City,
		manufacturer.Location.District,
		manufacturer.Location.PostalCode,
		manufacturer.Location.Address,
		manufacturer.ContactInfo.Phone,
		manufacturer.ContactInfo.Email,
		strconv.FormatFloat(float64(manufacturer.Location.Coordinate.Longitude), 'f', -1, 32),
		strconv.FormatFloat(float64(manufacturer.Location.Coordinate.Latitude), 'f', -1, 32),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
