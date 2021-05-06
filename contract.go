package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

type User struct {
	Name string `json:"name"`
	ID string `json:"id"`
	Assets []string `json:"assets"`
}

type asset struct {
	Name string `json:"name"`
	ID string `json:"name"`
	MetaData string `json:"metadata"`
}

type assetHistory struct {
	AssetID string `json:"asset_id"`
	OriginOwnerID string `json:"origin_owner_id"`
	CurrentOwnerID string `json:"current_owner_id"`
}

const (
	originOwner = "originOwnerPlaceHolder"
)

func constructUserKey(userId string) string {
	return fmt.Sprintf("user_%s", userId)
}

func constructAssetKey(assetId string) string {
	return fmt.Sprintf("asset_%s", assetId)
}

func userRegister(stub shim.ChaincodeStubInterface, args []string)	peer.Response {
	if len(args) != 2 {
		return shim.Error("Error number of args")
	}
	name := args[0]
	id := args[1]
	if name == "" || id == "" {
		return shim.Error("Invalid args")
	}
	if userBytes, err := stub.GetState(constructUserKey(id)); err != nil ||
		len(userBytes) != 0 {
			return shim.Error("User already exist")
		}
	user := User{
		Name: name,
		ID: id,
		Assets: make([]string, 0),
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}
	err = stub.PutState(constructUserKey(id), userBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("put user error %s", err))
	}
	return shim.Success(nil)
}

