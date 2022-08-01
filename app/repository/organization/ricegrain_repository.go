package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type riceGrainRepository struct {
	Fabric *client.Gateway
}

func NewRiceGrainRepository(fabric *client.Gateway) model.RiceGrainRepository {
	return &riceGrainRepository{
		Fabric: fabric,
	}
}

func (r *riceGrainRepository) FindAll(channelID, orgID string) (*[]model.RiceGrain, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	queryString, err := json.Marshal(gin.H{
		"selector": gin.H{
			"doc_type":    "ricegrain",
			"producer_id": orgID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse querystring: %w", err)
	}

	riceGrainsJSON, err := contract.EvaluateTransaction("RiceGrainContract:QueryRiceGrains", string(queryString))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var riceGrains *[]model.RiceGrain
	err = json.Unmarshal(riceGrainsJSON, &riceGrains)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceGrains, nil
}

func (r *riceGrainRepository) FindByID(channelID, ID string) (*model.RiceGrain, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	riceGrainJSON, err := contract.EvaluateTransaction("RiceGrainContract:ReadRiceGrain", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	riceGrain := &model.RiceGrain{}
	err = json.Unmarshal(riceGrainJSON, riceGrain)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return riceGrain, nil
}

func (r *riceGrainRepository) Create(channelID string, riceGrain *model.RiceGrain) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("RiceGrainContract:CreateRiceGrain", riceGrain.ID, riceGrain.ProducerID, riceGrain.VarietyName, riceGrain.GrainShape, riceGrain.GrainColor)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceGrainRepository) Update(channelID string, riceGrain *model.RiceGrain) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("RiceGrainContract:UpdateRiceGrain", riceGrain.ID, riceGrain.ProducerID, riceGrain.VarietyName, riceGrain.GrainShape, riceGrain.GrainColor)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *riceGrainRepository) Delete(channelID, ID string) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("RiceGrainContract:DeleteRiceGrain", ID)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
