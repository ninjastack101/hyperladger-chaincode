package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type Wallet struct {
	ObjectType string  `json:"docType"`
	ID         string  `json:"id"`
	Amount     float64 `json:"amount"`
	MobileHash string  `json:"mobileHash"`
}

type WalletTransaction struct {
	ObjectType     string  `json:"docType"`
	WalletID       string  `json:"walletId"`
	TxID           string  `json:"txId"`
	Type           string  `json:"type"`
	Action         string  `json:"action"`
	ActionEntityID string  `json:"actionEntityId"`
	Amount         float64 `json:"amount"`
	CreationDate   int64   `json:"creationDate"`
	Customer       string  `json:"customer"`
}

func (s *SmartContract) createWallet(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	argsLength := len(args)
	if argsLength < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	customer := DefaultCustomer
	if argsLength >= 6 {
		customer = args[5]
	}

	action := DefaultAction
	actionEntityID := DefaultActionEntityId

	if argsLength >= 5 {
		action = args[3]
		actionEntityID = args[4]
	}

	var amount float64
	if argsLength > 2 && args[2] != "" {
		val, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return shim.Error(err.Error())
		}
		amount = val

	} else {
		options, err := s.getOptionsObject(stub, customer)
		if err != nil {
			return shim.Error(err.Error())
		}
		amount = options.Registration

	}

	var wallet = new(Wallet)
	wallet.ObjectType = WalletObjectType
	wallet.ID = args[0]
	wallet.MobileHash = args[1]
	wallet.Amount = amount

	Key, err := stub.CreateCompositeKey(WalletObjectType, []string{wallet.ID})
	if err != nil {
		return shim.Error(err.Error())
	}

	asBytes, err := stub.GetState(Key)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(asBytes) != 0 {
		return shim.Error("Wallet with id " + args[0] + " already exists")
	}

	err = s.updateTreasureBalance(stub, -amount, "registration", stub.GetTxID(), action, actionEntityID, customer)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = s.createWalletTransaction(stub, amount, wallet.ID, "registration", stub.GetTxID(), action, actionEntityID, customer)
	if err != nil {
		if err == errDoubleHit {
			return sc.Response{
				Status:  int32(409),
				Message: err.Error(),
			}
		}
		return shim.Error(err.Error())
	}

	walletAsBytes, err := json.Marshal(wallet)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(Key, walletAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(walletAsBytes)
}

func (s *SmartContract) getWallet(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key, err := stub.CreateCompositeKey(WalletObjectType, []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}

	walletAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(walletAsBytes)
}

func (s *SmartContract) purchaseCoins(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	customer := DefaultCustomer
	if len(args) >= 5 {
		customer = args[4]
	}

	var walletID = args[0]
	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	action := args[2]
	actionEntityID := args[3]

	err = s.updateTreasureBalance(stub, -amount, "purchase", stub.GetTxID(), action, actionEntityID, customer)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = s.updateWalletBalance(stub, amount, walletID, "purchase", stub.GetTxID(), action, actionEntityID, customer)
	if err != nil {
		if err == errDoubleHit {
			return sc.Response{
				Status:  int32(409),
				Message: err.Error(),
			}
		}
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) spendCoins(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	customer := DefaultCustomer
	if len(args) >= 5 {
		customer = args[4]
	}

	var walletID = args[0]
	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	action := args[2]
	actionEntityID := args[3]

	err = s.updateWalletBalance(stub, -amount, walletID, "spend", stub.GetTxID(), action, actionEntityID, customer)
	if err != nil {
		if err == errDoubleHit {
			return sc.Response{
				Status:  int32(409),
				Message: err.Error(),
			}
		}
		return shim.Error(err.Error())
	}

	err = s.updateTreasureBalance(stub, amount, "spend", stub.GetTxID(), action, actionEntityID, customer)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) updateWalletBalance(stub shim.ChaincodeStubInterface,
	amount float64,
	walletID, transactionType, txnID, action, actionEntityID, customer string) error {

	var key, err = stub.CreateCompositeKey(WalletObjectType, []string{walletID})
	if err != nil {
		return err
	}

	var wallet = new(Wallet)
	byteWallets, err := stub.GetState(key)
	if err != nil {
		return err
	}

	if len(byteWallets) == 0 {
		return errors.New("Wallet with id " + walletID + " not found")
	}

	err = json.Unmarshal(byteWallets, wallet)
	if err != nil {
		return err
	}

	wallet.Amount = wallet.Amount + amount
	if wallet.Amount < 0 {
		return errors.New("insufficient funds")
	}

	err = s.createWalletTransaction(stub, amount, walletID, transactionType, txnID, action, actionEntityID, customer)
	if err != nil {
		return err
	}

	asBytes, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	err = stub.PutState(key, asBytes)
	return err
}

func (s *SmartContract) updateWalletMobileHash(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	walletID := args[0]
	mobileHash := args[1]
	key, err := stub.CreateCompositeKey(WalletObjectType, []string{walletID})
	if err != nil {
		return shim.Error(err.Error())
	}

	action := "MOBILE_UPDATE"
	actionEntityID := mobileHash + " | " + time.Now().String()

	walletAsBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(walletAsBytes) == 0 {
		return shim.Error("Wallet with id " + walletID + " not found")
	}

	var wallet = new(Wallet)
	err = json.Unmarshal(walletAsBytes, wallet)
	if err != nil {
		return shim.Error(err.Error())
	}

	wallet.MobileHash = mobileHash

	err = s.createWalletTransaction(stub, 0, walletID, "mobile update", stub.GetTxID(), action, actionEntityID, DefaultCustomer)
	if err != nil {
		if err == errDoubleHit {
			return sc.Response{
				Status:  int32(409),
				Message: err.Error(),
			}
		}
	}

	walletAsBytes, err = json.Marshal(wallet)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(key, walletAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(walletAsBytes)
}

func (s *SmartContract) createWalletTransaction(stub shim.ChaincodeStubInterface,
	amount float64,
	walletID, transactionType, txnID, action, actionEntityID, customer string) error {
	var key, err = stub.CreateCompositeKey(WalletTransactionObjectType, []string{walletID, action, actionEntityID})
	if err != nil {
		return err
	}

	walletTxn, err := stub.GetState(key)
	if len(walletTxn) != 0 {
		return errDoubleHit
	}

	var transaction = new(WalletTransaction)
	transaction.ObjectType = WalletTransactionObjectType
	transaction.TxID = txnID
	transaction.Type = transactionType
	transaction.WalletID = walletID
	transaction.Amount = amount
	transaction.Action = action
	transaction.ActionEntityID = actionEntityID
	transaction.Customer = customer
	transaction.CreationDate = time.Now().Unix()
	trAsBytes, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	err = stub.PutState(key, trAsBytes)
	return err
}
