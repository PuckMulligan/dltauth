package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// AuthChaincode example simple Chaincode implementation
type AuthChaincode struct {
}

type authRecord struct {
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
	Success   bool   `json:"success"`
}

// Init is called during chaincode instantiation to initialize any data.
func (t *AuthChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (t *AuthChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "recordAuth" {
		return t.recordAuth(stub, args)
	} else if function == "queryAuths" {
		return t.queryAuths(stub, args)
	}

	return shim.Error("Invalid function name.")
}

// recordAuth records an authentication attempt
func (t *AuthChaincode) recordAuth(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	userID := args[0]
	success, err := strconv.ParseBool(args[1])
	if err != nil {
		return shim.Error("Failed to parse 'success' flag")
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	auth := authRecord{UserID: userID, Timestamp: timestamp, Success: success}

	authJSONasBytes, err := json.Marshal(auth)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState("AUTH"+timestamp+userID, authJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// queryAuths returns all authentication records
func (t *AuthChaincode) queryAuths(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	startKey := "AUTH"
	endKey := "AUTH~"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(AuthChaincode))
	if err != nil {
		fmt.Printf("Error starting Auth chaincode: %s", err)
	}
}

// put blockchain authorization through a different hashing func that is *not* stored or *only* stored in some encrypted way by the dlt. the dlt will ensure integrity and confidentiality when the original
// has function cannot
