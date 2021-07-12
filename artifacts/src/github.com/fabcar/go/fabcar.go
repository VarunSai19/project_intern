package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"encoding/base64"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

type SmartContract struct {
	contractapi.Contract
}

var logger = flogging.MustGetLogger("fabcar_cc")

type TelcoData struct{
	AadharNumber string `json:"AadharNumber"`
	Name   string `json:"Name"`
	PhoneNumber  string `json:"PhoneNumber"`
	Status   string `json:"Status"`
	Money int64 `json:"Money"`
	Doc_type string `json:"Doc_type"`
}

type ServiceData struct{
	ServiceName   string `json:"ServiceName"`
	UserName  string `json:"UserName"`
	Doc_type string `json:"Doc_type"`
}

type AadharData struct {
	AadharNumber   string `json:"AadharNumber"`
	Address    string `json:"Address"`
	DateOfBirth   string `json:"DateOfBirth"`
	Name   string `json:"Name"`
	Gender   string `json:"Gender"`
}

type DrivingLicence struct {
	LicenceNumber  string `json:"LicenceNumber"`
	Address    string `json:"Address"`
	DateOfBirth   string `json:"DateOfBirth"`
	Name   string `json:"Name"`
	Gender   string `json:"Gender"`
	LicenceValidity   string `json:"LicenceValidity"`
}

type User_Data struct {
	UserName   string `json:"UserName"`
	AadharNumber   string `json:"AadharNumber"`
	LicenceNumber  string `json:"LicenceNumber"`
	Address    string `json:"Address"`
	DateOfBirth   string `json:"DateOfBirth"`
	Name   string `json:"Name"`
	Gender   string `json:"Gender"`
	LicenceValidity   string `json:"LicenceValidity"`
}

type Car struct {
	ID      string `json:"id"`
	Make    string `json:"make"`
	Model   string `json:"model"`
	Color   string `json:"color"`
	Owner   string `json:"owner"`
	AddedAt uint64 `json:"addedAt"`	
}

func (s *SmartContract) CreateData(ctx contractapi.TransactionContextInterface, Data string) (string, error) {
	if len(Data) == 0 {
		return "", fmt.Errorf("Please pass the correct data")
	}
	var data TelcoData
	err := json.Unmarshal([]byte(Data), &data)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling Data. %s", err.Error())
	}

	dataAsBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling Data. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", dataAsBytes)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(data.PhoneNumber, dataAsBytes)
}

func (s *SmartContract) ChangeData(ctx contractapi.TransactionContextInterface, Data string) error {
	if len(Data) == 0 {
		return "", fmt.Errorf("Please pass the correct data")
	}

	var newdata TelcoData
	err := json.Unmarshal([]byte(Data), &data)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling Data. %s", err.Error())
	}

	data,err := s.ReadAsset(ctx,newdata.PhoneNumber)

	data.AadharNumber = newdata.AadharNumber
	data.Name = newdata.Name;

	dataAsBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling Data. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", dataAsBytes)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(data.PhoneNumber, dataAsBytes)
}

func (s *SmartContract) AddMoney(ctx contractapi.TransactionContextInterface, Id string,amount int64) (string, error) {
	if len(Id) == 0 {
		return "", fmt.Errorf("Please pass the correct data")
	}
	asset,err := s.ReadAsset(ctx,Id);

	asset.Money = asset.Money + amount
	asset.Status = "Active"

	dataAsBytes, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling Data. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", dataAsBytes)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(asset.PhoneNumber, dataAsBytes)
}

func (s *SmartContract) BuyService(ctx contractapi.TransactionContextInterface, username string,servicename string,price string) (string, error) {
	if len(username) == 0 {
		return "", fmt.Errorf("Please pass the correct data.")
	}
	Price, err := strconv.ParseInt(price, 10, 64)
	if err!=nil{
		return "", fmt.Errorf("Please price of product correctly.")
	}

	asset,err := s.ReadAsset(ctx,username);
	if err!=nil{
		return "", fmt.Errorf("Problem while reading the data.")
	}
	if asset.Money < Price {
		return "", fmt.Errorf("Insufficient amount in wallet.")
	}

	data := &ServiceData{
		ServiceName:servicename,
		UserName: username+"_service",
		Doc_type: "service",
	}
	
	dataAsBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling Data. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", dataAsBytes)

	err = ctx.GetStub().PutState(data.UserName, dataAsBytes)

	if err != nil {
		return "",fmt.Errorf("Failed while pushing the transaction.")
	}

	asset.Money = asset.Money - Price;
	asset.Status = "Active"

	dataAsBytes1, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling Data. %s", err.Error())
	}
	
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(asset.PhoneNumber, dataAsBytes1)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, ID string) (*TelcoData, error) {
	if len(ID) == 0 {
		return nil, fmt.Errorf("Please provide correct contract Id")
	}
	dataAsBytes, err := ctx.GetStub().GetState(ID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if dataAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", ID)
	}
	data := new(TelcoData)
	_ = json.Unmarshal(dataAsBytes, data)

	return data, nil
}

func (s *SmartContract) GetDataByPhoneNumber(ctx contractapi.TransactionContextInterface, ID string) (*TelcoData, error) {
	if len(ID) == 0 {
		return nil, fmt.Errorf("Please provide correct contract Id")
	}
	// clientID,err_id := s.GetSubmittingClientIdentity(ctx)
	dataAsBytes, err := ctx.GetStub().GetState(ID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if dataAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", ID)
	}
	data := new(TelcoData)
	_ = json.Unmarshal(dataAsBytes, data)

	return data, nil
}

func (s *SmartContract) GetHistoryForAsset(ctx contractapi.TransactionContextInterface, ID string) (string, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(ID)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return string(buffer.Bytes()), nil
}

func (s *SmartContract) DeleteDataById(ctx contractapi.TransactionContextInterface, ID string) (string, error) {
	if len(ID) == 0 {
		return "", fmt.Errorf("Please provide correct contract Id")
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().DelState(ID)
}

func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	// x509::CN=telco-admin,OU=o 
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	res := string(decodeID)
	i:=0
	id:=""
	for ;i<len(res);i++{
		if res[i] == '='{
			break	
		}
	}
	for i=i+1;i<len(res);i++{
		if res[i] == ','{
			break	
		} 
		id += string(res[i])
	} 
	return id, nil
}


func (s *SmartContract) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]TelcoData, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []TelcoData{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newData := new(TelcoData)
		fmt.Print("Responce is ",response.Value,"\n")
		err = json.Unmarshal(response.Value, newData)
		if err == nil {
			results = append(results, *newData)
		}
	}
	return results, nil
}

func (s *SmartContract) QueryAllData(ctx contractapi.TransactionContextInterface, queryString string) ([]TelcoData, error) {
	err := ctx.GetClientIdentity().AssertAttributeValue("usertype", "telco-admin")
	if err != nil {
		return nil,fmt.Errorf("submitting client not authorized to perform this task.")
	}

	return s.getQueryResultForQueryString(ctx,queryString)
}


func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carData string) (string, error) {

	if len(carData) == 0 {
		return "", fmt.Errorf("Please pass the correct car data")
	}

	var car Car
	err := json.Unmarshal([]byte(carData), &car)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling car. %s", err.Error())
	}

	carAsBytes, err := json.Marshal(car)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling car. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", carAsBytes)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(car.ID, carAsBytes)
}


//
func (s *SmartContract) UpdateCarOwner(ctx contractapi.TransactionContextInterface, carID string, newOwner string) (string, error) {

	if len(carID) == 0 {
		return "", fmt.Errorf("Please pass the correct car id")
	}

	carAsBytes, err := ctx.GetStub().GetState(carID)

	if err != nil {
		return "", fmt.Errorf("Failed to get car data. %s", err.Error())
	}

	if carAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", carID)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	car.Owner = newOwner

	carAsBytes, err = json.Marshal(car)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling car. %s", err.Error())
	}

	//  txId := ctx.GetStub().GetTxID()

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(car.ID, carAsBytes)

}


func (s *SmartContract) GetCarById(ctx contractapi.TransactionContextInterface, carID string) (*Car, error) {
	if len(carID) == 0 {
		return nil, fmt.Errorf("Please provide correct contract Id")
		// return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, err := ctx.GetStub().GetState(carID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carID)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil

}

func (s *SmartContract) DeleteCarById(ctx contractapi.TransactionContextInterface, carID string) (string, error) {
	if len(carID) == 0 {
		return "", fmt.Errorf("Please provide correct contract Id")
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().DelState(carID)
}



func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincodes: %s", err.Error())
	}

}
