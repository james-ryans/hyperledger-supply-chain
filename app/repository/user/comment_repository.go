package userrepository

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type commentRepository struct {
	Fabric *gateway.Gateway
}

func NewCommentRepository(fabric *gateway.Gateway) usermodel.CommentRepository {
	return &commentRepository{
		Fabric: fabric,
	}
}

func (r *commentRepository) FindAllByRiceID(riceID string) ([]*usermodel.Comment, error) {
	network, err := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	if err != nil {
		return nil, fmt.Errorf("failed to get network %s: %w", os.Getenv("FABRIC_GLOBALCHANNEL_NAME"), err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_GLOBALCHAINCODE_NAME"))

	commentsJSON, err := contract.EvaluateTransaction("CommentContract:FindAllByRiceId", riceID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction %w", err)
	}

	var comments []*usermodel.Comment
	err = json.Unmarshal(commentsJSON, &comments)
	if err != nil {
		return nil, fmt.Errorf("failed to parsed result: %w", err)
	}

	return comments, nil
}

func (r *commentRepository) Create(comment *usermodel.Comment) error {
	network, err := r.Fabric.GetNetwork(os.Getenv("FABRIC_GLOBALCHANNEL_NAME"))
	if err != nil {
		return fmt.Errorf("failed to get network %s: %w", os.Getenv("FABRIC_GLOBALCHANNEL_NAME"), err)
	}
	contract := network.GetContract(os.Getenv("FABRIC_GLOBALCHAINCODE_NAME"))

	_, err = contract.SubmitTransaction(
		"CommentContract:Create",
		comment.ID,
		comment.RiceID,
		comment.UserName,
		comment.Text,
		comment.CommentAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}
