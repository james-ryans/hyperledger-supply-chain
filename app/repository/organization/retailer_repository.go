package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type retailerRepository struct {
	Fabric *client.Gateway
}

func NewRetailerRepository(fabric *client.Gateway) model.RetailerRepository {
	return &retailerRepository{
		Fabric: fabric,
	}
}

func (r *retailerRepository) FindAll(channelID string) ([]*model.Retailer, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	retailersJSON, err := contract.EvaluateTransaction("RetailerContract:FindAll")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var retailers []*model.Retailer
	err = json.Unmarshal(retailersJSON, &retailers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return retailers, nil
}

func (r *retailerRepository) FindByID(channelID, ID string) (*model.Retailer, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	retailerJSON, err := contract.EvaluateTransaction("RetailerContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	retailer := model.Retailer{}
	err = json.Unmarshal(retailerJSON, &retailer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &retailer, nil
}

func (r *retailerRepository) Create(channelID string, retailer *model.Retailer) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction(
		"RetailerContract:Create",
		retailer.ID,
		retailer.Type,
		retailer.Name,
		retailer.Location.Province,
		retailer.Location.City,
		retailer.Location.District,
		retailer.Location.PostalCode,
		retailer.Location.Address,
		retailer.ContactInfo.Phone,
		retailer.ContactInfo.Email,
		strconv.FormatFloat(float64(retailer.Location.Coordinate.Longitude), 'f', -1, 32),
		strconv.FormatFloat(float64(retailer.Location.Coordinate.Latitude), 'f', -1, 32),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
