package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/fabric"
)

type dataSources struct {
	Gateway *client.Gateway
}

func initDS() (*dataSources, error) {
	log.Printf("Initializing data sources\n")
	fabricCreds := fabric.Credentials{
		MSPID:        os.Getenv("FABRIC_MSP_ID"),
		PeerEndpoint: os.Getenv("FABRIC_PEER_ENDPOINT"),
		GatewayPeer:  os.Getenv("FABRIC_GATEWAY_PEER"),
		CertPath:     os.Getenv("FABRIC_CERT_PATH"),
		KeyPath:      os.Getenv("FABRIC_KEY_PATH"),
		TLSCertPath:  os.Getenv("FABRIC_TLS_CERT_PATH"),
	}
	fabricConfig := fabric.DefaultConfig()

	log.Printf("Connecting to Fabric Gateway\n")
	gateway, err := fabric.Connect(fabricCreds, fabricConfig)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	return &dataSources{
		Gateway: gateway,
	}, nil
}
