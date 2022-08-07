#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORGDOTCOM}/$1/" \
        -e "s/\${ORG}/$2/" \
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
        -e "s/\${ORG}/$2/" \
        -e "s/\${P0PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        script/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORGDOTCOM=superadmin.com
ORG=Superadmin
P0PORT=5050
CAPORT=7054
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/superadmin.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/superadmin.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/superadmin.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/superadmin.com/peers/connection.yaml

ORGDOTCOM=supplier0.com
ORG=Supplier0
P0PORT=5052
CAPORT=7055
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/supplier0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/supplier0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/supplier0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/supplier0.com/peers/connection.yaml

ORGDOTCOM=producer0.com
ORG=Producer0
P0PORT=5054
CAPORT=7056
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/producer0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/producer0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/producer0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/producer0.com/peers/connection.yaml

ORGDOTCOM=manufacturer0.com
ORG=Manufacturer0
P0PORT=5056
CAPORT=7057
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/manufacturer0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/manufacturer0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/manufacturer0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/manufacturer0.com/peers/connection.yaml

ORGDOTCOM=distributor0.com
ORG=Distributor0
P0PORT=5058
CAPORT=7058
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/distributor0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/distributor0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/distributor0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/distributor0.com/peers/connection.yaml

ORGDOTCOM=retailer0.com
ORG=Retailer0
P0PORT=5060
CAPORT=7059
PEERPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/retailer0.com/fabric-ca/ca-cert.pem
CAPEM=/Users/james_ryans/go/src/github.com/meneketehe/hehe/organizations/retailer0.com/fabric-ca/ca-cert.pem

echo "$(json_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/retailer0.com/peers/connection.json
echo "$(yaml_ccp $ORGDOTCOM $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/retailer0.com/peers/connection.yaml
