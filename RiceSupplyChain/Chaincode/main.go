package main

import (
	"ricesupplychain/contracts"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func main() {

	riceContract := new(contracts.RiceContract)
		
    chaincode, err := contractapi.NewChaincode(riceContract)

    if err != nil {
        panic("Error creating chaincode: " + err.Error())
    }

    if err := chaincode.Start(); err != nil {
        panic("Error starting chaincode: " + err.Error())
    }
}

