package contract

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type CommentContract struct {
	contractapi.Contract
}

type CommentDoc struct {
	DocType string `json:"doc_type"`
	usermodel.Comment
}

func NewCommentDoc(comment usermodel.Comment) CommentDoc {
	return CommentDoc{
		DocType: "comment",
		Comment: comment,
	}
}

func (c *CommentContract) FindAllByRiceId(ctx contractapi.TransactionContextInterface, riceId string) ([]*usermodel.Comment, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"comment","rice_id":"%s"}}`, riceId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	comments := make([]*usermodel.Comment, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var comment usermodel.Comment
		err = json.Unmarshal(result.Value, &comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (c *CommentContract) Create(ctx contractapi.TransactionContextInterface, id, riceId, userName, text string, commentAt time.Time) error {
	commentDoc := NewCommentDoc(
		usermodel.Comment{
			ID:        id,
			RiceID:    riceId,
			UserName:  userName,
			Text:      text,
			CommentAt: commentAt,
		},
	)
	commentDocJSON, err := json.Marshal(commentDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, commentDocJSON)
}
