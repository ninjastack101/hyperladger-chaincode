package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Invoke - Our entry point for Invocations
// ========================================
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	logger.Info(fmt.Sprintf("Starting ninjastackcoin smart contract Invoke for %s and arguments passed are %v", function, args))

	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	case "createWallet":
		return s.createWallet(stub, args)
	case "getWallet":
		return s.getWallet(stub, args)
	case "searchWallets":
		return s.searchWallets(stub, args)
	case "updateWalletMobileHash":
		return s.updateWalletMobileHash(stub, args)
	case "searchWalletTransactions":
		return s.searchWalletTransactions(stub, args)
	case "searchTreasureTransactions":
		return s.searchTreasureTransactions(stub, args)
	case "createTreasure":
		return s.createTreasure(stub, args)
	case "getTreasure":
		return s.getTreasure(stub, args)
	case "purchaseCoins":
		return s.purchaseCoins(stub, args)
	case "spendCoins":
		return s.spendCoins(stub, args)
	case "setOptions":
		return s.setOptions(stub, args)
	case "getOptions":
		return s.getOptions(stub, args)
	default:
		return shim.Error("Invalid Smart contract function name.")
	}
}
