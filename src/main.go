package main

import (
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type privateWorld struct {
}

func main() {
	shim.Start(new(privateWorld))
}

// Init is called during Instantiate transaction
func (cc *privateWorld) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction
func (cc *privateWorld) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "readCollection":
		transient, err := stub.GetTransient()
		if err != nil {
			return shim.Error(err.Error())
		}
		return cc.readCollection(stub, transient)
	case "writeCollection":
		transient, err := stub.GetTransient()
		if err != nil {
			return shim.Error(err.Error())
		}
		return cc.writeCollection(stub, transient)
	case "write":
		return cc.write(stub, args)
	case "read":
		return cc.read(stub, args)
	default:
		return shim.Error("Valid methods are 'writeCollection|readCollection|write|read'!")
	}
}

// Read text by ID in private collection
func (cc *privateWorld) readCollection(stub shim.ChaincodeStubInterface, transient map[string][]byte) peer.Response {

	col := string(transient["collection"])
	id := strings.ToLower(string(transient["id"]))

	if value, err := stub.GetPrivateData(col, id); err == nil {
		if value != nil {
			return shim.Success(value)
		}
		return shim.Error("Not Found")
	} else {
		return shim.Error(err.Error())
	}
}

// Write text into private collection
func (cc *privateWorld) writeCollection(stub shim.ChaincodeStubInterface, transient map[string][]byte) peer.Response {

	col := string(transient["collection"])
	id := strings.ToLower(string(transient["id"]))
	txt := string(transient["value"])

	if err := stub.PutPrivateData(col, id, []byte(txt)); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Read data from the ledger by ID
func (cc *privateWorld) read(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Parameter Mismatch")
	}
	id := strings.ToLower(args[0])

	if value, err := stub.GetState(id); err == nil {
		if value != nil {
			return shim.Success(value)
		}
		return shim.Error("Not Found")
	} else {
		return shim.Error(err.Error())
	}
}

// Write data into the ledger
func (cc *privateWorld) write(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 || len(args[0]) == 0 || len(args[1]) == 0 {
		return shim.Error("Parameter Mismatch")
	}
	id := strings.ToLower(args[0])
	txt := args[1]

	if err := stub.PutState(id, []byte(txt)); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
