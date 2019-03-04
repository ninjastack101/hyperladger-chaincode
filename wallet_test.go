package main

import (
	"encoding/json"
	"testing"
)

var defaultWalletID = "default_wallet_id"
var defaultMobileHash = "9ba3878af953abfc6d91e50f01d7bded407c601101497235dd4b60a20b25ecee"

func TestCreateWallet(t *testing.T) {
	t.Log("Test createWallet")
	response := stub.MockInvoke("1", [][]byte{[]byte("createWallet"),
		[]byte(defaultWalletID),
		[]byte(defaultMobileHash)})
	equals(t, int32(200), response.GetStatus())

	var wallet = new(Wallet)
	err := json.Unmarshal(response.GetPayload(), wallet)
	ok(t, err)
	equals(t, DefaultRegistrationAmount, int(wallet.Amount))
	equals(t, defaultWalletID, wallet.ID)
	equals(t, defaultMobileHash, wallet.MobileHash)
}

func TestGetWallet(t *testing.T) {
	t.Log("Test getWallet")
	response := stub.MockInvoke("1", [][]byte{[]byte("getWallet"), []byte(defaultWalletID)})
	equals(t, int32(200), response.GetStatus())

	var wallet = new(Wallet)
	err := json.Unmarshal(response.GetPayload(), wallet)
	ok(t, err)
	equals(t, DefaultRegistrationAmount, int(wallet.Amount))
	equals(t, defaultWalletID, wallet.ID)
	equals(t, defaultMobileHash, wallet.MobileHash)
}

func TestPurchaseCoins(t *testing.T) {
	t.Log("Test purchaseCoins")
	response := stub.MockInvoke("1", [][]byte{[]byte("purchaseCoins"),
		[]byte(defaultWalletID),
		[]byte("200"),
		[]byte("MAGIC_BOX"),
		[]byte("BOX_NUMBER_1")})
	equals(t, int32(200), response.GetStatus())

	response = stub.MockInvoke("1", [][]byte{[]byte("getWallet"), []byte(defaultWalletID)})
	equals(t, int32(200), response.GetStatus())

	var wallet = new(Wallet)
	err := json.Unmarshal(response.GetPayload(), wallet)
	ok(t, err)
	equals(t, 310, int(wallet.Amount))
	equals(t, defaultWalletID, wallet.ID)
	equals(t, defaultMobileHash, wallet.MobileHash)
}

func TestSpendCoins(t *testing.T) {
	t.Log("Test spendCoins")
	response := stub.MockInvoke("1", [][]byte{[]byte("spendCoins"),
		[]byte(defaultWalletID),
		[]byte("200"),
		[]byte("PREDICTION"),
		[]byte("P_NUMBER_1")})
	equals(t, int32(200), response.GetStatus())

	response = stub.MockInvoke("1", [][]byte{[]byte("getWallet"), []byte(defaultWalletID)})
	equals(t, int32(200), response.GetStatus())

	var wallet = new(Wallet)
	err := json.Unmarshal(response.GetPayload(), wallet)
	ok(t, err)
	equals(t, DefaultRegistrationAmount, int(wallet.Amount))
	equals(t, defaultWalletID, wallet.ID)
	equals(t, defaultMobileHash, wallet.MobileHash)
}

// ------------------------------------- Negative Cases --------------------------------------------------------

func TestPurchaseCoinsNegative(t *testing.T) {
	t.Log("Test purchaseCoins Negative")
	response := stub.MockInvoke("1", [][]byte{[]byte("purchaseCoins"),
		[]byte(defaultWalletID),
		[]byte("200"),
		[]byte("MAGIC_BOX"),
		[]byte("BOX_NUMBER_1")})
	equals(t, int32(409), response.GetStatus())

	response = stub.MockInvoke("1", [][]byte{[]byte("getWallet"), []byte(defaultWalletID)})
	equals(t, int32(200), response.GetStatus())

	var wallet = new(Wallet)
	err := json.Unmarshal(response.GetPayload(), wallet)
	ok(t, err)
	equals(t, DefaultRegistrationAmount, int(wallet.Amount))
	equals(t, defaultWalletID, wallet.ID)
	equals(t, defaultMobileHash, wallet.MobileHash)
}

func TestSpendCoinsNegative(t *testing.T) {
	t.Log("Test spendCoins Negative")
	response := stub.MockInvoke("1", [][]byte{[]byte("spendCoins"),
		[]byte(defaultWalletID),
		[]byte("50"),
		[]byte("PREDICTION"),
		[]byte("P_NUMBER_1")})
	equals(t, int32(409), response.GetStatus())

	response = stub.MockInvoke("1", [][]byte{[]byte("getWallet"), []byte(defaultWalletID)})
	equals(t, int32(200), response.GetStatus())

	var wallet = new(Wallet)
	err := json.Unmarshal(response.GetPayload(), wallet)
	ok(t, err)
	equals(t, DefaultRegistrationAmount, int(wallet.Amount))
	equals(t, defaultWalletID, wallet.ID)
	equals(t, defaultMobileHash, wallet.MobileHash)
}
