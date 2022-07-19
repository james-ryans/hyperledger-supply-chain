package userrepository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type riceSackRepository struct {
	Fabric *client.Gateway
}

func NewRiceSackRepository(fabric *client.Gateway) usermodel.RiceSackRepository {
	return &riceSackRepository{
		Fabric: fabric,
	}
}

func (r *riceSackRepository) FindByCode(userID, code string) (*usermodel.RiceSack, error) {
	network := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	contract := network.GetContract(os.Getenv("FABRIC_GLOBALCHAINCODE_NAME"))

	riceSackJSON, err := contract.SubmitTransaction("RiceSackContract:FindByCode", userID, code)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %w", err)
	}

	riceSack, err := usermodel.UnmarshalRiceSack(riceSackJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parsed result: %w", err)
	}

	return riceSack, nil
}

func (r *riceSackRepository) Create(sack *usermodel.RiceSack) error {
	network := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	contract := network.GetContract(os.Getenv("FABRIC_GLOBALCHAINCODE_NAME"))

	sackJSON, err := json.Marshal(sack)
	if err != nil {
		return err
	}

	_, err = contract.SubmitTransaction("RiceSackContract:Create", string(sackJSON))
	if err != nil {
		return fmt.Errorf("failed to submit transaction %w", err)
	}

	return nil
}
