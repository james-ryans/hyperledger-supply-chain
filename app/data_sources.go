package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/meneketehe/hehe/app/fabric"
	"github.com/meneketehe/hehe/app/model"
	"github.com/meneketehe/hehe/app/model/enum"
	service "github.com/meneketehe/hehe/app/service/organization"
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

	if os.Getenv("ORG_ROLE") == "superadmin" {
		_ = couch.CreateDB(context.TODO(), "users")
		_ = couch.CreateDB(context.TODO(), "organizations")
		_ = couch.CreateDB(context.TODO(), "channels")
	}

	err = couch.CreateDB(context.TODO(), "admins")
	if err == nil {
		db := couch.DB(context.TODO(), "admins")

		hashedPassword, err := service.HashPassword(os.Getenv("ORG_PASSWORD"))
		if err != nil {
			return nil, fmt.Errorf("error hash password: %v", err)
		}

		account := model.OrganizationAccount{
			ID:             os.Getenv("ORG_ADMIN_ID"),
			Role:           os.Getenv("ORG_ROLE"),
			Type:           enum.OrgAccAdmin,
			OrganizationID: os.Getenv("ORG_ID"),
			Email:          os.Getenv("ORG_EMAIL"),
			Password:       hashedPassword,
			RegisteredAt:   time.Now(),
		}
		_, err = db.Put(context.TODO(), account.ID, account)
		if err != nil {
			return nil, fmt.Errorf("error seeding admin account: %v", err)
		}
	}

	return &dataSources{
		Gateway: gateway,
		Couch:   couch,
	}, nil
}
