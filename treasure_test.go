package main

import (
	"encoding/json"
	"strconv"
	"testing"
)

var treasureID = TreasureID
var treasureAmount = DefaultTreasureAmount

func TestGetTreasure(t *testing.T) {
	t.Log("Test getTreasure")
	response := stub.MockInvoke("1", [][]byte{[]byte("getTreasure"), []byte(treasureID)})
	equals(t, int32(200), response.GetStatus())

	var treasure = new(Treasure)
	err := json.Unmarshal(response.GetPayload(), treasure)
	ok(t, err)
	equals(t, TreasureObjectType, treasure.ObjectType)
	equals(t, treasureAmount, strconv.FormatFloat(treasure.Balance, 'f', -1, 64))
}

func TestCreateTreasure(t *testing.T) {
	t.Log("Test createTreasure")
	treasureID = "test-ninjastack"
	treasureAmount = "5421000"
	response := stub.MockInvoke("1", [][]byte{[]byte("createTreasure"),
		[]byte(treasureAmount), []byte(treasureID)})
	equals(t, int32(200), response.GetStatus())

	var treasure = new(Treasure)
	err := json.Unmarshal(response.GetPayload(), treasure)
	ok(t, err)
	equals(t, TreasureObjectType, treasure.ObjectType)
	equals(t, treasureAmount, strconv.FormatFloat(treasure.Balance, 'f', -1, 64))

	// Test getTreasure again with new values
	TestGetTreasure(t)
}
