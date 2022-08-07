package userrepository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type userRepository struct {
	Fabric *gateway.Gateway
}

func NewUserRepository(fabric *gateway.Gateway) usermodel.UserRepository {
	return &userRepository{
		Fabric: fabric,
	}
}

func (r *userRepository) FindByID(ID string) (*usermodel.User, error) {
	network, err := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", os.Getenv("FABRIC_GLOBALCHANNEL_NAME"), err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_GLOBALCHAINCODE_NAME"))

	userJSON, err := contract.EvaluateTransaction("UserContract:FindById", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var user usermodel.User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(user *usermodel.User) error {
	network, err := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", os.Getenv("FABRIC_GLOBALCHANNEL_NAME"), err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_GLOBALCHAINCODE_NAME"))

	_, err = contract.SubmitTransaction(
		"UserContract:Create",
		user.ID,
		user.Name,
		user.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to submit transanction: %w", err)
	}

	return nil
}

func (r *userRepository) Update(user *usermodel.User) error {
	network, err := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", os.Getenv("FABRIC_GLOBALCHANNEL_NAME"), err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_GLOBALCHAINCODE_NAME"))

	_, err = contract.SubmitTransaction(
		"UserContract:Update",
		user.ID,
		user.Name,
		user.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to submit transanction: %w", err)
	}

	return nil
}
