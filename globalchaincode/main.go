package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/globalchaincode/contract"
)

func main() {
	userContract := &contract.UserContract{}
	riceSackContract := &contract.RiceSackContract{}
	scanHistoryContract := &contract.ScanHistoryContract{}
	commentContract := &contract.CommentContract{}

	chaincode, err := contractapi.NewChaincode(
		userContract,
		riceSackContract,
		scanHistoryContract,
		commentContract,
	)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
