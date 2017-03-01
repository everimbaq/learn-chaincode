package main

import (

	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


type HelloWorldChaincode struct {
}

func (t *HelloWorldChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("HelloWorld - Init called with function %s!\n", function)

	return nil, nil
}

func (t *HelloWorldChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("HelloWorld - Invoke called with function %s!\n", function)
	message := "Hello World"
	return []byte(message), nil
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