/*
Copyright IBM Corp. 2019 All Rights Reserved.

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
	"fmt"
	"errors"
	"strconv"
	"bytes"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


type License struct {
	Title string `json:"title"`
	Artist string `json:"artist"`
	ISRC string `json:"ISRC"`
	ISWC string `json:"ISWC"`
	ReleaseDate string `json:"releaseDate"`
	UseType string `json:"useType"`
	RateAdvance int `json:"rateAdvance"`
	MasterArtist1 string `json:"masterArtist1"`
	MasterPercentage1 int `json:"masterPercentage1"`
	PublishingArtist1 string `json:"publishingArtist1"`
	PublishingPercentage1 int `json:"publishingPercentage1"`
	MasterArtist2 string `json:"masterArtist2"`
	MasterPercentage2 int `json:"masterPercentage2"`
	PublishingArtist2 string `json:"publishingArtist2"`
	PublishingPercentage2 int `json:"publishingPercentage2"`
	LimitDistribution bool `json:"limitDistribution"`
	LimitDistributionEarliest string `json:"limitDistributionEarliest"`
	LimitDistributionLatest string `json:"limitDistributionLatest"`
}

type TransferOfRights struct {
	title string `json:"title"`
	licensee string `json:"licensee"`
}

const (
	NETWORK_USER_ID										= "network-user"
	NEWTORK_FEE_PERCENTAGE						= 10
	NETWORK_WALLET_DEFAULT_BALANCE		= 0
	WALLET_KEY_PREFIX									= "wallet--"
	LICENCE_KEY_PREFIX								= "license--"
	TRANSFER_OF_RIGHTS_KEY_PREFIX			= "transferOfRights--"
)




// DRMChaincode implementation
type DRMChaincode struct {

}

func (t *DRMChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	var err error


	fmt.Println("Initializing DRM")

	// create network wallet with default balance
	err = t.setWalletBalance(stub, NETWORK_USER_ID, NETWORK_WALLET_DEFAULT_BALANCE)
	if err != nil { return shim.Error(err.Error()) }

	return shim.Success(nil)
}


func (t *DRMChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	fmt.Println("DRM Invoke: ", function)

	if function == "createWallet" {
		return t.createWallet(stub, args)
	} else if function == "createLicense" {
		return t.createLicense(stub, args)
	} else if function == "searchForLicense" {
		return t.searchForLicense(stub, args)
	} else if function == "transferRights" {
		return t.transferRights(stub, args)
	}

	return shim.Error("Invalid invoke function name")
}

func (t *DRMChaincode) getLicenseKey(stub shim.ChaincodeStubInterface, title string) (string, error) {
	licenseKey, err := stub.CreateCompositeKey(LICENCE_KEY_PREFIX, []string{title})
	if err != nil {
		return "", err
	} else {
		return licenseKey, nil
	}
}

func (t *DRMChaincode) getWalletKey(stub shim.ChaincodeStubInterface, userId string) (string, error) {
	walletKey, err := stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{userId})
	if err != nil {
		return "", err
	} else {
		return walletKey, nil
	}
}

func (t *DRMChaincode) getTransferOfRightsKey(stub shim.ChaincodeStubInterface, title string, licenseeUserId string) (string, error) {
	transferOfRightsKey, err := stub.CreateCompositeKey(TRANSFER_OF_RIGHTS_KEY_PREFIX, []string{title, licenseeUserId})
	if err != nil {
		return "", err
	} else {
		return transferOfRightsKey, nil
	}
}


func (t *DRMChaincode) setWalletBalance(stub shim.ChaincodeStubInterface, userId string, balance int) error {

	var err error
	var walletKey string


	walletKey, err = t.getWalletKey(stub, userId)
	if err != nil {
		return err
	}

	err = stub.PutState(walletKey, []byte(strconv.Itoa(balance)))
	if err != nil {
		return err
	}

	fmt.Printf("Wallet %s balance %d recorded\n", walletKey, balance)

	return nil
}

func (t *DRMChaincode) getWalletBalance(stub shim.ChaincodeStubInterface, userId string) (int, error) {

	var err error
	var balanceBytes []byte
	var balance int
	var walletKey string


	walletKey, err = t.getWalletKey(stub, userId)
	if err != nil {
		return 0, err
	}

	balanceBytes, err = stub.GetState(walletKey)
	if err != nil {
		return 0, errors.New("{\"Error\":\"Failed to get state for " + walletKey + "\"}")
	}

	if len(balanceBytes) == 0 {
		return 0, errors.New("{\"Error\":\"No record found for " + walletKey + "\"}")
	}

	balance, err = strconv.Atoi(string(balanceBytes))
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (t *DRMChaincode) createWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var userId string
	var defaultBalance int
	var err error


	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2 {userId, defaultBalance} - found %d", len(args)))
		return shim.Error(err.Error())
	}

	userId = string(args[0])

	defaultBalance, err = strconv.Atoi(string(args[1]))
	if err != nil { return shim.Error(err.Error()) }

	err = t.setWalletBalance(stub, userId, defaultBalance)
	if err != nil { return shim.Error(err.Error()) }

	return shim.Success(nil)
}

func (t *DRMChaincode) createLicense(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var licenseKey string
	var license *License
	var licenseBytes []byte
	var err error
	var rateAdvance int
	var masterPercentage1 int
	var publishingPercentage1 int
	var masterPercentage2 int
	var publishingPercentage2 int
	var limitDistribution bool


	if len(args) != 18 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 18 - found %d", len(args)))
		return shim.Error(err.Error())
	}

	rateAdvance, err = strconv.Atoi(string(args[6]))
	if err != nil { return shim.Error(err.Error()) }

	masterPercentage1, err = strconv.Atoi(string(args[8]))
	if err != nil { return shim.Error(err.Error()) }

	publishingPercentage1, err = strconv.Atoi(string(args[10]))
	if err != nil { return shim.Error(err.Error()) }

	masterPercentage2, err = strconv.Atoi(string(args[12]))
	if err != nil { return shim.Error(err.Error()) }

	publishingPercentage2, err = strconv.Atoi(string(args[14]))
	if err != nil { return shim.Error(err.Error()) }

	limitDistribution, err = strconv.ParseBool(string(args[15]))
	if err != nil { return shim.Error(err.Error()) }

	license = &License{
		args[0],
		args[1],
		args[2],
		args[3],
		args[4],
		args[5],
		rateAdvance,
		args[7],
		masterPercentage1,
		args[9],
		publishingPercentage1,
		args[11],
		masterPercentage2,
		args[13],
		publishingPercentage2,
		limitDistribution,
		args[16],
		args[17],
	}

	licenseBytes, err = json.Marshal(license)
	if err != nil {
		return shim.Error("Error marshaling License structure")
	}


	// Write the license to the ledger
	licenseKey, err = t.getLicenseKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(licenseKey, licenseBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("License %s recorded\n", args[0])

	return shim.Success(nil)
}

func (t *DRMChaincode) searchForLicense(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var useType string


	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1 {useType} - found %d", len(args)))
		return shim.Error(err.Error())
	}

	useType = string(args[0])

	queryString :=
	`{
		"selector": {
				"useType": ` + useType + `
		}
	}`

	fmt.Printf("queryString:\n%s\n", queryString)

  // Invoke query
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	// Iterate through all returned assets
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("queryResult:\n%s\n", buffer.String())

	return shim.Success([]byte(buffer.String()))
}

// transfer of rights - calculates "advance" fees only
func (t *DRMChaincode) transferRights(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var title string
	var licenseKey string
	var license *License
	var licenseBytes []byte
	var transferOfRightsKey string
	var transferOfRights *TransferOfRights
	var transferOfRightsBytes []byte
	var err error
	var licenseeUserId string
	var networkFees int
	var licensorMaster1Fees int = 0
	var licensorMaster2Fees int = 0
	var licensorPublisher1Fees int = 0
	var licensorPublisher2Fees int = 0
	var totalFees int
	var networkBalance int
	var licenseeBalance int
	var licensorMaster1Balance int
	var licensorMaster2Balance int
	var licensorPublisher1Balance int
	var licensorPublisher2Balance int


	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2 {title, licenseeUserId} - found %d", len(args)))
		return shim.Error(err.Error())
	}

	title = args[0]
	licenseeUserId = args[1]


	// get the license from the ledger
	licenseKey, err = t.getLicenseKey(stub, title)
	if err != nil {	return shim.Error(err.Error()) }

	licenseBytes, err = stub.GetState(licenseKey)
	if err != nil { return shim.Error(err.Error()) }

	if len(licenseBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for title %s", title))
		return shim.Error(err.Error())
	}

	err = json.Unmarshal(licenseBytes, &license)
	if err != nil {	return shim.Error(err.Error()) }


	// get current balances of wallets
	networkBalance, err = t.getWalletBalance(stub, NETWORK_USER_ID)
	if err != nil { return shim.Error(err.Error()) }
	licenseeBalance, err = t.getWalletBalance(stub, licenseeUserId)
	if err != nil { return shim.Error(err.Error()) }
	licensorMaster1Balance, err = t.getWalletBalance(stub, license.MasterArtist1)
	if err != nil { return shim.Error(err.Error()) }
	licensorMaster2Balance, err = t.getWalletBalance(stub, license.MasterArtist2)
	if err != nil { return shim.Error(err.Error()) }
	licensorPublisher1Balance, err = t.getWalletBalance(stub, license.PublishingArtist1)
	if err != nil { return shim.Error(err.Error()) }
	licensorPublisher2Balance, err = t.getWalletBalance(stub, license.PublishingArtist2)
	if err != nil { return shim.Error(err.Error()) }


	// calculate fees
	networkFees = (license.RateAdvance * NEWTORK_FEE_PERCENTAGE) / 100
	if (license.MasterPercentage1 > 0) {
		licensorMaster1Fees = (license.RateAdvance * license.MasterPercentage1) / 100
	}
	if (license.MasterPercentage2 > 0) {
		licensorMaster2Fees = (license.RateAdvance * license.MasterPercentage2) / 100
	}
	if (license.PublishingPercentage1 > 0) {
		licensorPublisher1Fees = (license.RateAdvance * license.PublishingPercentage1) / 100
	}
	if (license.PublishingPercentage2 > 0) {
		licensorPublisher2Fees = (license.RateAdvance * license.PublishingPercentage2) / 100
	}

	totalFees = networkFees + licensorMaster1Fees + licensorMaster2Fees + licensorPublisher1Fees + licensorPublisher2Fees

	if (licenseeBalance < totalFees){
		err = errors.New(fmt.Sprintf("Insufficient funds in wallet - userid %s", licenseeUserId))
		return shim.Error(err.Error())
	}


	// calculate new balances
	networkBalance += networkFees
	licenseeBalance -= totalFees
	licensorMaster1Balance += licensorMaster1Fees
	licensorMaster2Balance += licensorMaster2Fees
	licensorPublisher1Balance += licensorPublisher1Fees
	licensorPublisher2Balance += licensorPublisher2Fees


	// record new balances
	err = t.setWalletBalance(stub, NETWORK_USER_ID, networkBalance)
	if err != nil { return shim.Error(err.Error()) }
	err = t.setWalletBalance(stub, licenseeUserId, licenseeBalance)
	if err != nil { return shim.Error(err.Error()) }
	err = t.setWalletBalance(stub, license.MasterArtist1, licensorMaster1Balance)
	if err != nil { return shim.Error(err.Error()) }
	err = t.setWalletBalance(stub, license.MasterArtist2, licensorMaster2Balance)
	if err != nil { return shim.Error(err.Error()) }
	err = t.setWalletBalance(stub, license.PublishingArtist1, licensorPublisher1Balance)
	if err != nil { return shim.Error(err.Error()) }
	err = t.setWalletBalance(stub, license.PublishingArtist2, licensorPublisher2Balance)
	if err != nil { return shim.Error(err.Error()) }


	// record transfer of rights
	transferOfRights = &TransferOfRights{
		title,
		licenseeUserId,
	}

	transferOfRightsBytes, err = json.Marshal(transferOfRights)
	if err != nil {
		return shim.Error("Error marshaling TransferOfRights structure")
	}

	transferOfRightsKey, err = t.getTransferOfRightsKey(stub, title, licenseeUserId)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(transferOfRightsKey, transferOfRightsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	fmt.Printf("Transfer of rights of %s to %s completed\n", title, licenseeUserId)

	return shim.Success(nil)
}



func main() {
	drmc := new(DRMChaincode)
	err := shim.Start(drmc)
	if err != nil {
		fmt.Printf("Error starting DRM chaincode: %s", err)
	}
}
