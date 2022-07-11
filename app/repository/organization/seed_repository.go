package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/model"
)

type seedRepository struct {
	Fabric *client.Gateway
}

func NewSeedRepository(fabric *client.Gateway) model.SeedRepository {
	return &seedRepository{
		Fabric: fabric,
	}
}

func (r *seedRepository) FindAll(channelID, orgID string) (*[]model.Seed, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	queryString, err := json.Marshal(gin.H{
		"selector": gin.H{
			"doc_type":    "seed",
			"supplier_id": orgID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse querystring: %w", err)
	}

	seedsJSON, err := contract.EvaluateTransaction("SeedContract:QuerySeeds", string(queryString))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var seeds *[]model.Seed
	err = json.Unmarshal(seedsJSON, &seeds)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return seeds, nil
}

func (r *seedRepository) FindByID(channelID, ID string) (*model.Seed, error) {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	seedJSON, err := contract.EvaluateTransaction("SeedContract:ReadSeed", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	seed := &model.Seed{}
	err = json.Unmarshal(seedJSON, seed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return seed, nil
}

func (r *seedRepository) Create(channelID string, seed *model.Seed) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("SeedContract:CreateSeed", seed.ID, seed.SupplierID, seed.VarietyName, strconv.FormatFloat(float64(seed.PlantAge), 'f', -1, 32), seed.PlantShape, strconv.FormatFloat(float64(seed.PlantHeight), 'f', -1, 32), seed.LeafShape)
	if err != nil {
		return fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return nil
}

func (r *seedRepository) Update(channelID string, seed *model.Seed) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("SeedContract:UpdateSeed", seed.ID, seed.SupplierID, seed.VarietyName, strconv.FormatFloat(float64(seed.PlantAge), 'f', -1, 32), seed.PlantShape, strconv.FormatFloat(float64(seed.PlantHeight), 'f', -1, 32), seed.LeafShape)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

func (r *seedRepository) Delete(channelID, ID string) error {
	network := r.Fabric.GetNetwork(channelID)
	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))

	_, err := contract.SubmitTransaction("SeedContract:DeleteSeed", ID)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
