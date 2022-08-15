#!/bin/bash

setPeerEnv() {
  local ORG=$1
  if [ "$1" = "superadmin.com" ]; then
    local MSPID="SuperadminMSP"
    local ADDRESS="localhost:5050"

  elif [ "$1" = "supplier0.com" ]; then
    local MSPID="Supplier0MSP"
    local ADDRESS="localhost:5051"

  elif [ "$1" = "producer0.com" ]; then
    local MSPID="Producer0MSP"
    local ADDRESS="localhost:5052"

  elif [ "$1" = "manufacturer0.com" ]; then
    local MSPID="Manufacturer0MSP"
    local ADDRESS="localhost:5053"

  elif [ "$1" = "distributor0.com" ]; then
    local MSPID="Distributor0MSP"
    local ADDRESS="localhost:5054"

  elif [ "$1" = "retailer0.com" ]; then
    local MSPID="Retailer0MSP"
    local ADDRESS="localhost:5055"
  fi

  export FABRIC_CFG_PATH="${PWD}/organizations/${ORG}/peercfg/"
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="${MSPID}"
  export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}/organizations/${ORG}/tlsca/tlsca.${ORG}-cert.pem"
  export CORE_PEER_MSPCONFIGPATH="${PWD}/organizations/${ORG}/users/Admin@${ORG}/msp"
  export CORE_PEER_ADDRESS="${ADDRESS}"
}

setOrdererEnv() {
  local ORG=$1
  if [ "$ORG" = "superadmin.com" ]; then
    local ADDRESS="localhost:6050"
    local HOSTNAME="orderer.superadmin.com"

  elif [ "$ORG" = "supplier0.com" ]; then
    local ADDRESS="localhost:6051"
    local HOSTNAME="orderer.supplier0.com"

  elif [ "$ORG" = "producer0.com" ]; then
    local ADDRESS="localhost:6052"
    local HOSTNAME="orderer.producer0.com"

  elif [ "$ORG" = "manufacturer0.com" ]; then
    local ADDRESS="localhost:6053"
    local HOSTNAME="orderer.manufacturer0.com"

  elif [ "$ORG" = "distributor0.com" ]; then
    local ADDRESS="localhost:6054"
    local HOSTNAME="orderer.distributor0.com"

  elif [ "$ORG" = "retailer0.com" ]; then
    local ADDRESS="localhost:6055"
    local HOSTNAME="orderer.retailer0.com"
  fi

  export ORDERER_ADDRESS="${ADDRESS}"
  export ORDERER_HOSTNAME="${HOSTNAME}"
  export ORDERER_CA="${PWD}/organizations/${ORG}/tlsca/tlsca.${ORG}-cert.pem"
}

parsePeerConnectionParameters() {
  PEER_CONN_PARMS=()
  PEERS=""
  while [ "$#" -gt 0 ]; do
    setPeerEnv $1
    PEER="$1"
    if [ -z "$PEERS" ]
    then
      PEERS="$PEER"
    else
      PEERS="$PEERS $PEER"
    fi
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" --peerAddresses "${CORE_PEER_ADDRESS}")
    CA="${CORE_PEER_TLS_ROOTCERT_FILE}"
    infoln $CA
    TLSINFO=(--tlsRootCertFiles "${CA}")
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" "${TLSINFO[@]}")
    shift
  done
}
