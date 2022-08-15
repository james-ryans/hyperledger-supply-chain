package fabric

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/meneketehe/hehe/app/helper"
)

type Gateway struct {
	Client *gateway.Gateway
}

func Connect(orgDomain string) (*gateway.Gateway, error) {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if orgDomain == "" {
		orgDomain = os.Getenv("ORG_DOMAIN")
	}
	if !wallet.Exists(orgDomain) {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(helper.BasePath, "organizations", orgDomain, "peers", "connection.yaml")

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, orgDomain),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}

	return gw, nil
}

func populateWallet(wallet *gateway.Wallet) error {
	orgDomain := os.Getenv("ORG_DOMAIN")
	credPath := filepath.Join(helper.BasePath, "organizations", orgDomain, "users", fmt.Sprintf("Admin@%s", orgDomain), "msp")

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	orgMSP := os.Getenv("FABRIC_MSP_ID")
	id := gateway.NewX509Identity(orgMSP, string(cert), string(key))

	return wallet.Put(orgDomain, id)
}
