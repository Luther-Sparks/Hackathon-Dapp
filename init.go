package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type SimpleContract struct {
}

func (t * SimpleContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments.Expecting a key and a value")
	}
	_, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for state")
	}

	err = stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create state: %s", args[0]))
	}
	return shim.Success(nil)
}

func (t *SimpleContract) get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}
	var key string
	var err error
	key = args[0]
	valBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get state: %s", key))
	}
	if valBytes == nil {
		return shim.Error(fmt.Sprintf("Nil amount: %s", key))
	}
	return shim.Success(valBytes)
}

func (t *SimpleContract) set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")	
	}
	_, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for state")
	}
	err = stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to set state: %s", args[0]))
	}
	return shim.Success(nil)
}

func (t *SimpleContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "get" {
		return t.get(stub, args)
	} else if function == "set" {
		return t.set(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"get\" \"set\"")
}

func main() {
	if err := shim.Start(new(SimpleContract)); err != nil {
		fmt.Printf("Error starting SimpleContract chaincode: %s", err)
	}
}