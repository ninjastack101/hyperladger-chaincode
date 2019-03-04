package main

import (
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Error declaration
var (
	errDoubleHit = errors.New("Transaction with given action and action entity already exists for the given wallet id")
)

type SmartContract struct {
}

const DefaultTreasureAmount = "210000000"
const DefaultRegistrationAmount = 110
const DefaultCustomer = "ninjastack"
const DefaultAction = "action"
const DefaultActionEntityId = "action-entity-id"

const OptionsObjectType = "options"
const WalletObjectType = "wallet"
const TreasureObjectType = "treasure"
const WalletTransactionObjectType = "walletTransaction"
const TreasureTransactionObjectType = "treasureTransaction"
const OptionsID = "Options"
const TreasureID = "Treasure"

// Init initializes chaincode
// ========================================
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	_, args := stub.GetFunctionAndParameters()
	var treasureAmount = DefaultTreasureAmount
	if len(args) > 0 {
		treasureAmount = args[0]
	}

	var registration float64 = DefaultRegistrationAmount
	if len(args) >= 2 {
		var value, err = strconv.ParseFloat(args[1], 64)
		if err == nil {
			registration = value
		}
	}

	s.createTreasure(stub, []string{treasureAmount})

	var options = make([]string, 2)
	options[0] = strconv.FormatFloat(registration, 'f', -1, 64)
	options[1] = DefaultCustomer
	return s.setOptions(stub, options)
}
