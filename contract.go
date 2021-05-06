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

type Asset struct {
	Name string `json:"name"`
	ID string `json:"name"`
	MetaData string `json:"metadata"`
}

type AssetHistory struct {
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
	//将对象进行序列化
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

// 资产登记
func assetEnroll(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 4 {
		return shim.Error("Error number of args")
	}

	assetName := args[0]
	assetId := args[1]
	metaData := args[2]
	ownerId := args[3]

	if assetName == "" || assetId == "" || ownerId == "" {
		return shim.Error("Invalid args")
	}

	userBytes, err := stub.GetState(constructUserKey(ownerId))
	if err == nil || len(userBytes) == 0 {
		return shim.Error("User not found")
	}

	if assetBytes, err := stub.GetState(constructAssetKey(assetId)); 
		err == nil && len(assetBytes) != 0 {
			return shim.Error("Assert already exist")
		}

	asset := &Asset{
		Name: assetName,
		ID: assetId,
		MetaData: metaData,
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal asset error: %s", err))
	}
	if err := stub.PutState(constructAssetKey(assetId), assetBytes); err != nil {
		return shim.Error(fmt.Sprintf("save asset error: %s", err))
	}
	
	//更新用户信息
	user := new(User)
	if err := json.Unmarshal(userBytes, user); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error: %s", err))
	}
	user.Assets = append(user.Assets, assetId)
	//序列化user
	userBytes, err = json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshall user error: %s", err))
	}
	if err := stub.PutState(constructUserKey(user.ID), userBytes); err != nil {
		return shim.Error(fmt.Sprintf("update user error: %s", err))
	}

	history := &AssetHistory{
		AssetID: assetId,
		OriginOwnerID: originOwner,
		CurrentOwnerID: ownerId,
	}
	historyBytes, err := json.Marshal(history)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal asset history error: %s", err))
	}
	
	historyKey, err := stub.CreateCompositeKey("history", []string{
		assetId,
		originOwner,
		ownerId,
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create key error: %s", err))
	}
	if err := stub.PutState(historyKey, historyBytes); err != nil {
		return shim.Error(fmt.Sprintf("save asset history error: %s", err))
	}
	return shim.Success(historyBytes)
}