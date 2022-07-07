package main

import (
    "errors"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SimpleContract struct {
    contractapi.Contract
}

func (sc *SimpleContract) Create(ctx contractapi.TransactionContextInterface, key string, value string) error {
    existing, err := ctx.GetStub().GetState(key)
    if err != nil {
        return errors.New("unable to interact with world state")
    }

    if existing != nil {
        return fmt.Errorf("cannot create world state pair with key %s. Already exists", key)
    }

    err = ctx.GetStub().PutState(key, []byte(value))
    if err != nil {
        return errors.New("unable to interact with world state")
    }

    return nil
}

func (sc *SimpleContract) Update(ctx contractapi.TransactionContextInterface, key string, value string) error {
    existing, err := ctx.GetStub().GetState(key)
    if err != nil {
        return errors.New("unable to interact with world state")
    }

    if existing == nil {
        return fmt.Errorf("cannot update world state pair with key %s. Does not exist", key)
    }

    err = ctx.GetStub().PutState(key, []byte(value))
    if err != nil {
        return errors.New("unable to interact with world state")
    }

    return nil
}

func (sc *SimpleContract) Read(ctx contractapi.TransactionContextInterface, key string) (string, error) {
    existing, err := ctx.GetStub().GetState(key)
    if err != nil {
        return "", errors.New("unable to interact with world state")
    }

    if existing == nil {
        return "", fmt.Errorf("cannot read world state pair with key %s. Does not exist", key)
    }

    return string(existing), nil
}
