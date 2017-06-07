/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type Mandate struct {
    Name string `json:"name"`
    Bank  string `json:"bank"`
	DateOfBirth string `json:"dateOfBirth"`
}
var SampleMandates = []Mandate{
    {
        Name: "Adil Haris",
        Bank: "HDFC Bank",
		DateOfBirth: "11th July 1993",
    },{
        Name: "John Johny Johnson",
        Bank: "ICICI Bank",
		DateOfBirth: "12th July 1993",
    },
}

var mandateCount int
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	err := stub.PutState("Init", []byte(args[0]))
	if err != nil {
		return nil, err
	}
	mandateCount = 2
	mandateCountString := strconv.Itoa(mandateCount)
	err1 := stub.PutState("MandateCount", []byte(mandateCountString))
	if err1 != nil {
		return nil, err1
	}
	MandatesBytes, _ := json.Marshal(&SampleMandates)
	err2 := stub.PutState("SampleMandates", MandatesBytes)
	if err2 != nil {
		return nil, err2
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string)([]byte, error) {
	fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "newMandate" {
		return t.newMandate(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) newMandate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var newmandate Mandate
	fmt.Println("running write()")
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name, bank and date of birth of the investor are required")
	}
	newmandate.Name = args[0]
	newmandate.Bank = args[1]
	newmandate.DateOfBirth = args[2]
	mandateCount++
	mandateCountString := strconv.Itoa(mandateCount)
	err := stub.PutState("MandateCount", []byte(mandateCountString))
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ := json.Marshal(newmandate)
	err1 := stub.PutState(mandateCountString, jsonAsBytes) //write the variable into the chaincode state
	if err1 != nil {
		return nil, err1
	}
	SampleMandates = append(SampleMandates,newmandate)
	MandatesBytes, _ := json.Marshal(&SampleMandates)
	err2 := stub.PutState("SampleMandates", MandatesBytes)
	if err2 != nil {
		return nil, err2
	}
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string)([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)
	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, nil
}

func addNewMandateEntry(name string,bank string,DoB string) {
	newMandate := Mandate{
		Name:    name,
		Bank: bank,
		DateOfBirth: DoB,
	}
	SampleMandates = append(SampleMandates, newMandate)
}
