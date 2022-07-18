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

function dropVendorDep() {
  infoln "Drop vendored Go dependencies at ${CC_SRC_PATH}"
  rm -rf "${CC_SRC_PATH}/vendor"
  successln "Finished drop vendored Go dependencies"
}

function packageChaincode() {
  setPeerEnv "superadmin.com"

  find "${CC_SRC_PATH}" -type f -iname '*.tar.gz' -delete

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
  peer lifecycle chaincode approveformyorg  -o localhost:6050 --ordererTLSHostnameOverride orderer.superadmin.com --tls --cafile /Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/superadmin.com/tlsca/tlsca.superadmin.com-cert.pem --channelID "${CHANNEL_NAME}" --name "${CC_NAME}" --version "${CC_VERSION}" --package-id "${PACKAGE_ID}" --sequence ${CC_SEQUENCE} >&log.txt
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
dropVendorDep

installChaincode superadmin.com
installChaincode retailer0.com

queryInstalled superadmin.com
approveForMyOrg superadmin.com

queryInstalled retailer0.com
approveForMyOrg retailer0.com

commitChaincodeDefinition superadmin.com \
  retailer0.com

queryCommitted superadmin.com
queryCommitted retailer0.com
