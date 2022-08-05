package contract

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type RiceSackContract struct {
	contractapi.Contract
}

type RiceSackDoc struct {
	DocType string `json:"doc_type"`
	usermodel.RiceSack
}

func NewRiceSackDoc(riceSack usermodel.RiceSack) RiceSackDoc {
	return RiceSackDoc{
		DocType:  "ricesack",
		RiceSack: riceSack,
	}
}

func (c *RiceSackContract) FindByCode(ctx contractapi.TransactionContextInterface, userId, code string) (string, error) {
	sack, err := getRiceSack(ctx, code)
	if err != nil {
		return "", err
	}
	if sack == nil {
		return "", fmt.Errorf("the rice sack %s does not exist", code)
	}

	sackJSON, err := json.Marshal(sack)
	if err != nil {
		return "", fmt.Errorf("failed to parse result: %w", err)
	}

	id, err := newDeterministicUuid(fmt.Sprintf("%s", ctx.GetStub().GetTxID()))
	if err != nil {
		return "", err
	}

	scanAt, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", err
	}

	lastTrace := sack.Traces[len(sack.Traces)-1]

	riceJSON, err := json.Marshal(lastTrace.Commodity)
	if err != nil {
		return "", err
	}

	var rice usermodel.Rice
	err = json.Unmarshal(riceJSON, &rice)
	if err != nil {
		return "", err
	}

	if userId != "" {
		scanHistoryDoc := NewScanHistoryDoc(usermodel.ScanHistory{
			ID:           id,
			UserID:       userId,
			RiceSackCode: sack.Code,
			Name:         fmt.Sprintf("%s - %s - %s", rice.BrandName, lastTrace.Organization.Name, sack.Code),
			ScanAt:       scanAt.AsTime(),
		})
		scanHistoryDocJSON, err := json.Marshal(scanHistoryDoc)
		if err != nil {
			return "", err
		}

		err = ctx.GetStub().PutState(scanHistoryDoc.ID, scanHistoryDocJSON)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s", sackJSON), nil
}

func (c *RiceSackContract) Create(ctx contractapi.TransactionContextInterface, riceSackJSON string) error {
	riceSack, err := usermodel.UnmarshalRiceSack([]byte(riceSackJSON))
	if err != nil {
		return err
	}

	riceSackDoc := NewRiceSackDoc(*riceSack)
	riceSackDocJSON, err := json.Marshal(riceSackDoc)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(riceSackDoc.Code, riceSackDocJSON)
}

func getRiceSack(ctx contractapi.TransactionContextInterface, code string) (*usermodel.RiceSack, error) {
	riceSackJSON, err := ctx.GetStub().GetState(code)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if riceSackJSON == nil {
		return nil, nil
	}

	riceSack, err := usermodel.UnmarshalRiceSack(riceSackJSON)
	if err != nil {
		return nil, err
	}

	return riceSack, nil
}

func newDeterministicUuid(input string) (string, error) {
	id, err := uuid.NewRandomFromReader(strings.NewReader(fmt.Sprintf("%16s", input)))
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
