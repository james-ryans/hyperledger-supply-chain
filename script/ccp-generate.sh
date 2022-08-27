#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORGDOTCOM}/$1/" \
        -e "s/\${ORGMSPID}/$2/" \
        -e "s/\${P0PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        script/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORGDOTCOM}/$1/" \
        -e "s/\${ORGMSPID}/$2/" \
        -e "s/\${P0PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        script/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORGDOTCOM=superadmin.com
ORGMSPID=SuperadminMSP
P0PORT=5050
CAPORT=7050
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/superadmin.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/superadmin.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/superadmin.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/superadmin.com/peers/connection.yaml

ORGDOTCOM=supplier0.com
ORGMSPID=Supplier0MSP
P0PORT=5051
CAPORT=7051
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/supplier0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/supplier0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/supplier0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/supplier0.com/peers/connection.yaml

ORGDOTCOM=producer0.com
ORGMSPID=Producer0MSP
P0PORT=5052
CAPORT=7052
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/producer0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/producer0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/producer0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/producer0.com/peers/connection.yaml

ORGDOTCOM=manufacturer0.com
ORGMSPID=Manufacturer0MSP
P0PORT=5053
CAPORT=7053
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/manufacturer0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/manufacturer0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/manufacturer0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/manufacturer0.com/peers/connection.yaml

ORGDOTCOM=distributor0.com
ORGMSPID=Distributor0MSP
P0PORT=5054
CAPORT=7054
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/distributor0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/distributor0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/distributor0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/distributor0.com/peers/connection.yaml

ORGDOTCOM=retailer0.com
ORGMSPID=Retailer0MSP
P0PORT=5055
CAPORT=7055
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/retailer0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/retailer0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/retailer0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORGMSPID $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/retailer0.com/peers/connection.yaml
