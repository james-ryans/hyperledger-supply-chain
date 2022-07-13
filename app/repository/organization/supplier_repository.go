package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type supplierRepository struct {
	Fabric *client.Gateway
}

func NewSupplierRepository(fabric *client.Gateway) model.SupplierRepository {
	return &supplierRepository{
		Fabric: fabric,
	}
}

func (r *supplierRepository) FindAll(channelID string) ([]*model.Supplier, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	suppliersJSON, err := contract.EvaluateTransaction("SupplierContract:FindAll")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var suppliers []*model.Supplier
	err = json.Unmarshal(suppliersJSON, &suppliers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return suppliers, nil
}

func (r *supplierRepository) FindByID(channelID, ID string) (*model.Supplier, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	supplierJSON, err := contract.EvaluateTransaction("SupplierContract:FindByID", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	supplier := model.Supplier{}
	err = json.Unmarshal(supplierJSON, &supplier)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &supplier, nil
}

func (r *supplierRepository) Create(channelID string, supplier *model.Supplier) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction(
		"SupplierContract:Create",
		supplier.ID,
		supplier.Type,
		supplier.Name,
		supplier.Location.Province,
		supplier.Location.City,
		supplier.Location.District,
		supplier.Location.PostalCode,
		supplier.Location.Address,
		supplier.ContactInfo.Phone,
		supplier.ContactInfo.Email,
		strconv.FormatFloat(float64(supplier.Location.Coordinate.Longitude), 'f', -1, 32),
		strconv.FormatFloat(float64(supplier.Location.Coordinate.Latitude), 'f', -1, 32),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
