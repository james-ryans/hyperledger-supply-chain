package service

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/google/uuid"
	"github.com/meneketehe/hehe/app/fabric"
	"github.com/meneketehe/hehe/app/helper"
	"github.com/meneketehe/hehe/app/model"
)

type globalChannelService struct {
	GlobalChannelRepository model.GlobalChannelRepository
}

type GlobalChannelServiceConfig struct {
	GlobalChannelRepository model.GlobalChannelRepository
}

func NewGlobalChannelService(c *GlobalChannelServiceConfig) model.GlobalChannelService {
	return &globalChannelService{
		GlobalChannelRepository: c.GlobalChannelRepository,
	}
}

func (s *globalChannelService) GetChannelNameByFile(filename string) (string, error) {
	filePath := filepath.Join(helper.BasePath, "app", filename)
	JSON, err := exec.Command("configtxgen", "-inspectBlock", filePath).Output()
	if err != nil {
		return "", err
	}

	inspect, err := gabs.ParseJSON(JSON)
	if err != nil {
		return "", err
	}

	name, ok := inspect.Path("data.data.0.payload.header.channel_header.channel_id").Data().(string)
	if !ok {
		return "", fmt.Errorf("name not found in json")
	}

	return name, nil
}

func (s *globalChannelService) GetAllChannels() ([]*model.GlobalChannel, error) {
	chs, err := s.GlobalChannelRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return chs, nil
}

func (s *globalChannelService) GetChannel(ID string) (*model.GlobalChannel, error) {
	ch, err := s.GlobalChannelRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (s *globalChannelService) GetChannelByName(name string) (*model.GlobalChannel, error) {
	ch, err := s.GlobalChannelRepository.FindByName(name)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (s *globalChannelService) CheckNameExists(name string) (bool, error) {
	ch, err := s.GlobalChannelRepository.FindByName(name)
	if err != nil {
		return false, err
	}

	return ch != nil, nil
}

func (s *globalChannelService) CreateChannel(ch *model.GlobalChannel) (*model.GlobalChannel, error) {
	ch.ID = uuid.New().String()

	ch, err := s.GlobalChannelRepository.Create(ch)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (s *globalChannelService) CreateChannelBlock(ch *model.GlobalChannel, orgs []*model.GlobalOrganization) error {
	configtx, err := model.NewConfigTx(ch.Name, orgs)
	if err != nil {
		return err
	}

	configtxPath := filepath.Join(helper.BasePath, "app", "storage", ch.Name, "configtx.yaml")
	if err := os.MkdirAll(filepath.Dir(configtxPath), os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(configtxPath, []byte(configtx), 0644); err != nil {
		return err
	}

	outputPath := filepath.Join(helper.BasePath, "app", "storage", ch.Name, "genesis.block")
	if err := os.Setenv("FABRIC_CFG_PATH", filepath.Dir(outputPath)); err != nil {
		return err
	}
	cmd := exec.Command("configtxgen", "-profile", ch.Name+"Genesis", "-outputBlock", outputPath, "-channelID", ch.Name)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (s *globalChannelService) GetJoinedChannels() ([]string, error) {
	seq, _ := strconv.Atoi(os.Getenv("ORG_SEQ"))
	setPeerEnv(os.Getenv("ORG_DOMAIN"), os.Getenv("FABRIC_MSP_ID"), seq)

	out, _ := exec.Command("peer", "channel", "list").Output()
	log.Println(string(out))

	chs := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		ch := scanner.Text()
		if ch == "globalchannel" {
			continue
		}

		chs = append(chs, ch)
	}

	return chs[1:], nil
}

func (s *globalChannelService) JoinChannel(name, blockPath string, orgs []*model.GlobalOrganization) error {
	out, _ := exec.Command("exec", "$SHELL").CombinedOutput()
	log.Println(string(out))

	for _, org := range orgs {
		osnAdminJoin(org, name, blockPath)
		peerJoin(org, blockPath)
	}

	for _, org := range orgs {
		approveForMyOrg(org, name)
	}
	commitCCDef(orgs, name)

	for _, org := range orgs {
		gw, err := fabric.Connect(org.Domain)
		if err != nil {
			return fmt.Errorf("error opening fabric gateway: %w", err)
		}

		network, err := gw.GetNetwork(name)
		if err != nil {
			return fmt.Errorf("failed to get network %s: %w", name, err)
		}
		contract := network.GetContract("cc")

		switch org.Role {
		case "supplier":
			_, err = contract.SubmitTransaction(
				"SupplierContract:Create",
				org.ID,
				org.Role,
				org.Name,
				org.Location.Province,
				org.Location.City,
				org.Location.District,
				org.Location.PostalCode,
				org.Location.Address,
				org.ContactInfo.Phone,
				org.ContactInfo.Email,
				strconv.FormatFloat(float64(org.Location.Coordinate.Longitude), 'f', -1, 32),
				strconv.FormatFloat(float64(org.Location.Coordinate.Latitude), 'f', -1, 32),
			)
		case "producer":
			_, err = contract.SubmitTransaction(
				"ProducerContract:Create",
				org.ID,
				org.Role,
				org.Name,
				org.Location.Province,
				org.Location.City,
				org.Location.District,
				org.Location.PostalCode,
				org.Location.Address,
				org.ContactInfo.Phone,
				org.ContactInfo.Email,
				strconv.FormatFloat(float64(org.Location.Coordinate.Longitude), 'f', -1, 32),
				strconv.FormatFloat(float64(org.Location.Coordinate.Latitude), 'f', -1, 32),
			)
		case "manufacturer":
			_, err = contract.SubmitTransaction(
				"ManufacturerContract:Create",
				org.ID,
				org.Role,
				org.Name,
				org.Code,
				org.Location.Province,
				org.Location.City,
				org.Location.District,
				org.Location.PostalCode,
				org.Location.Address,
				org.ContactInfo.Phone,
				org.ContactInfo.Email,
				strconv.FormatFloat(float64(org.Location.Coordinate.Longitude), 'f', -1, 32),
				strconv.FormatFloat(float64(org.Location.Coordinate.Latitude), 'f', -1, 32),
			)
		case "distributor":
			_, err = contract.SubmitTransaction(
				"DistributorContract:Create",
				org.ID,
				org.Role,
				org.Name,
				org.Location.Province,
				org.Location.City,
				org.Location.District,
				org.Location.PostalCode,
				org.Location.Address,
				org.ContactInfo.Phone,
				org.ContactInfo.Email,
				strconv.FormatFloat(float64(org.Location.Coordinate.Longitude), 'f', -1, 32),
				strconv.FormatFloat(float64(org.Location.Coordinate.Latitude), 'f', -1, 32),
			)
		case "retailer":
			_, err = contract.SubmitTransaction(
				"RetailerContract:Create",
				org.ID,
				org.Role,
				org.Name,
				org.Location.Province,
				org.Location.City,
				org.Location.District,
				org.Location.PostalCode,
				org.Location.Address,
				org.ContactInfo.Phone,
				org.ContactInfo.Email,
				strconv.FormatFloat(float64(org.Location.Coordinate.Longitude), 'f', -1, 32),
				strconv.FormatFloat(float64(org.Location.Coordinate.Latitude), 'f', -1, 32),
			)
		}
		if err != nil {
			return fmt.Errorf("failed to submit transaction: %w", err)
		}
	}

	return nil
}

func osnAdminJoin(org *model.GlobalOrganization, name, blockPath string) {
	orgDomain := org.Domain
	orgSeq := org.Seq

	out, _ := exec.Command("osnadmin", "channel", "join",
		"--channelID", name,
		"--config-block", blockPath,
		"-o", "localhost:"+strconv.FormatInt(int64(4050+orgSeq), 10),
		"--ca-file", filepath.Join(helper.BasePath, "organizations", orgDomain, "tlsca", fmt.Sprintf("tlsca.%s-cert.pem", orgDomain)),
		"--client-cert", filepath.Join(helper.BasePath, "organizations", orgDomain, "orderers", "tls", "server.crt"),
		"--client-key", filepath.Join(helper.BasePath, "organizations", orgDomain, "orderers", "tls", "server.key"),
	).CombinedOutput()
	log.Println(string(out))
}

func peerJoin(org *model.GlobalOrganization, blockPath string) {
	setPeerEnv(org.Domain, org.MSPID, org.Seq)

	out, _ := exec.Command("peer", "channel", "join", "-b", blockPath).CombinedOutput()
	log.Println(string(out))
}

func approveForMyOrg(org *model.GlobalOrganization, name string) {
	orgDomain := org.Domain
	orgMSPID := org.MSPID
	orgSeq := org.Seq

	setPeerEnv(orgDomain, orgMSPID, orgSeq)

	out, _ := exec.Command("peer", "lifecycle", "chaincode", "queryinstalled", "--output", "json").Output()
	log.Println(string(out))
	installedCC, _ := gabs.ParseJSON(out)

	packageId := ""
	for idx := range installedCC.Path(fmt.Sprintf("installed_chaincodes")).Data().([]interface{}) {
		if installedCC.Path(fmt.Sprintf("installed_chaincodes.%d.label", idx)).Data().(string) == "cc_1.1" {
			packageId = installedCC.Path(fmt.Sprintf("installed_chaincodes.%d.package_id", idx)).Data().(string)
		}
	}
	log.Println(packageId)

	for true {
		out, _ = exec.Command("peer", "lifecycle", "chaincode", "approveformyorg",
			"-o", fmt.Sprintf("localhost:%d", 6050+orgSeq),
			"--ordererTLSHostnameOverride", fmt.Sprintf("orderer.%s", orgDomain),
			"--tls",
			"--cafile", filepath.Join(helper.BasePath, "organizations", orgDomain, "tlsca", fmt.Sprintf("tlsca.%s-cert.pem", orgDomain)),
			"--channelID", name,
			"--name", "cc",
			"--version", "1.1",
			"--package-id", packageId,
			"--sequence", "1",
		).CombinedOutput()
		log.Println(string(out))

		if string(out) != "Error: failed to send transaction: got unexpected status: SERVICE_UNAVAILABLE -- no Raft leader\n" {
			break
		}
	}
}

func commitCCDef(orgs []*model.GlobalOrganization, name string) {
	orgDomain := orgs[0].Domain
	orgSeq := orgs[0].Seq

	args := make([]string, 0)
	args = append(args,
		"lifecycle", "chaincode", "commit",
		"-o", fmt.Sprintf("localhost:%d", 6050+orgSeq),
		"--ordererTLSHostnameOverride", fmt.Sprintf("orderer.%s", orgDomain),
		"--tls",
		"--cafile", filepath.Join(helper.BasePath, "organizations", orgDomain, "tlsca", fmt.Sprintf("tlsca.%s-cert.pem", orgDomain)),
		"--channelID", name,
		"--name", "cc",
	)
	for _, org := range orgs {
		args = append(args, "--peerAddresses", fmt.Sprintf("localhost:%d", 5050+org.Seq))
		args = append(args, "--tlsRootCertFiles", filepath.Join(helper.BasePath, "organizations", org.Domain, "tlsca", fmt.Sprintf("tlsca.%s-cert.pem", org.Domain)))
	}
	args = append(args,
		"--version", "1.1",
		"--sequence", "1",
	)

	test := fmt.Sprintf("%s", "peer")
	for _, arg := range args {
		test += fmt.Sprintf(" %s", arg)
	}
	log.Println(test)

	out, _ := exec.Command("peer", args...).CombinedOutput()
	log.Println(string(out))
}

// func execCmd(c chan string, name string, args ...string) {
// 	out, _ := exec.Command(name, args...).CombinedOutput()
// 	c <- string(out)
// }

func setPeerEnv(orgDomain string, orgMSPID string, orgSeq int) {
	_ = os.Setenv("FABRIC_CFG_PATH", filepath.Join(helper.BasePath, "organizations", orgDomain, "peercfg"))
	_ = os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	_ = os.Setenv("CORE_PEER_LOCALMSPID", orgMSPID)
	_ = os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", filepath.Join(helper.BasePath, "organizations", orgDomain, "tlsca", fmt.Sprintf("tlsca.%s-cert.pem", orgDomain)))
	_ = os.Setenv("CORE_PEER_MSPCONFIGPATH", filepath.Join(helper.BasePath, "organizations", orgDomain, "users", fmt.Sprintf("Admin@%s", orgDomain), "msp"))
	_ = os.Setenv("CORE_PEER_ADDRESS", fmt.Sprintf("localhost:%d", 5050+orgSeq))
}
