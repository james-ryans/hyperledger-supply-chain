package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/meneketehe/hehe/chaincode/contract"
)

func main() {
	supplierContract := &contract.SupplierContract{}
	producerContract := &contract.ProducerContract{}
	manufacturerContract := &contract.ManufacturerContract{}
	distributorContract := &contract.DistributorContract{}
	retailerContract := &contract.RetailerContract{}

	seedContract := &contract.SeedContract{}
	riceGrainContract := &contract.RiceGrainContract{}
	riceContract := &contract.RiceContract{}

	chaincode, err := contractapi.NewChaincode(
		supplierContract,
		producerContract,
		manufacturerContract,
		distributorContract,
		retailerContract,
		seedContract,
		riceGrainContract,
		riceContract,
	)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
