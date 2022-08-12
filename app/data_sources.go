package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/meneketehe/hehe/app/fabric"
	"github.com/meneketehe/hehe/app/model"
	"github.com/meneketehe/hehe/app/model/enum"
	service "github.com/meneketehe/hehe/app/service/organization"
)

type dataSources struct {
	Gateway *gateway.Gateway
	Couch   *kivik.Client
}

func initDS() (*dataSources, error) {
	log.Printf("Initializing data sources\n")

	log.Printf("Connecting to Fabric Gateway\n")
	gw, err := fabric.Connect()
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
			Code:           os.Getenv("ORG_CODE"),
			Name:           os.Getenv("ORG_NAME"),
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
		Gateway: gw,
		Couch:   couch,
	}, nil
}
