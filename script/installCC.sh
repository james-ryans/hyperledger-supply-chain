#!/bin/bash

. ./script/utils.sh
. ./script/envVar.sh

CC_NAME=${1}
CC_SRC_PATH=${2}
CC_VERSION=${3}

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

vendorDep
packageChaincode
dropVendorDep

installChaincode supplier0.com
installChaincode producer0.com
installChaincode manufacturer0.com
installChaincode distributor0.com
installChaincode retailer0.com
