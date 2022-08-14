package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/meneketehe/hehe/app/model"
)

type riceRepository struct {
	Fabric *gateway.Gateway
}

func NewRiceRepository(fabric *gateway.Gateway) model.RiceRepository {
	return &riceRepository{
		Fabric: fabric,
	}
}

func (r *riceRepository) FindAll(channelID, orgID string) (*[]model.Rice, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	queryString, err := json.Marshal(gin.H{
		"selector": gin.H{
			"doc_type":        "rice",
			"manufacturer_id": orgID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse querystring: %w", err)
	}

	ricesJSON, err := contract.EvaluateTransaction("RiceContract:QueryRices", string(queryString))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var rices *[]model.Rice
	err = json.Unmarshal(ricesJSON, &rices)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return rices, nil
}

func (r *riceRepository) FindByID(channelID, ID string) (*model.Rice, error) {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceJSON, err := contract.EvaluateTransaction("RiceContract:ReadRice", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	rice := &model.Rice{}
	err = json.Unmarshal(riceJSON, rice)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return rice, nil
}

func (r *riceRepository) Create(channelID string, rice *model.Rice) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction("RiceContract:CreateRice", rice.ID, rice.ManufacturerID, rice.Code, rice.BrandName, strconv.FormatFloat(float64(rice.Weight), 'f', -1, 32), rice.Texture, strconv.FormatFloat(float64(rice.AmyloseRate), 'f', -1, 32))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceRepository) Update(channelID string, rice *model.Rice) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction("RiceContract:UpdateRice", rice.ID, rice.ManufacturerID, rice.Code, rice.BrandName, strconv.FormatFloat(float64(rice.Weight), 'f', -1, 32), rice.Texture, strconv.FormatFloat(float64(rice.AmyloseRate), 'f', -1, 32))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceRepository) Delete(channelID, ID string) error {
	network, err := r.Fabric.GetNetwork(channelID)
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", channelID, err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err = contract.SubmitTransaction("RiceContract:DeleteRice", ID)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
