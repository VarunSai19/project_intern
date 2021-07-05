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
}

type HistoryQueryResult struct {
	Record    *TelcoData    `json:"record"`
	TxId     string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
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
	
	// err1 := ctx.GetClientIdentity().AssertAttributeValue("usertype", "customer")
	// if err1 != nil {
	// 	return "",fmt.Errorf("submitting client not authorized to create asset.")
	// }

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


func (s *SmartContract) GetDataByPhoneNumber(ctx contractapi.TransactionContextInterface, ID string) (*TelcoData, error) {
	if len(ID) == 0 {
		return nil, fmt.Errorf("Please provide correct contract Id")
	}

	// err := ctx.GetClientIdentity().AssertAttributeValue("usertype", "telco-admin")
	clientID,err_id := s.GetSubmittingClientIdentity(ctx)

	if err_id != nil{
		return nil,err_id
	}


	// if err != nil {
	// 	return "",fmt.Errorf("submitting client not authorized to create asset.")
	// }

	fmt.Print("Id of the user is ",clientID,"\n")

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

func (s *SmartContract) DeleteDataByUserName(ctx contractapi.TransactionContextInterface, ID string) (string, error) {
	if len(ID) == 0 {
		return "", fmt.Errorf("Please provide correct contract Id")
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().DelState(ID)
}

func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
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
