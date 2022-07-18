package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/globalchaincode/contract"
)

func main() {
	riceSackContract := &contract.RiceSackContract{}

	chaincode, err := contractapi.NewChaincode(riceSackContract)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
