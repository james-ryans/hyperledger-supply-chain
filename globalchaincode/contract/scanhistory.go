package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type ScanHistoryContract struct {
	contractapi.Contract
}

type ScanHistoryDoc struct {
	DocType string `json:"doc_type"`
	usermodel.ScanHistory
}

func NewScanHistoryDoc(scanHistory usermodel.ScanHistory) ScanHistoryDoc {
	return ScanHistoryDoc{
		DocType:     "scanhistory",
		ScanHistory: scanHistory,
	}
}

func (c *ScanHistoryContract) FindAll(ctx contractapi.TransactionContextInterface, userId string) ([]*usermodel.ScanHistory, error) {
	query := fmt.Sprintf(`{"selector":{"doc_type":"scanhistory","user_id":"%s"}}`, userId)
	resultIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	scanHistories := make([]*usermodel.ScanHistory, 0)
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var scanHistory usermodel.ScanHistory
		err = json.Unmarshal(result.Value, &scanHistory)
		if err != nil {
			return nil, err
		}
		scanHistories = append(scanHistories, &scanHistory)
	}

	return scanHistories, nil
}
