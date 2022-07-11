package fabric

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/meneketehe/hehe/app/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Credentials struct {
	MSPID        string
	PeerEndpoint string
	GatewayPeer  string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
}

type Config struct {
	EvaluateTimeout     int32
	EndorseTimeout      int32
	SubmitTimeout       int32
	CommitStatusTimeout int32
}

type Gateway struct {
	*Config
	Client *client.Gateway
}

func Connect(cred Credentials, conf Config) (*client.Gateway, error) {
	clientConnection, err := newGrpcConnection(cred)
	if err != nil {
		return nil, err
	}

	id, err := newIdentity(cred)
	if err != nil {
		return nil, err
	}

	sign, err := newSign(cred)
	if err != nil {
		return nil, err
	}

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(time.Duration(conf.EvaluateTimeout)*time.Second),
		client.WithEndorseTimeout(time.Duration(conf.EndorseTimeout)*time.Second),
		client.WithSubmitTimeout(time.Duration(conf.SubmitTimeout)*time.Second),
		client.WithCommitStatusTimeout(time.Duration(conf.CommitStatusTimeout)*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return gateway, nil
}

func DefaultConfig() Config {
	return Config{
		EvaluateTimeout:     5,
		EndorseTimeout:      15,
		SubmitTimeout:       5,
		CommitStatusTimeout: 60,
	}
}

func newGrpcConnection(cred Credentials) (*grpc.ClientConn, error) {
	certificate, err := loadCertificate(filepath.Join(helper.BasePath, cred.TLSCertPath))
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, cred.GatewayPeer)

	connection, err := grpc.Dial(cred.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return connection, nil
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePem, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePem)
}

func newIdentity(cred Credentials) (*identity.X509Identity, error) {
	certificate, err := loadCertificate(filepath.Join(helper.BasePath, cred.CertPath))
	if err != nil {
		return nil, err
	}

	id, err := identity.NewX509Identity(cred.MSPID, certificate)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func newSign(cred Credentials) (identity.Sign, error) {
	files, err := ioutil.ReadDir(filepath.Join(helper.BasePath, cred.KeyPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read private key directory: %w", err)
	}
	privateKeyPEM, err := ioutil.ReadFile(path.Join(helper.BasePath, cred.KeyPath, files[0].Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, err
	}

	return sign, nil
}
