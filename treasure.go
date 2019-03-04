package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type Treasure struct {
	ObjectType string  `json:"docType"`
	Balance    float64 `json:"balance"`
}

type TreasureTransaction struct {
	ObjectType     string  `json:"docType"`
	TxID           string  `json:"txId"`
	Type           string  `json:"type"`
	Action         string  `json:"action"`
	ActionEntityID string  `json:"actionEntityId"`
	Amount         float64 `json:"amount"`
	CreationDate   int64   `json:"creationDate"`
	Customer       string  `json:"customer"`
}

func (s *SmartContract) createTreasure(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	var balance float64
	if len(args) >= 1 && args[0] != "" {
		b, err := strconv.ParseFloat(args[0], 64)
		if err == nil {
			balance = b
		}
	}

	TreasureID := TreasureID
	if len(args) >= 2 {
		TreasureID = args[1]
	}

	var treasure = new(Treasure)
	treasure.ObjectType = TreasureObjectType
	treasure.Balance = balance

	key, err := stub.CreateCompositeKey(TreasureObjectType, []string{TreasureID})
	if err != nil {
		return shim.Error(err.Error())
	}

	treasureAsBytes, err := json.Marshal(treasure)
	if err != nil {
		return shim.Error(err.Error())
	}

	stub.PutState(key, treasureAsBytes)

	uuid := DefaultActionEntityId

	err = s.createTreasureTransaction(stub, balance, "createTreasure", stub.GetTxID(), "genesis transaction", uuid, DefaultCustomer)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(treasureAsBytes)
}

func (s *SmartContract) getTreasure(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	TreasureID := TreasureID
	if len(args) >= 1 {
		TreasureID = args[0]
	}

	key, err := stub.CreateCompositeKey(TreasureObjectType, []string{TreasureID})
	if err != nil {
		return shim.Error(err.Error())
	}

	treasureAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(treasureAsBytes)
}

func (s *SmartContract) updateTreasureBalance(stub shim.ChaincodeStubInterface,
	amount float64,
	transactionType, txnID, action, actionEntityID, customer string) error {

	key, err := stub.CreateCompositeKey(TreasureObjectType, []string{TreasureID})
	if err != nil {
		return err
	}

	var treasure = new(Treasure)
	treasureAsBytes, err := stub.GetState(key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(treasureAsBytes, treasure)
	if err != nil {
		return err
	}

	treasure.Balance += amount
	if treasure.Balance < 0 {
		return errors.New("insufficient funds on treasure")
	}

	err = s.createTreasureTransaction(stub, amount, transactionType, txnID, action, actionEntityID, customer)
	if err != nil {
		return err
	}

	treasureAsBytes, err = json.Marshal(treasure)
	if err != nil {
		return err
	}

	err = stub.PutState(key, treasureAsBytes)
	return err
}

func (s *SmartContract) createTreasureTransaction(stub shim.ChaincodeStubInterface,
	amount float64,
	transactionType, txnID, action, actionEntityID, customer string) error {

	var key, err = stub.CreateCompositeKey(TreasureTransactionObjectType, []string{txnID, transactionType})
	if err != nil {
		return err
	}

	var transaction = new(TreasureTransaction)
	transaction.ObjectType = TreasureTransactionObjectType
	transaction.Type = transactionType
	transaction.TxID = txnID
	transaction.Action = action
	transaction.ActionEntityID = actionEntityID
	transaction.Amount = amount
	transaction.Customer = customer
	transaction.CreationDate = time.Now().Unix()
	trAsBytes, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	err = stub.PutState(key, trAsBytes)
	return err
}
