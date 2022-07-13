package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type distributorRepository struct {
	Fabric *client.Gateway
}

func NewDistributorRepository(fabric *client.Gateway) model.DistributorRepository {
	return &distributorRepository{
		Fabric: fabric,
	}
}

func (r *distributorRepository) FindAll(channelID string) ([]*model.Distributor, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	distributorsJSON, err := contract.EvaluateTransaction("DistributorContract:FindAll")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var distributors []*model.Distributor
	err = json.Unmarshal(distributorsJSON, &distributors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return distributors, nil
}

func (r *distributorRepository) FindByID(channelID, ID string) (*model.Distributor, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	distributorJSON, err := contract.EvaluateTransaction("DistributorContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	distributor := model.Distributor{}
	err = json.Unmarshal(distributorJSON, &distributor)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &distributor, nil
}

func (r *distributorRepository) Create(channelID string, distributor *model.Distributor) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction(
		"DistributorContract:Create",
		distributor.ID,
		distributor.Type,
		distributor.Name,
		distributor.Location.Province,
		distributor.Location.City,
		distributor.Location.District,
		distributor.Location.PostalCode,
		distributor.Location.Address,
		distributor.ContactInfo.Phone,
		distributor.ContactInfo.Email,
		strconv.FormatFloat(float64(distributor.Location.Coordinate.Longitude), 'f', -1, 32),
		strconv.FormatFloat(float64(distributor.Location.Coordinate.Latitude), 'f', -1, 32),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
