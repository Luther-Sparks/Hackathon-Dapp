package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	// "github.com/hyperledger/fabric/core/chaincode/shim"
	// "github.com/hyperledger/fabric/protos/peer"
	// "github.com/stretchr/testify/require"
	// "golang.org/x/text/cases"
	// "golang.org/x/tools/go/analysis/passes/nilfunc"
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

// 实现资产转让
func assetExchange(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Error number of args")
	}

	ownerId := args[0]
	assetId := args[1]
	currentOwnerId := args[2]

	if ownerId == "" || assetId == "" || currentOwnerId == "" {
		return shim.Error("Invalid args")
	}
	// 验证数据是否存在
	originOwnerBytes, err := stub.GetState(constructUserKey(ownerId))
	if err != nil || len(originOwnerBytes) == 0 {
		return shim.Error("user not found")
	}
	currentOwnerBytes, err := stub.GetState(constructUserKey(currentOwnerId))
	if err != nil || len(currentOwnerBytes) == 0 {
		return shim.Error("user not found")
	}
	assetBytes, err := stub.GetState(constructAssetKey(assetId))
	if err != nil || len(assetBytes) == 0 {
		return shim.Error("asset not found")
	}

	// 检验卖家是否存在指定资产
	originOwner := new(User)
	// 反序列化user
	if err := json.Unmarshal(originOwnerBytes, originOwner); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error: %s", err))
	}
	assetid_exist := false
	for _, assetid := range originOwner.Assets {
		if assetid == assetId {
			assetid_exist = true
			break
		}
	}
	if !assetid_exist{
		return shim.Error("asset owner doesn't have the aim asset")
	}

	// 将结果写入
	totolAssetId := make([]string, 0)
	for _, assetid := range originOwner.Assets {
		if assetid == assetId {
			continue
		}
		totolAssetId = append(totolAssetId, assetid)
	}
	originOwner.Assets = totolAssetId

	originOwnerBytes, err = json.Marshal(originOwner)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error: %s", err))
	}
	if err := stub.PutState(constructUserKey(ownerId), originOwnerBytes); err != nil {
		return shim.Error(fmt.Sprintf("update user error: %s", err))
	}

	//买家添加资产id
	currentOwner := new(User)
	//反序列化user
	if err := json.Unmarshal(currentOwnerBytes, currentOwner); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error: %s", err))
	}
	currentOwner.Assets = append(currentOwner.Assets, assetId)

	currentOwnerBytes, err = json.Marshal(currentOwner)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error: %s", err))
	}
	if err := stub.PutState(constructUserKey(currentOwnerId), currentOwnerBytes); err != nil {
		return shim.Error(fmt.Sprintf("update user error: %s", err))
	}

	//插入资产变更记录
	history := &AssetHistory{
		AssetID: assetId,
		OriginOwnerID: ownerId,
		CurrentOwnerID: currentOwnerId,
	}
	historyBytes, err := json.Marshal(history)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal asset history error: %s", err))
	}
	
	historyKey, err := stub.CreateCompositeKey("history", []string{
		assetId,
		ownerId,
		currentOwnerId,
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create key error: %s", err))
	}
	if err := stub.PutState(historyKey, historyBytes); err != nil {
		return shim.Error(fmt.Sprintf("save asset history error: %s", err))
	}

	return shim.Success(nil)
}

// 基于id查询用户
func queryUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Error number of args")
	}
	userId := args[0]
	if userId == "" {
		return shim.Error("Invalid args")
	}
	userBytes, err := stub.GetState(constructUserKey(userId))
	if err != nil || len(userBytes) == 0 {
		return shim.Error("user not found")
	}
	return shim.Success(userBytes)
}

// 基于assetId查询资产信息
func queryAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Error number of args")
	}
	assetId := args[0]
	if assetId == "" {
		return shim.Error("Invalid args")
	}
	assetBytes, err := stub.GetState(constructAssetKey(assetId))
	if err != nil || len(assetBytes) == 0 {
		return shim.Error("asset not found")
	}
	return shim.Success(assetBytes)
}

// 资产交易记录查询
func queryAssetHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 && len(args) != 1 {
		return shim.Error("Error number of args")
	}
	assetId := args[0]
	queryType := "all"
	if len(args) == 2 {
		queryType = args[1]
	}

	if queryType != "all" && queryType != "enroll" && queryType != "exchange" {
		return shim.Error(fmt.Sprintf("queryType unknown %s", queryType))
	}
	assetBytes, err := stub.GetState(constructAssetKey(assetId))
	if err != nil || len(assetBytes) == 0 {
		return shim.Error("asset not found")
	}

	keys := make([]string, 0)
	keys = append(keys, assetId)
	switch queryType {
	case "enroll":
		keys = append(keys, originOwner)
	case "exchange", "all":
	default:
		return shim.Error(fmt.Sprintf("unsupport queryType: %s", queryType))
	}
	result, err := stub.GetStateByPartialCompositeKey("history", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("query history error: %s", err))
	}
	defer result.Close()

	histories := make([]*AssetHistory, 0)
	for result.HasNext() {
		historyVal, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query error: %s", err))
		}

		history := new(AssetHistory)
		if err := json.Unmarshal(historyVal.GetValue(), history); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}

		// 过滤掉不是资产转让的记录
		if queryType == "exchange" && history.OriginOwnerID == originOwner {
			continue
		}
		
		histories = append(histories, history)
	}

	historiesBytes, err := json.Marshal(histories)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error: %s", err))
	}

	return shim.Success(historiesBytes)
}

type AssetExchangeChainCode struct{

}

func (t *AssetExchangeChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *AssetExchangeChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	functionName, args := stub.GetFunctionAndParameters()
	switch functionName {
	case "userRegister":
		return userRegister(stub, args)
	case "assetEnroll":
		return assetEnroll(stub, args)
	case "assetExchange":
		return assetExchange(stub, args)
	case "queryUser":
		return queryUser(stub, args)
	case "queryAsset":
		return queryAsset(stub, args)
	case "queryAssetHistory":
		return queryAssetHistory(stub, args)
	default:
		return shim.Error(fmt.Sprintf("unsupported function: %s", functionName))
	}
}

func main() {
	err := shim.Start(new(AssetExchangeChainCode))
	if err != nil {
		fmt.Printf("Error starting AssetExchange chaincode: %s", err)
	}
}
