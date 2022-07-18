#!/bin/bash

function networkUp() {
  createSuperadminOrg
  createSupplier0Org
  createProducer0Org
  createManufacturer0Org
  createDistributor0Org
  createRetailer0Org

  createChannel0
  createGlobalChannel
}

function createSuperadminOrg() {
  infoln "Generating superadmin certificates using FABRIC CA"

  docker compose -f organizations/superadmin.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/superadmin.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . ./script/registerEnroll.sh
  enrollSuperadmin

  docker compose -f organizations/superadmin.com/docker-compose-orderer.yaml \
    -f organizations/superadmin.com/docker-compose-couch.yaml \
    -f organizations/superadmin.com/docker-compose-peer.yaml up -d 2>&1
}

function createSupplier0Org() {
  infoln "Generating supplier0 certificates using FABRIC CA"

  docker compose -f organizations/supplier0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/supplier0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . ./script/registerEnroll.sh
  enrollSupplier0

  docker compose -f organizations/supplier0.com/docker-compose-orderer.yaml \
    -f organizations/supplier0.com/docker-compose-couch.yaml \
    -f organizations/supplier0.com/docker-compose-peer.yaml up -d 2>&1
}

function createProducer0Org() {
  infoln "Generating producer0 certificates using FABRIC CA"

  docker compose -f organizations/producer0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/producer0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . ./script/registerEnroll.sh
  enrollProducer0

  docker compose -f organizations/producer0.com/docker-compose-orderer.yaml \
    -f organizations/producer0.com/docker-compose-couch.yaml \
    -f organizations/producer0.com/docker-compose-peer.yaml up -d 2>&1
}

function createManufacturer0Org() {
  infoln "Generating manufacturer0 certificates using FABRIC CA"

  docker compose -f organizations/manufacturer0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/manufacturer0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . ./script/registerEnroll.sh
  enrollManufacturer0

  docker compose -f organizations/manufacturer0.com/docker-compose-orderer.yaml \
    -f organizations/manufacturer0.com/docker-compose-couch.yaml \
    -f organizations/manufacturer0.com/docker-compose-peer.yaml up -d 2>&1
}

function createDistributor0Org() {
  infoln "Generating distributor0 certificates using FABRIC CA"

  docker compose -f organizations/distributor0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/distributor0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . ./script/registerEnroll.sh
  enrollDistributor0

  docker compose -f organizations/distributor0.com/docker-compose-orderer.yaml \
    -f organizations/distributor0.com/docker-compose-couch.yaml \
    -f organizations/distributor0.com/docker-compose-peer.yaml up -d 2>&1
}

function createRetailer0Org() {
  infoln "Generating retailer0 certificates using FABRIC CA"

  docker compose -f organizations/retailer0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/retailer0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . ./script/registerEnroll.sh
  enrollRetailer0

  docker compose -f organizations/retailer0.com/docker-compose-orderer.yaml \
    -f organizations/retailer0.com/docker-compose-couch.yaml \
    -f organizations/retailer0.com/docker-compose-peer.yaml up -d 2>&1
}

function createGlobalChannel() {
  infoln "Creating Global Channel"

  export FABRIC_CFG_PATH=${PWD}/organizations/superadmin.com/global-channel-config

  set -x
  configtxgen -profile GlobalGenesis -outputBlock ./organizations/superadmin.com/channel-artifacts/globalchannel.block -channelID globalchannel
  { set +x; } 2>/dev/null

  . ./script/createChannel.sh
  ordererJoinChannel globalchannel superadmin.com localhost:4050

  peerJoinChannel globalchannel superadmin.com SuperadminMSP localhost:5050
  peerJoinChannel globalchannel retailer0.com Retailer0MSP localhost:5060
}

function createChannel0() {
  infoln "Creating Channel 0"

  export FABRIC_CFG_PATH=${PWD}/organizations/superadmin.com/channel0-config

  set -x
  configtxgen -profile Channel0Genesis -outputBlock ./organizations/superadmin.com/channel-artifacts/channel0.block -channelID channel0
  { set +x; } 2>/dev/null

  . ./script/createChannel.sh
  ordererJoinChannel channel0 supplier0.com localhost:4051
  ordererJoinChannel channel0 producer0.com localhost:4052
  ordererJoinChannel channel0 manufacturer0.com localhost:4053
  ordererJoinChannel channel0 distributor0.com localhost:4054
  ordererJoinChannel channel0 retailer0.com localhost:4055

  peerJoinChannel channel0 supplier0.com Supplier0MSP localhost:5052
  peerJoinChannel channel0 producer0.com Producer0MSP localhost:5054
  peerJoinChannel channel0 manufacturer0.com Manufacturer0MSP localhost:5056
  peerJoinChannel channel0 distributor0.com Distributor0MSP localhost:5058
  peerJoinChannel channel0 retailer0.com Retailer0MSP localhost:5060
}

function deployCC() {
  ./script/deployCC.sh "${1}" "${2}" "${3}" "${4}" "${5}"
}

function deployGCC() {
  ./script/deployGCC.sh "${1}" "${2}" "${3}" "${4}" "${5}"
}

function networkDown() {
  downDockerContainers retailer0.com
  downDockerContainers distributor0.com
  downDockerContainers manufacturer0.com
  downDockerContainers producer0.com
  downDockerContainers supplier0.com
  downDockerContainers superadmin.com

  removeGeneratedFiles retailer0.com
  removeGeneratedFiles distributor0.com
  removeGeneratedFiles manufacturer0.com
  removeGeneratedFiles producer0.com
  removeGeneratedFiles supplier0.com
  removeGeneratedFiles superadmin.com
}

function downDockerContainers() {
  ORG=$1

  docker compose -f organizations/"$ORG"/docker-compose-ca.yaml \
    -f organizations/"$ORG"/docker-compose-orderer.yaml \
    -f organizations/"$ORG"/docker-compose-couch.yaml \
    -f organizations/"$ORG"/docker-compose-peer.yaml \
    down --volumes --remove-orphans
}

function removeGeneratedFiles() {
  ORG=$1

  rm -rf organizations/"$ORG"/channel-artifacts \
    organizations/"$ORG"/fabric-ca/msp \
    organizations/"$ORG"/fabric-ca/ca-cert.pem \
    organizations/"$ORG"/fabric-ca/tls-cert.pem \
    organizations/"$ORG"/fabric-ca/fabric-ca-server.db \
    organizations/"$ORG"/fabric-ca/IssuerPublicKey \
    organizations/"$ORG"/fabric-ca/IssuerRevocationPublicKey \
    organizations/"$ORG"/msp \
    organizations/"$ORG"/orderers \
    organizations/"$ORG"/peers \
    organizations/"$ORG"/ca \
    organizations/"$ORG"/tlsca \
    organizations/"$ORG"/users \
    organizations/"$ORG"/fabric-ca-client-config.yaml
}

cd "$PWD/.." || exit

. ./script/utils.sh

MODE=$1

if [ "$MODE" == "up" ]; then
  infoln "Starting network"
  networkUp
elif [ "$MODE" == "down" ]; then
  infoln "Stopping network"
  networkDown
elif [ "$MODE" == "deployCC" ]; then
  infoln "Deploying chaincode on org channel"
  deployCC "${2}" "${3}" "${4}" "${5}" "${6}"
elif [ "$MODE" == "deployGCC" ]; then
  infoln "Deploying chaincode on global"
  deployGCC "${2}" "${3}" "${4}" "${5}" "${6}"
fi
