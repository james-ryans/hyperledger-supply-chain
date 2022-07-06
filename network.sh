#!/bin/bash

. utils.sh

function networkUp() {
  createSuperadminOrg
  createGlobalChannel

  createProducer0Org
  createSupplier0Org
  createManufacturer0Org
  createDistributor0Org
  createRetailer0Org
  createChannel0
}

function createSuperadminOrg() {
  infoln "Generating superadmin certificates using FABRIC CA"

  docker compose -f organizations/superadmin.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/superadmin.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . registerEnroll.sh
  enrollSuperadmin

  docker compose -f organizations/superadmin.com/docker-compose-orderer.yaml up -d 2>&1
}

function createProducer0Org() {
  infoln "Generating producer0 certificates using FABRIC CA"

  docker compose -f organizations/producer0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/producer0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . registerEnroll.sh
  enrollProducer0

  docker compose -f organizations/producer0.com/docker-compose-orderer.yaml \
    -f organizations/producer0.com/docker-compose-peer.yaml up -d 2>&1
}

function createSupplier0Org() {
  infoln "Generating supplier0 certificates using FABRIC CA"

  docker compose -f organizations/supplier0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/supplier0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . registerEnroll.sh
  enrollSupplier0

  docker compose -f organizations/supplier0.com/docker-compose-orderer.yaml \
    -f organizations/supplier0.com/docker-compose-peer.yaml up -d 2>&1
}

function createManufacturer0Org() {
  infoln "Generating manufacturer0 certificates using FABRIC CA"

  docker compose -f organizations/manufacturer0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/manufacturer0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . registerEnroll.sh
  enrollManufacturer0

  docker compose -f organizations/manufacturer0.com/docker-compose-orderer.yaml \
    -f organizations/manufacturer0.com/docker-compose-peer.yaml up -d 2>&1
}

function createDistributor0Org() {
  infoln "Generating distributor0 certificates using FABRIC CA"

  docker compose -f organizations/distributor0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/distributor0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . registerEnroll.sh
  enrollDistributor0

  docker compose -f organizations/distributor0.com/docker-compose-orderer.yaml \
    -f organizations/distributor0.com/docker-compose-peer.yaml up -d 2>&1
}

function createRetailer0Org() {
  infoln "Generating retailer0 certificates using FABRIC CA"

  docker compose -f organizations/retailer0.com/docker-compose-ca.yaml up -d 2>&1
  until [ -f "${PWD}/organizations/retailer0.com/fabric-ca/tls-cert.pem" ]; do
    sleep 1
  done;

  . registerEnroll.sh
  enrollRetailer0

  docker compose -f organizations/retailer0.com/docker-compose-orderer.yaml \
    -f organizations/retailer0.com/docker-compose-peer.yaml up -d 2>&1
}

function createGlobalChannel() {
  infoln "Creating Global Channel"

  export FABRIC_CFG_PATH=${PWD}/organizations/superadmin.com/global-channel-config

  set -x
  configtxgen -profile GlobalGenesis -outputBlock ./organizations/superadmin.com/channel-artifacts/globalchannel.block -channelID globalchannel
  { set +x; } 2>/dev/null

  local rc=1
  local COUNTER=1
  local DELAY=3
  local MAX_RETRY=10
  while [[ $rc -ne 0 ]] && [[ $COUNTER -lt $MAX_RETRY ]] ; do
    sleep $DELAY
    set -x
    osnadmin channel join --channelID globalchannel \
      --config-block ./organizations/superadmin.com/channel-artifacts/globalchannel.block \
      -o localhost:4050 \
      --ca-file "${PWD}/organizations/superadmin.com/tlsca/tlsca.superadmin.com-cert.pem" \
      --client-cert "${PWD}/organizations/superadmin.com/orderers/tls/server.crt" \
      --client-key "${PWD}/organizations/superadmin.com/orderers/tls/server.key" \
      >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=$res
    COUNTER=$(expr $COUNTER + 1)
  done

  cat log.txt
}

function createChannel0() {
  infoln "Creating Channel 0"

  export FABRIC_CFG_PATH=${PWD}/organizations/superadmin.com/channel0-config

  set -x
  configtxgen -profile Channel0Genesis -outputBlock ./organizations/superadmin.com/channel-artifacts/channel0.block -channelID channel0
  { set +x; } 2>/dev/null

  . createChannel.sh
  ordererJoinChannel channel0 localhost:4051 producer0.com
  ordererJoinChannel channel0 localhost:4052 supplier0.com
  ordererJoinChannel channel0 localhost:4053 manufacturer0.com
  ordererJoinChannel channel0 localhost:4054 distributor0.com
  ordererJoinChannel channel0 localhost:4055 retailer0.com

  peerJoinChannel channel0 localhost:5051 producer0.com Producer0MSP
  peerJoinChannel channel0 localhost:5053 supplier0.com Supplier0MSP
  peerJoinChannel channel0 localhost:5055 manufacturer0.com Manufacturer0MSP
  peerJoinChannel channel0 localhost:5057 distributor0.com Distributor0MSP
  peerJoinChannel channel0 localhost:5059 retailer0.com Retailer0MSP
}

function networkDown() {
  docker compose -f organizations/retailer0.com/docker-compose-ca.yaml \
    -f organizations/retailer0.com/docker-compose-orderer.yaml \
    -f organizations/retailer0.com/docker-compose-peer.yaml \
    down --volumes --remove-orphans

  docker compose -f organizations/distributor0.com/docker-compose-ca.yaml \
    -f organizations/distributor0.com/docker-compose-orderer.yaml \
    -f organizations/distributor0.com/docker-compose-peer.yaml \
    down --volumes --remove-orphans

  docker compose -f organizations/manufacturer0.com/docker-compose-ca.yaml \
    -f organizations/manufacturer0.com/docker-compose-orderer.yaml \
    -f organizations/manufacturer0.com/docker-compose-peer.yaml \
    down --volumes --remove-orphans

  docker compose -f organizations/supplier0.com/docker-compose-ca.yaml \
      -f organizations/supplier0.com/docker-compose-orderer.yaml \
      -f organizations/supplier0.com/docker-compose-peer.yaml \
      down --volumes --remove-orphans

  docker compose -f organizations/producer0.com/docker-compose-ca.yaml \
    -f organizations/producer0.com/docker-compose-orderer.yaml \
    -f organizations/producer0.com/docker-compose-peer.yaml \
    down --volumes --remove-orphans

  docker compose -f organizations/superadmin.com/docker-compose-ca.yaml \
    -f organizations/superadmin.com/docker-compose-orderer.yaml \
    down --volumes --remove-orphans

  rm -rf organizations/retailer0.com/fabric-ca \
    organizations/retailer0.com/msp \
    organizations/retailer0.com/orderers \
    organizations/retailer0.com/peers \
    organizations/retailer0.com/ca \
    organizations/retailer0.com/tlsca \
    organizations/retailer0.com/users \
    organizations/retailer0.com/fabric-ca-client-config.yaml

  rm -rf organizations/distributor0.com/fabric-ca \
    organizations/distributor0.com/msp \
    organizations/distributor0.com/orderers \
    organizations/distributor0.com/peers \
    organizations/distributor0.com/ca \
    organizations/distributor0.com/tlsca \
    organizations/distributor0.com/users \
    organizations/distributor0.com/fabric-ca-client-config.yaml

  rm -rf organizations/manufacturer0.com/fabric-ca \
    organizations/manufacturer0.com/msp \
    organizations/manufacturer0.com/orderers \
    organizations/manufacturer0.com/peers \
    organizations/manufacturer0.com/ca \
    organizations/manufacturer0.com/tlsca \
    organizations/manufacturer0.com/users \
    organizations/manufacturer0.com/fabric-ca-client-config.yaml

  rm -rf organizations/supplier0.com/fabric-ca \
    organizations/supplier0.com/msp \
    organizations/supplier0.com/orderers \
    organizations/supplier0.com/peers \
    organizations/supplier0.com/ca \
    organizations/supplier0.com/tlsca \
    organizations/supplier0.com/users \
    organizations/supplier0.com/fabric-ca-client-config.yaml

  rm -rf organizations/producer0.com/fabric-ca \
    organizations/producer0.com/msp \
    organizations/producer0.com/orderers \
    organizations/producer0.com/peers \
    organizations/producer0.com/ca \
    organizations/producer0.com/tlsca \
    organizations/producer0.com/users \
    organizations/producer0.com/fabric-ca-client-config.yaml

  rm -rf organizations/superadmin.com/channel-artifacts \
    organizations/superadmin.com/fabric-ca \
    organizations/superadmin.com/msp \
    organizations/superadmin.com/orderers \
    organizations/superadmin.com/tlsca \
    organizations/superadmin.com/users \
    organizations/superadmin.com/fabric-ca-client-config.yaml
}

MODE=$1

if [ "$MODE" == "up" ]; then
  infoln "Starting network"
  networkUp
elif [ "$MODE" == "down" ]; then
  infoln "Stopping network"
  networkDown
fi
