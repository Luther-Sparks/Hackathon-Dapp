package main

import (
	"fmt"
	"encodin/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type User struct {
	Name string `json:"name"`
	ID string `json:"id"`
	Asserts []string `json:"asserts"`
}

type Assert struct {
	Name string `json:"name"`
	ID string `json:"name"`
	MetaData string `json:"metadata"`
}

type AssertHistory struct {
	AssertID string `json:"assert_id"`
	OriginOwnerID string `json:"origin_owner_id"`
	CurrentOwnerID string `json:"current_owner_id"`
}

const (
	
)