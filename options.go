package main

import (
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type Options struct {
	ObjectType   string  `json:"docType"`
	Registration float64 `json:"registration"`
	Customer     string  `json:"customer"`
}

func (s *SmartContract) setOptions(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var options = new(Options)
	options.ObjectType = OptionsObjectType
	options.Registration, _ = strconv.ParseFloat(args[0], 64)
	options.Customer = DefaultCustomer
	if len(args) >= 2 {
		options.Customer = args[1]
	}

	var key, err = stub.CreateCompositeKey(OptionsObjectType, []string{OptionsID, options.Customer})
	if err != nil {
		return shim.Error(err.Error())
	}

	asBytes, err := json.Marshal(options)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(key, asBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(asBytes)
}

func (s *SmartContract) getOptionsAsByte(stub shim.ChaincodeStubInterface, customer string) ([]byte, error) {

	if len(customer) == 0 {
		customer = DefaultCustomer
	}

	key, err := stub.CreateCompositeKey(OptionsObjectType, []string{OptionsID, customer})
	if err != nil {
		return nil, err
	}

	options, err := stub.GetState(key)
	if options == nil && customer != DefaultCustomer {
		return s.getOptionsAsByte(stub, "")
	}

	return options, err
}

func (s *SmartContract) getOptions(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	customer := DefaultCustomer
	if len(args) >= 1 {
		customer = args[0]
	}

	options, err := s.getOptionsAsByte(stub, customer)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(options)
}

func (s *SmartContract) getOptionsObject(stub shim.ChaincodeStubInterface, customer string) (Options, error) {

	var options = new(Options)
	var optionsBytes, err = s.getOptionsAsByte(stub, customer)
	if err != nil {
		return *options, err
	}

	err = json.Unmarshal(optionsBytes, options)
	return *options, err
}
