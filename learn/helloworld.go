package main

import (

	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
)


type HelloWorldChaincode struct {
}

func (t *HelloWorldChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("HelloWorld - Init called with function %s!\n", function)

	return nil, nil
}

func (t *HelloWorldChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *HelloWorldChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *HelloWorldChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("HelloWorld - Query called with function %s!\n", function)

	message := "Hello World"
	return []byte(message), nil;
}

func main() {
	err := shim.Start(new(HelloWorldChaincode))
	if err != nil {
		fmt.Printf("Error starting Hello World chaincode: %s", err)
	}
}