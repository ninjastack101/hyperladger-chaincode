package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) searchWallets(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	return s.searchEntities(stub, WalletObjectType, args)
}

func (s *SmartContract) searchWalletTransactions(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	return s.searchEntities(stub, WalletTransactionObjectType, args)
}

func (s *SmartContract) searchTreasureTransactions(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	return s.searchEntities(stub, TreasureTransactionObjectType, args)
}

func (s *SmartContract) searchEntities(stub shim.ChaincodeStubInterface, DocType string, args []string) sc.Response {

	var parametersString = ""
	if len(args) > 0 && args[0] != "" {
		parametersString = parametersString + ", " + args[0]
	}
	limit := 10
	skip := 0
	if len(args) >= 3 {
		page, err := strconv.Atoi(args[1])
		if err != nil {
			return shim.Error(err.Error())
		}

		size, err := strconv.Atoi(args[2])
		if err != nil {
			return shim.Error(err.Error())
		}

		skip = (page - 1) * size
		limit = size
	}

	query := "{\"selector\":{\"docType\":\"" + DocType + "\"" + parametersString + "}, \"limit\": " + strconv.Itoa(limit) + ",\"skip\":" + strconv.Itoa(skip) + "}"

	// if len(args) >= 2 {
	// 	pageSize, err := strconv.ParseInt(args[1], 10, 32)
	// 	if err != nil {
	// 		return shim.Error(err.Error())
	// 	}

	// 	bookmark := ""
	// 	if len(args) >= 3 {
	// 		bookmark = args[2]
	// 	}

	// 	return s.queryDataWithBookmark(stub, query, int32(pageSize), bookmark)
	// }

	return s.queryData(stub, query)
}

func (s *SmartContract) queryData(stub shim.ChaincodeStubInterface, query string) sc.Response {

	fmt.Printf(" query:%s\n", query)

	resultsIterator, err := stub.GetQueryResult(query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- queryData:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

// To be include only in version 1.4
// func (s *SmartContract) queryDataWithBookmark(stub shim.ChaincodeStubInterface, query string, pageSize int32, bookmark string) sc.Response {

// 	fmt.Printf(" query:%s\n", query)
// 	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(query, pageSize, bookmark)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	defer resultsIterator.Close()

// 	buffer, err := constructQueryResponseFromIterator(resultsIterator)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)

// 	fmt.Printf("- queryData:\n%s\n", bufferWithPaginationInfo.String())
// 	return shim.Success(bufferWithPaginationInfo.Bytes())
// }

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	return &buffer, nil
}

// To be include only in version 1.4
// func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *sc.QueryResponseMetadata) *bytes.Buffer {

// 	buffer.WriteString("[{\"ResponseMetadata\":{\"RecordsCount\":")
// 	buffer.WriteString("\"")
// 	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
// 	buffer.WriteString("\"")
// 	buffer.WriteString(", \"Bookmark\":")
// 	buffer.WriteString("\"")
// 	buffer.WriteString(responseMetadata.Bookmark)
// 	buffer.WriteString("\"}}]")

// 	return buffer
// }
