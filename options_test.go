package main

import (
	"encoding/json"
	"strconv"
	"testing"
)

var defaultCustomer = DefaultCustomer
var defaultRegistration float64 = DefaultRegistrationAmount

func TestGetOptions(t *testing.T) {
	t.Log("Test getOptions")
	response := stub.MockInvoke("1", [][]byte{[]byte("getOptions"),
		[]byte(defaultCustomer)})
	equals(t, int32(200), response.GetStatus())

	var options = new(Options)
	err := json.Unmarshal(response.GetPayload(), options)
	ok(t, err)
	equals(t, defaultRegistration, options.Registration)
	equals(t, defaultCustomer, options.Customer)
}

func TestSetOptions(t *testing.T) {
	t.Log("Test setOptions")
	defaultCustomer = "New-ninjastack"
	defaultRegistration = float64(200)
	response := stub.MockInvoke("2", [][]byte{[]byte("setOptions"),
		[]byte(strconv.FormatFloat(defaultRegistration, 'f', -1, 64)),
		[]byte(defaultCustomer)})
	equals(t, int32(200), response.GetStatus())

	var options = new(Options)
	err := json.Unmarshal(response.GetPayload(), options)
	ok(t, err)
	equals(t, defaultRegistration, options.Registration)
	equals(t, defaultCustomer, options.Customer)

	// Test getOptions again with new values
	TestGetOptions(t)
}

// ------------------------------------- Negative Cases --------------------------------------------------------

func TestFalseCase(t *testing.T) {
	t.Log("Test false cases")
	response := stub.MockInvoke("3", [][]byte{[]byte("getOption"),
		[]byte(defaultCustomer)})
	equals(t, int32(500), response.GetStatus())

	response = stub.MockInvoke("4", [][]byte{[]byte("setOptions")})
	equals(t, int32(500), response.GetStatus())
}
