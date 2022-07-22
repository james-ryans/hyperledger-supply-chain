package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type UserContract struct {
	contractapi.Contract
}

type UserDoc struct {
	DocType string `json:"doc_type"`
	usermodel.User
}

func NewUserDoc(user usermodel.User) UserDoc {
	return UserDoc{
		DocType: "user",
		User:    user,
	}
}

func (c *UserContract) FindById(ctx contractapi.TransactionContextInterface, id string) (*usermodel.User, error) {
	user, err := getUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("the user %s does not exist", id)
	}

	return user, nil
}

func (c *UserContract) Create(ctx contractapi.TransactionContextInterface, id, name, email string) error {
	userDoc := NewUserDoc(
		usermodel.User{
			ID:            id,
			Name:          name,
			Email:         email,
			ScanHistories: []usermodel.ScanHistory{},
		},
	)
	userDocJSON, err := json.Marshal(userDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, userDocJSON)
}

func (c *UserContract) Update(ctx contractapi.TransactionContextInterface, id, name, email string) error {
	user, err := getUser(ctx, id)
	if err != nil {
		return fmt.Errorf("user %s does not exist", id)
	}

	userDoc := NewUserDoc(*user)
	userDoc.Name = name
	userDoc.Email = email
	userDocJSON, err := json.Marshal(userDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(userDoc.ID, userDocJSON)
}

func getUser(ctx contractapi.TransactionContextInterface, id string) (*usermodel.User, error) {
	userJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if userJSON == nil {
		return nil, nil
	}

	var user usermodel.User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
