#!/bin/bash

DELAY=3
MAX_RETRY=5

ordererJoinChannel() {
  CHANNEL_NAME="$1"
  ADDRESS="$2"
  ORG="$3"

  local rc=1
  local COUNTER=1
  while [[ $rc -ne 0 ]] && [[ $COUNTER -lt $MAX_RETRY ]] ; do
    sleep $DELAY
    set -x
    osnadmin channel join --channelID "${CHANNEL_NAME}" \
      --config-block "./organizations/superadmin.com/channel-artifacts/${CHANNEL_NAME}.block" \
      -o "${ADDRESS}" \
      --ca-file "${PWD}/organizations/${ORG}/tlsca/tlsca.${ORG}-cert.pem" \
      --client-cert "${PWD}/organizations/${ORG}/orderers/tls/server.crt" \
      --client-key "${PWD}/organizations/${ORG}/orderers/tls/server.key" \
      >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=$res
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
}

peerJoinChannel() {
  CHANNEL_NAME="$1"
  FABRIC_CFG_PATH="${PWD}/organizations/$3/peercfg/"
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="$4"
  export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}/organizations/$3/tlsca/tlsca.$3-cert.pem"
  export CORE_PEER_MSPCONFIGPATH="${PWD}/organizations/$3/users/Admin@$3/msp"
  export CORE_PEER_ADDRESS="$2"

  local rc=1
  local COUNTER=1
  while [[ $rc -ne 0 ]] && [[ $COUNTER -lt $MAX_RETRY ]] ; do
    sleep $DELAY
    set -x
    peer channel join -b "./organizations/superadmin.com/channel-artifacts/${CHANNEL_NAME}.block" >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=$res
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
}
