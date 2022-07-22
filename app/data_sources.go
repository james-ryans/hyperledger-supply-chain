package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/fabric"
)

type dataSources struct {
	Gateway *client.Gateway
	Couch   *kivik.Client
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
		return nil, fmt.Errorf("error opening fabric gateway: %w", err)
	}

	couchUrl := fmt.Sprintf("http://%s:%s@%s", os.Getenv("COUCHDB_USERNAME"), os.Getenv("COUCHDB_PASSWORD"), os.Getenv("COUCHDB_URL"))

	log.Printf("Connecting to CouchDB %s", couchUrl)
	couch, err := kivik.New("couch", couchUrl)
	if err != nil {
		return nil, fmt.Errorf("error opening couch db: %w", err)
	}

	_ = couch.CreateDB(context.TODO(), "users")

	return &dataSources{
		Gateway: gateway,
		Couch:   couch,
	}, nil
}
