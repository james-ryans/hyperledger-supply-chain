package model

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/meneketehe/hehe/app/helper"
)

func NewConfigTx(name string, orgs []*GlobalOrganization) (string, error) {
	tempPath := filepath.Join(helper.BasePath, "app", "fixture", "configtx", "configtx-template.txt")
	temp, err := ioutil.ReadFile(tempPath)
	if err != nil {
		return "", err
	}

	orgconfigtx, err := NewOrgConfigTx(orgs)
	if err != nil {
		return "", err
	}

	ordAddConfigTx, err := NewOrdererAddressConfigTx(orgs)
	if err != nil {
		return "", err
	}

	ordConConfigTx, err := NewOrdererConsenterConfigTx(orgs)
	if err != nil {
		return "", err
	}

	orgAliasConfigTx, err := NewOrgAliasConfigTx(orgs)
	if err != nil {
		return "", err
	}

	configtx := fmt.Sprintf(string(temp), name, orgconfigtx, ordAddConfigTx, ordConConfigTx, orgAliasConfigTx)

	return configtx, nil
}

func NewOrgConfigTx(orgs []*GlobalOrganization) (string, error) {
	tempPath := filepath.Join(helper.BasePath, "app", "fixture", "configtx", "org-template.txt")
	temp, err := ioutil.ReadFile(tempPath)
	if err != nil {
		return "", err
	}

	configtx := ""
	for _, org := range orgs {
		mspDir := filepath.Join(helper.BasePath, "organizations", org.Domain, "msp")
		orgconfigtx := fmt.Sprintf(string(temp), org.MSPID, org.Domain, mspDir, org.Seq)

		configtx += orgconfigtx
	}

	return configtx, nil
}

func NewOrdererAddressConfigTx(orgs []*GlobalOrganization) (string, error) {
	tempPath := filepath.Join(helper.BasePath, "app", "fixture", "configtx", "ord-add-template.txt")
	temp, err := ioutil.ReadFile(tempPath)
	if err != nil {
		return "", err
	}

	configtx := ""
	for _, org := range orgs {
		orgconfigtx := fmt.Sprintf(string(temp), org.Domain, org.Seq)

		configtx += orgconfigtx
	}

	return configtx, nil
}

func NewOrdererConsenterConfigTx(orgs []*GlobalOrganization) (string, error) {
	tempPath := filepath.Join(helper.BasePath, "app", "fixture", "configtx", "ord-con-template.txt")
	temp, err := ioutil.ReadFile(tempPath)
	if err != nil {
		return "", err
	}

	configtx := ""
	for _, org := range orgs {
		orgconfigtx := fmt.Sprintf(string(temp), org.Domain, org.Seq, filepath.Join(helper.BasePath, "organizations", org.Domain, "orderers", "tls", "server.crt"))

		configtx += orgconfigtx
	}

	return configtx, nil
}

func NewOrgAliasConfigTx(orgs []*GlobalOrganization) (string, error) {
	tempPath := filepath.Join(helper.BasePath, "app", "fixture", "configtx", "org-alias-template.txt")
	temp, err := ioutil.ReadFile(tempPath)
	if err != nil {
		return "", err
	}

	configtx := ""
	for _, org := range orgs {
		orgconfigtx := fmt.Sprintf(string(temp), org.MSPID)

		configtx += orgconfigtx
	}

	return configtx, nil
}
