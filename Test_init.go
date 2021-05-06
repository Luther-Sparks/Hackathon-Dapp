package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)


func checkInit(t *testing.T, stub *shimtest.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shimtest.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkGet(t *testing.T, stub *shimtest.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("get"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("get", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("get", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("get value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkSet(t *testing.T, stub *shimtest.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func Test_Init(t *testing.T) {
	cc := new(SimpleContract)
	stub := shimtest.NewMockStub("sccc", cc)
	checkInit(t, stub, [][]byte{[]byte("a"), []byte("10")})
	checkGet(t, stub, "a", "10")
}

func Test_Set(t *testing.T) {
	cc := new(SimpleContract)
	stub := shimtest.NewMockStub("sccc", cc)
	checkInit(t, stub, [][]byte{[]byte("a"), []byte("10")})
	checkSet(t, stub, [][]byte{[]byte("set"), []byte("a"), []byte("20")})
	checkGet(t, stub, "a", "20")
}