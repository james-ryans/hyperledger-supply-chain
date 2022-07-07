#!/bin/bash

function enrollSuperadmin() {
  infoln "Enrolling the CA admin"

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/superadmin.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-superadmin --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-superadmin.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-superadmin.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-superadmin.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-superadmin.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/superadmin.com/msp/config.yaml"

  mkdir -p "${PWD}/organizations/superadmin.com/msp/tlscacerts"
  cp "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/superadmin.com/msp/tlscacerts/tlsca.superadmin.com-cert.pem"

  mkdir -p "${PWD}/organizations/superadmin.com/tlsca"
  cp "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/superadmin.com/tlsca/tlsca.superadmin.com-cert.pem"

  mkdir -p "${PWD}/organizations/superadmin.com/ca"
  cp "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/superadmin.com/ca/ca.superadmin.com-cert.pem"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-superadmin --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering peer"
  set -x
  fabric-ca-client register --caname ca-superadmin --id.name peer --id.secret peerpw --id.type peer --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-superadmin --id.name user --id.secret userpw --id.type client --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-superadmin --id.name superadminadmin --id.secret superadminadminpw --id.type admin --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7054 --caname ca-superadmin -M "${PWD}/organizations/superadmin.com/orderers/msp" --csr.hosts orderer.superadmin.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/superadmin.com/msp/config.yaml" "${PWD}/organizations/superadmin.com/orderers/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7054 --caname ca-superadmin -M "${PWD}/organizations/superadmin.com/orderers/tls" --enrollment.profile tls --csr.hosts orderer.superadmin.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/superadmin.com/orderers/tls/tlscacerts/"* "${PWD}/organizations/superadmin.com/orderers/tls/ca.crt"
  cp "${PWD}/organizations/superadmin.com/orderers/tls/signcerts/"* "${PWD}/organizations/superadmin.com/orderers/tls/server.crt"
  cp "${PWD}/organizations/superadmin.com/orderers/tls/keystore/"* "${PWD}/organizations/superadmin.com/orderers/tls/server.key"

  infoln "Generating the peer msp"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7054 --caname ca-superadmin -M "${PWD}/organizations/superadmin.com/peers/msp" --csr.hosts peer.superadmin.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/superadmin.com/msp/config.yaml" "${PWD}/organizations/superadmin.com/peers/msp/config.yaml"

  infoln "Generating the peer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7054 --caname ca-superadmin -M "${PWD}/organizations/superadmin.com/peers/tls" --enrollment.profile tls --csr.hosts peer.superadmin.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/superadmin.com/peers/tls/tlscacerts/"* "${PWD}/organizations/superadmin.com/peers/tls/ca.crt"
  cp "${PWD}/organizations/superadmin.com/peers/tls/signcerts/"* "${PWD}/organizations/superadmin.com/peers/tls/server.crt"
  cp "${PWD}/organizations/superadmin.com/peers/tls/keystore/"* "${PWD}/organizations/superadmin.com/peers/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user:userpw@localhost:7054 --caname ca-superadmin -M "${PWD}/organizations/superadmin.com/users/User@superadmin.com/msp" --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/superadmin.com/msp/config.yaml" "${PWD}/organizations/superadmin.com/peers/users/User@superadmin.com/msp/config.yaml"

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://superadminadmin:superadminadminpw@localhost:7054 --caname ca-superadmin -M "${PWD}/organizations/superadmin.com/users/Admin@superadmin.com/msp" --tls.certfiles "${PWD}/organizations/superadmin.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/superadmin.com/msp/config.yaml" "${PWD}/organizations/superadmin.com/users/Admin@superadmin.com/msp/config.yaml"
}

function enrollProducer0() {
  infoln "Enrolling the CA admin"

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/producer0.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7055 --caname ca-producer0 --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7055-ca-producer0.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7055-ca-producer0.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7055-ca-producer0.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7055-ca-producer0.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/producer0.com/msp/config.yaml"

  mkdir -p "${PWD}/organizations/producer0.com/msp/tlscacerts"
  cp "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/producer0.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/producer0.com/tlsca"
  cp "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/producer0.com/tlsca/tlsca.producer0.com-cert.pem"

  mkdir -p "${PWD}/organizations/producer0.com/ca"
  cp "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/producer0.com/ca/ca.producer0.com-cert.pem"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-producer0 --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering peer"
  set -x
  fabric-ca-client register --caname ca-producer0 --id.name peer --id.secret peerpw --id.type peer --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-producer0 --id.name user --id.secret userpw --id.type client --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-producer0 --id.name producer0admin --id.secret producer0adminpw --id.type admin --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7055 --caname ca-producer0 -M "${PWD}/organizations/producer0.com/orderers/msp" --csr.hosts orderer.producer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/producer0.com/msp/config.yaml" "${PWD}/organizations/producer0.com/orderers/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7055 --caname ca-producer0 -M "${PWD}/organizations/producer0.com/orderers/tls" --enrollment.profile tls --csr.hosts orderer.producer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/producer0.com/orderers/tls/tlscacerts/"* "${PWD}/organizations/producer0.com/orderers/tls/ca.crt"
  cp "${PWD}/organizations/producer0.com/orderers/tls/signcerts/"* "${PWD}/organizations/producer0.com/orderers/tls/server.crt"
  cp "${PWD}/organizations/producer0.com/orderers/tls/keystore/"* "${PWD}/organizations/producer0.com/orderers/tls/server.key"

  infoln "Generating the peer msp"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7055 --caname ca-producer0 -M "${PWD}/organizations/producer0.com/peers/msp" --csr.hosts peer.producer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/producer0.com/msp/config.yaml" "${PWD}/organizations/producer0.com/peers/msp/config.yaml"

  infoln "Generating the peer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7055 --caname ca-producer0 -M "${PWD}/organizations/producer0.com/peers/tls" --enrollment.profile tls --csr.hosts peer.producer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/producer0.com/peers/tls/tlscacerts/"* "${PWD}/organizations/producer0.com/peers/tls/ca.crt"
  cp "${PWD}/organizations/producer0.com/peers/tls/signcerts/"* "${PWD}/organizations/producer0.com/peers/tls/server.crt"
  cp "${PWD}/organizations/producer0.com/peers/tls/keystore/"* "${PWD}/organizations/producer0.com/peers/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user:userpw@localhost:7055 --caname ca-producer0 -M "${PWD}/organizations/producer0.com/users/User@producer0.com/msp" --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/producer0.com/msp/config.yaml" "${PWD}/organizations/producer0.com/peers/users/User@producer0.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://producer0admin:producer0adminpw@localhost:7055 --caname ca-producer0 -M "${PWD}/organizations/producer0.com/users/Admin@producer0.com/msp" --tls.certfiles "${PWD}/organizations/producer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/producer0.com/msp/config.yaml" "${PWD}/organizations/producer0.com/users/Admin@producer0.com/msp/config.yaml"
}

function enrollSupplier0() {
  infoln "Enrolling the CA admin"

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/supplier0.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7056 --caname ca-supplier0 --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7056-ca-supplier0.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7056-ca-supplier0.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7056-ca-supplier0.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7056-ca-supplier0.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/supplier0.com/msp/config.yaml"

  mkdir -p "${PWD}/organizations/supplier0.com/msp/tlscacerts"
  cp "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/supplier0.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/supplier0.com/tlsca"
  cp "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/supplier0.com/tlsca/tlsca.supplier0.com-cert.pem"

  mkdir -p "${PWD}/organizations/supplier0.com/ca"
  cp "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/supplier0.com/ca/ca.supplier0.com-cert.pem"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-supplier0 --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering peer"
  set -x
  fabric-ca-client register --caname ca-supplier0 --id.name peer --id.secret peerpw --id.type peer --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-supplier0 --id.name user --id.secret userpw --id.type client --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-supplier0 --id.name supplier0admin --id.secret supplier0adminpw --id.type admin --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7056 --caname ca-supplier0 -M "${PWD}/organizations/supplier0.com/orderers/msp" --csr.hosts orderer.supplier0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/supplier0.com/msp/config.yaml" "${PWD}/organizations/supplier0.com/orderers/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7056 --caname ca-supplier0 -M "${PWD}/organizations/supplier0.com/orderers/tls" --enrollment.profile tls --csr.hosts orderer.supplier0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/supplier0.com/orderers/tls/tlscacerts/"* "${PWD}/organizations/supplier0.com/orderers/tls/ca.crt"
  cp "${PWD}/organizations/supplier0.com/orderers/tls/signcerts/"* "${PWD}/organizations/supplier0.com/orderers/tls/server.crt"
  cp "${PWD}/organizations/supplier0.com/orderers/tls/keystore/"* "${PWD}/organizations/supplier0.com/orderers/tls/server.key"

  infoln "Generating the peer msp"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7056 --caname ca-supplier0 -M "${PWD}/organizations/supplier0.com/peers/msp" --csr.hosts peer.supplier0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/supplier0.com/msp/config.yaml" "${PWD}/organizations/supplier0.com/peers/msp/config.yaml"

  infoln "Generating the peer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7056 --caname ca-supplier0 -M "${PWD}/organizations/supplier0.com/peers/tls" --enrollment.profile tls --csr.hosts peer.supplier0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/supplier0.com/peers/tls/tlscacerts/"* "${PWD}/organizations/supplier0.com/peers/tls/ca.crt"
  cp "${PWD}/organizations/supplier0.com/peers/tls/signcerts/"* "${PWD}/organizations/supplier0.com/peers/tls/server.crt"
  cp "${PWD}/organizations/supplier0.com/peers/tls/keystore/"* "${PWD}/organizations/supplier0.com/peers/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user:userpw@localhost:7056 --caname ca-supplier0 -M "${PWD}/organizations/supplier0.com/users/User@supplier0.com/msp" --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/supplier0.com/msp/config.yaml" "${PWD}/organizations/supplier0.com/peers/users/User@supplier0.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://supplier0admin:supplier0adminpw@localhost:7056 --caname ca-supplier0 -M "${PWD}/organizations/supplier0.com/users/Admin@supplier0.com/msp" --tls.certfiles "${PWD}/organizations/supplier0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/supplier0.com/msp/config.yaml" "${PWD}/organizations/supplier0.com/users/Admin@supplier0.com/msp/config.yaml"
}

function enrollManufacturer0() {
  infoln "Enrolling the CA admin"

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/manufacturer0.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7057 --caname ca-manufacturer0 --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7057-ca-manufacturer0.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7057-ca-manufacturer0.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7057-ca-manufacturer0.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7057-ca-manufacturer0.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/manufacturer0.com/msp/config.yaml"

  mkdir -p "${PWD}/organizations/manufacturer0.com/msp/tlscacerts"
  cp "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/manufacturer0.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/manufacturer0.com/tlsca"
  cp "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/manufacturer0.com/tlsca/tlsca.manufacturer0.com-cert.pem"

  mkdir -p "${PWD}/organizations/manufacturer0.com/ca"
  cp "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/manufacturer0.com/ca/ca.manufacturer0.com-cert.pem"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-manufacturer0 --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering peer"
  set -x
  fabric-ca-client register --caname ca-manufacturer0 --id.name peer --id.secret peerpw --id.type peer --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-manufacturer0 --id.name user --id.secret userpw --id.type client --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-manufacturer0 --id.name manufacturer0admin --id.secret manufacturer0adminpw --id.type admin --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7057 --caname ca-manufacturer0 -M "${PWD}/organizations/manufacturer0.com/orderers/msp" --csr.hosts orderer.manufacturer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/manufacturer0.com/msp/config.yaml" "${PWD}/organizations/manufacturer0.com/orderers/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7057 --caname ca-manufacturer0 -M "${PWD}/organizations/manufacturer0.com/orderers/tls" --enrollment.profile tls --csr.hosts orderer.manufacturer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/manufacturer0.com/orderers/tls/tlscacerts/"* "${PWD}/organizations/manufacturer0.com/orderers/tls/ca.crt"
  cp "${PWD}/organizations/manufacturer0.com/orderers/tls/signcerts/"* "${PWD}/organizations/manufacturer0.com/orderers/tls/server.crt"
  cp "${PWD}/organizations/manufacturer0.com/orderers/tls/keystore/"* "${PWD}/organizations/manufacturer0.com/orderers/tls/server.key"

  infoln "Generating the peer msp"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7057 --caname ca-manufacturer0 -M "${PWD}/organizations/manufacturer0.com/peers/msp" --csr.hosts peer.manufacturer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/manufacturer0.com/msp/config.yaml" "${PWD}/organizations/manufacturer0.com/peers/msp/config.yaml"

  infoln "Generating the peer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7057 --caname ca-manufacturer0 -M "${PWD}/organizations/manufacturer0.com/peers/tls" --enrollment.profile tls --csr.hosts peer.manufacturer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/manufacturer0.com/peers/tls/tlscacerts/"* "${PWD}/organizations/manufacturer0.com/peers/tls/ca.crt"
  cp "${PWD}/organizations/manufacturer0.com/peers/tls/signcerts/"* "${PWD}/organizations/manufacturer0.com/peers/tls/server.crt"
  cp "${PWD}/organizations/manufacturer0.com/peers/tls/keystore/"* "${PWD}/organizations/manufacturer0.com/peers/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user:userpw@localhost:7057 --caname ca-manufacturer0 -M "${PWD}/organizations/manufacturer0.com/users/User@manufacturer0.com/msp" --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/manufacturer0.com/msp/config.yaml" "${PWD}/organizations/manufacturer0.com/peers/users/User@manufacturer0.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://manufacturer0admin:manufacturer0adminpw@localhost:7057 --caname ca-manufacturer0 -M "${PWD}/organizations/manufacturer0.com/users/Admin@manufacturer0.com/msp" --tls.certfiles "${PWD}/organizations/manufacturer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/manufacturer0.com/msp/config.yaml" "${PWD}/organizations/manufacturer0.com/users/Admin@manufacturer0.com/msp/config.yaml"
}

function enrollDistributor0() {
  infoln "Enrolling the CA admin"

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/distributor0.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7058 --caname ca-distributor0 --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7058-ca-distributor0.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7058-ca-distributor0.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7058-ca-distributor0.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7058-ca-distributor0.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/distributor0.com/msp/config.yaml"

  mkdir -p "${PWD}/organizations/distributor0.com/msp/tlscacerts"
  cp "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/distributor0.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/distributor0.com/tlsca"
  cp "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/distributor0.com/tlsca/tlsca.distributor0.com-cert.pem"

  mkdir -p "${PWD}/organizations/distributor0.com/ca"
  cp "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/distributor0.com/ca/ca.distributor0.com-cert.pem"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-distributor0 --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering peer"
  set -x
  fabric-ca-client register --caname ca-distributor0 --id.name peer --id.secret peerpw --id.type peer --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-distributor0 --id.name user --id.secret userpw --id.type client --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-distributor0 --id.name distributor0admin --id.secret distributor0adminpw --id.type admin --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7058 --caname ca-distributor0 -M "${PWD}/organizations/distributor0.com/orderers/msp" --csr.hosts orderer.distributor0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/distributor0.com/msp/config.yaml" "${PWD}/organizations/distributor0.com/orderers/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7058 --caname ca-distributor0 -M "${PWD}/organizations/distributor0.com/orderers/tls" --enrollment.profile tls --csr.hosts orderer.distributor0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/distributor0.com/orderers/tls/tlscacerts/"* "${PWD}/organizations/distributor0.com/orderers/tls/ca.crt"
  cp "${PWD}/organizations/distributor0.com/orderers/tls/signcerts/"* "${PWD}/organizations/distributor0.com/orderers/tls/server.crt"
  cp "${PWD}/organizations/distributor0.com/orderers/tls/keystore/"* "${PWD}/organizations/distributor0.com/orderers/tls/server.key"

  infoln "Generating the peer msp"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7058 --caname ca-distributor0 -M "${PWD}/organizations/distributor0.com/peers/msp" --csr.hosts peer.distributor0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/distributor0.com/msp/config.yaml" "${PWD}/organizations/distributor0.com/peers/msp/config.yaml"

  infoln "Generating the peer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7058 --caname ca-distributor0 -M "${PWD}/organizations/distributor0.com/peers/tls" --enrollment.profile tls --csr.hosts peer.distributor0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/distributor0.com/peers/tls/tlscacerts/"* "${PWD}/organizations/distributor0.com/peers/tls/ca.crt"
  cp "${PWD}/organizations/distributor0.com/peers/tls/signcerts/"* "${PWD}/organizations/distributor0.com/peers/tls/server.crt"
  cp "${PWD}/organizations/distributor0.com/peers/tls/keystore/"* "${PWD}/organizations/distributor0.com/peers/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user:userpw@localhost:7058 --caname ca-distributor0 -M "${PWD}/organizations/distributor0.com/users/User@distributor0.com/msp" --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/distributor0.com/msp/config.yaml" "${PWD}/organizations/distributor0.com/peers/users/User@distributor0.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://distributor0admin:distributor0adminpw@localhost:7058 --caname ca-distributor0 -M "${PWD}/organizations/distributor0.com/users/Admin@distributor0.com/msp" --tls.certfiles "${PWD}/organizations/distributor0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/distributor0.com/msp/config.yaml" "${PWD}/organizations/distributor0.com/users/Admin@distributor0.com/msp/config.yaml"
}

function enrollRetailer0() {
  infoln "Enrolling the CA admin"

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/retailer0.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7059 --caname ca-retailer0 --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7059-ca-retailer0.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7059-ca-retailer0.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7059-ca-retailer0.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7059-ca-retailer0.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/retailer0.com/msp/config.yaml"

  mkdir -p "${PWD}/organizations/retailer0.com/msp/tlscacerts"
  cp "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/retailer0.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/retailer0.com/tlsca"
  cp "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/retailer0.com/tlsca/tlsca.retailer0.com-cert.pem"

  mkdir -p "${PWD}/organizations/retailer0.com/ca"
  cp "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem" "${PWD}/organizations/retailer0.com/ca/ca.retailer0.com-cert.pem"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-retailer0 --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering peer"
  set -x
  fabric-ca-client register --caname ca-retailer0 --id.name peer --id.secret peerpw --id.type peer --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-retailer0 --id.name user --id.secret userpw --id.type client --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-retailer0 --id.name retailer0admin --id.secret retailer0adminpw --id.type admin --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7059 --caname ca-retailer0 -M "${PWD}/organizations/retailer0.com/orderers/msp" --csr.hosts orderer.retailer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/retailer0.com/msp/config.yaml" "${PWD}/organizations/retailer0.com/orderers/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:7059 --caname ca-retailer0 -M "${PWD}/organizations/retailer0.com/orderers/tls" --enrollment.profile tls --csr.hosts orderer.retailer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/retailer0.com/orderers/tls/tlscacerts/"* "${PWD}/organizations/retailer0.com/orderers/tls/ca.crt"
  cp "${PWD}/organizations/retailer0.com/orderers/tls/signcerts/"* "${PWD}/organizations/retailer0.com/orderers/tls/server.crt"
  cp "${PWD}/organizations/retailer0.com/orderers/tls/keystore/"* "${PWD}/organizations/retailer0.com/orderers/tls/server.key"

  infoln "Generating the peer msp"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7059 --caname ca-retailer0 -M "${PWD}/organizations/retailer0.com/peers/msp" --csr.hosts peer.retailer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/retailer0.com/msp/config.yaml" "${PWD}/organizations/retailer0.com/peers/msp/config.yaml"

  infoln "Generating the peer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer:peerpw@localhost:7059 --caname ca-retailer0 -M "${PWD}/organizations/retailer0.com/peers/tls" --enrollment.profile tls --csr.hosts peer.retailer0.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/retailer0.com/peers/tls/tlscacerts/"* "${PWD}/organizations/retailer0.com/peers/tls/ca.crt"
  cp "${PWD}/organizations/retailer0.com/peers/tls/signcerts/"* "${PWD}/organizations/retailer0.com/peers/tls/server.crt"
  cp "${PWD}/organizations/retailer0.com/peers/tls/keystore/"* "${PWD}/organizations/retailer0.com/peers/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user:userpw@localhost:7059 --caname ca-retailer0 -M "${PWD}/organizations/retailer0.com/users/User@retailer0.com/msp" --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/retailer0.com/msp/config.yaml" "${PWD}/organizations/retailer0.com/peers/users/User@retailer0.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://retailer0admin:retailer0adminpw@localhost:7059 --caname ca-retailer0 -M "${PWD}/organizations/retailer0.com/users/Admin@retailer0.com/msp" --tls.certfiles "${PWD}/organizations/retailer0.com/fabric-ca/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/retailer0.com/msp/config.yaml" "${PWD}/organizations/retailer0.com/users/Admin@retailer0.com/msp/config.yaml"
}
