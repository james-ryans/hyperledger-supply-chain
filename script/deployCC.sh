#!/bin/bash

. ./script/utils.sh
. ./script/envVar.sh

DELAY=0.5
MAX_RETRY=10
CHANNEL_NAME=${1}
CC_NAME=${2}
CC_SRC_PATH=${3}
CC_VERSION=${4}
CC_SEQUENCE=${5}

function vendorDep() {
  infoln "Vendoring Go dependencies at ${CC_SRC_PATH}"
  pushd $CC_SRC_PATH
  GO111MODULE=on go mod vendor
  popd
  successln "Finished vendoring Go dependencies"
}

function packageChaincode() {
  setPeerEnv "superadmin.com" "SuperadminMSP" "localhost:5049"

  set -x
  peer lifecycle chaincode package "${CC_SRC_PATH}${CC_NAME}.tar.gz" --path ${CC_SRC_PATH} --lang golang --label ${CC_NAME}_${CC_VERSION} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode packaging has failed"
  successln "Chaincode is packaged"
}

function installChaincode() {
  ORG="${1}"
  setPeerEnv "${ORG}"

  set -x
  peer lifecycle chaincode install "${CC_SRC_PATH}${CC_NAME}.tar.gz" >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode installation on ${ORG} has failed"
  successln "Chaincode is installed on ${ORG}"
}

function queryInstalled() {
  ORG=$1
  setPeerEnv "${ORG}"

  set -x
  peer lifecycle chaincode queryinstalled >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
  verifyResult $res "Query installed on ${ORG} has failed"
  successln "Query installed successful on ${ORG} on channel"
}

function approveForMyOrg() {
  ORG=$1
  setPeerEnv "${ORG}"
  setOrdererEnv "${ORG}"

  set -x
  peer lifecycle chaincode approveformyorg -o "${ORDERER_ADDRESS}" --ordererTLSHostnameOverride "${ORDERER_HOSTNAME}" --tls --cafile "${ORDERER_CA}" --channelID "${CHANNEL_NAME}" --name "${CC_NAME}" --version "${CC_VERSION}" --package-id "${PACKAGE_ID}" --sequence ${CC_SEQUENCE} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode definition approved on ${ORG} on channel '$CHANNEL_NAME' failed"
  successln "Chaincode definition approved on ${ORG} on channel '$CHANNEL_NAME'"
}

function checkCommitReadiness() {
  ORG=$1
  shift 1
  setPeerEnv $ORG
  infoln "Checking the commit readiness of the chaincode definition on ${ORG} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  while [[ $rc -ne 0 ]] && [[ $COUNTER -lt $MAX_RETRY ]] ; do
    sleep $DELAY
    infoln "Attempting to check the commit readiness of the chaincode definition on ${ORG}, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode checkcommitreadiness --channelID ${CHANNEL_NAME} --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} --output json >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=0
    for var in "$@"; do
      grep "$var" log.txt &>/dev/null || let rc=1
    done
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
    infoln "Checking the commit readiness of the chaincode definition successful on ${ORG} on channel '$CHANNEL_NAME'"
  else
    fatalln "After $MAX_RETRY attempts, Check commit readiness result on ${ORG} is INVALID!"
  fi
}

function commitChaincodeDefinition() {
  ORG=$1
  setOrdererEnv $ORG
  parsePeerConnectionParameters "$@"
  res=$?
  verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

  set -x
  peer lifecycle chaincode commit -o "${ORDERER_ADDRESS}" --ordererTLSHostnameOverride "${ORDERER_HOSTNAME}" --tls --cafile "${ORDERER_CA}" --channelID "${CHANNEL_NAME}" --name "${CC_NAME}" "${PEER_CONN_PARMS[@]}" --version "${CC_VERSION}" --sequence "${CC_SEQUENCE}" >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode definition commit failed on ${ORG} on channel '$CHANNEL_NAME' failed"
  successln "Chaincode definition committed on channel '$CHANNEL_NAME'"
}

function queryCommitted() {
  ORG=$1
  setPeerEnv $ORG
  EXPECTED_RESULT="Version: ${CC_VERSION}, Sequence: ${CC_SEQUENCE}, Endorsement Plugin: escc, Validation Plugin: vscc"
  infoln "Querying chaincode definition on ${ORG} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  while [[ $rc -ne 0 ]] && [[ $COUNTER -lt $MAX_RETRY ]] ; do
    sleep $DELAY
    infoln "Attempting to Query committed status on ${ORG}, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME} >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    test $res -eq 0 && VALUE=$(cat log.txt | grep -o '^Version: '$CC_VERSION', Sequence: [0-9]*, Endorsement Plugin: escc, Validation Plugin: vscc')
    test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
    successln "Query chaincode definition successful on ${ORG} on channel '$CHANNEL_NAME'"
  else
    fatalln "After $MAX_RETRY attempts, Query chaincode definition result on ${ORG} is INVALID!"
  fi
}

vendorDep
packageChaincode

installChaincode supplier0.com
installChaincode producer0.com
installChaincode manufacturer0.com
installChaincode distributor0.com
installChaincode retailer0.com

queryInstalled supplier0.com
approveForMyOrg supplier0.com

#checkCommitReadiness supplier0.com "\"Producer0MSP\": false" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness producer0.com "\"Producer0MSP\": false" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness manufacturer0.com "\"Producer0MSP\": false" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness distributor0.com "\"Producer0MSP\": false" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness retailer0.com "\"Producer0MSP\": false" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"

queryInstalled producer0.com
approveForMyOrg producer0.com

#checkCommitReadiness supplier0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness producer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness manufacturer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness distributor0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness retailer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": false" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"

queryInstalled manufacturer0.com
approveForMyOrg manufacturer0.com

#checkCommitReadiness supplier0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness producer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness manufacturer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness distributor0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"
#checkCommitReadiness retailer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": false" "\"Retailer0MSP\": false"

queryInstalled distributor0.com
approveForMyOrg distributor0.com

#checkCommitReadiness supplier0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": false"
#checkCommitReadiness producer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": false"
#checkCommitReadiness manufacturer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": false"
#checkCommitReadiness distributor0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": false"
#checkCommitReadiness retailer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": false"

queryInstalled retailer0.com
approveForMyOrg retailer0.com

#checkCommitReadiness supplier0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": true"
#checkCommitReadiness producer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": true"
#checkCommitReadiness manufacturer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": true"
#checkCommitReadiness distributor0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": true"
#checkCommitReadiness retailer0.com "\"Producer0MSP\": true" "\"Supplier0MSP\": true" "\"Manufacturer0MSP\": true" "\"Distributor0MSP\": true" "\"Retailer0MSP\": true"

commitChaincodeDefinition supplier0.com \
  producer0.com \
  manufacturer0.com \
  distributor0.com \
  retailer0.com

queryCommitted supplier0.com
queryCommitted producer0.com
queryCommitted manufacturer0.com
queryCommitted distributor0.com
queryCommitted retailer0.com
