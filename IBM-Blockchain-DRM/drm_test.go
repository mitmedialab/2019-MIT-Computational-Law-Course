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
	"testing"
	//"bytes"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	LICENSEE1 = "value-licensee1"
	LICENSEE2 = "value-licensee2"
	DEFAULT_WALLET_BALANCE = 1000
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkNoState(t *testing.T, stub *shim.MockStub, name string) {
	bytes := stub.State[name]
	if bytes != nil {
		fmt.Println("State", name, "should be absent; found value")
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was", string(bytes), "and not", value, "as expected")
		t.FailNow()
	}
}

func checkBadQuery(t *testing.T, stub *shim.MockStub, function string, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status == shim.OK {
		fmt.Println("Query", name, "unexpectedly succeeded")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		fmt.Println("Query value", name, "was", payload, "and not", value, "as expected")
		t.FailNow()
	}
}

func checkQueryArgs(t *testing.T, stub *shim.MockStub, args [][]byte, value string) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Query", string(args[1]), "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", string(args[1]), "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		fmt.Println("Query value", string(args[1]), "was", payload, "and not", value, "as expected")
		t.FailNow()
	}
}

func checkBadInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status == shim.OK {
		fmt.Println("Invoke", args, "unexpectedly succeeded")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}





func TestDRM_Init(t *testing.T) {
	fmt.Println("TestDRM_Init started")

	scc := new(DRMChaincode)
	stub := shim.NewMockStub("DRM Workflow", scc)

	checkInit(t, stub, nil)
	walletKey, _ := stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{NETWORK_USER_ID})
	checkState(t, stub, walletKey, strconv.Itoa(NETWORK_WALLET_DEFAULT_BALANCE))

	fmt.Println("TestDRM_Init passed")
}

func TestDRM_CreateWallet(t *testing.T) {
	fmt.Println("TestDRM_CreateWallet started")

	scc := new(DRMChaincode)
	stub := shim.NewMockStub("DRM Workflow", scc)

	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte(LICENSEE1), []byte(strconv.Itoa(DEFAULT_WALLET_BALANCE))})
	walletKey, _ := stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{LICENSEE1})
	checkState(t, stub, walletKey, strconv.Itoa(DEFAULT_WALLET_BALANCE))
	checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte(LICENSEE2), []byte(strconv.Itoa(DEFAULT_WALLET_BALANCE))})
	walletKey, _ = stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{LICENSEE2})
	checkState(t, stub, walletKey, strconv.Itoa(DEFAULT_WALLET_BALANCE))

	fmt.Println("TestDRM_CreateWallet passed")
}

func TestDRM_CreateLicense(t *testing.T) {
	fmt.Println("TestDRM_CreateLicense started")

	scc := new(DRMChaincode)
	stub := shim.NewMockStub("DRM Workflow", scc)

	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("createLicense"),
		[]byte("value-title"),
		[]byte("value-artist"),
		[]byte("value-ISRC"),
		[]byte("value-ISWC"),
		[]byte("value-releaseDate"),
		[]byte("value-useType"),
		[]byte("100"),
		[]byte("value-masterArtist1"),
		[]byte("10"),
		[]byte("value-publishingArtist1"),
		[]byte("20"),
		[]byte("value-masterArtist2"),
		[]byte("30"),
		[]byte("value-publishingArtist2"),
		[]byte("40"),
		[]byte("true"),
		[]byte("value-limitDistributionEarliest"),
		[]byte("value-limitDistributionLatest")},
	)

	licenseKey, _ := stub.CreateCompositeKey(LICENCE_KEY_PREFIX, []string{"value-title"})

	license := &License{
		"value-title",
		"value-artist",
		"value-ISRC",
		"value-ISWC",
		"value-releaseDate",
		"value-useType",
		100,
		"value-masterArtist1",
		10,
		"value-publishingArtist1",
		20,
		"value-masterArtist2",
		30,
		"value-publishingArtist2",
		40,
		true,
		"value-limitDistributionEarliest",
		"value-limitDistributionLatest",
	}

	licenseBytes, _ := json.Marshal(license)
	checkState(t, stub, licenseKey, string(licenseBytes))

	fmt.Println("TestDRM_CreateLicense passed")
}



/*

func TestDRM_SearchForLicense(t *testing.T) {
	fmt.Println("TestDRM_SearchForLicense started")

	scc := new(DRMChaincode)
	stub := shim.NewMockStub("DRM Workflow", scc)

	checkInit(t, stub, nil)


	// license 1
	checkInvoke(t, stub, [][]byte{[]byte("createLicense"),
		[]byte("value-title-1"),
		[]byte("value-artist"),
		[]byte("value-ISRC"),
		[]byte("value-ISWC"),
		[]byte("value-releaseDate"),
		[]byte("value-useType"),
		[]byte("100"),
		[]byte("value-masterArtist1"),
		[]byte("10"),
		[]byte("value-publishingArtist1"),
		[]byte("20"),
		[]byte("value-masterArtist2"),
		[]byte("30"),
		[]byte("value-publishingArtist2"),
		[]byte("40"),
		[]byte("true"),
		[]byte("value-limitDistributionEarliest"),
		[]byte("value-limitDistributionLatest")},
	)

	licenseKey1, _ := stub.CreateCompositeKey(LICENCE_KEY_PREFIX, []string{"value-title-1"})

	license1 := &License{
		"value-title-1",
		"value-artist",
		"value-ISRC",
		"value-ISWC",
		"value-releaseDate",
		"value-useType",
		100,
		"value-masterArtist1",
		10,
		"value-publishingArtist1",
		20,
		"value-masterArtist2",
		30,
		"value-publishingArtist2",
		40,
		true,
		"value-limitDistributionEarliest",
		"value-limitDistributionLatest",
	}

	licenseBytes1, _ := json.Marshal(license1)
	checkState(t, stub, licenseKey1, string(licenseBytes1))


	// license 2
	checkInvoke(t, stub, [][]byte{[]byte("createLicense"),
		[]byte("value-title-2"),
		[]byte("value-artist"),
		[]byte("value-ISRC"),
		[]byte("value-ISWC"),
		[]byte("value-releaseDate"),
		[]byte("value-useType"),
		[]byte("100"),
		[]byte("value-masterArtist1"),
		[]byte("10"),
		[]byte("value-publishingArtist1"),
		[]byte("20"),
		[]byte("value-masterArtist2"),
		[]byte("30"),
		[]byte("value-publishingArtist2"),
		[]byte("40"),
		[]byte("true"),
		[]byte("value-limitDistributionEarliest"),
		[]byte("value-limitDistributionLatest")},
	)

	licenseKey2, _ := stub.CreateCompositeKey(LICENCE_KEY_PREFIX, []string{"value-title-2"})

	license2 := &License{
		"value-title-2",
		"value-artist",
		"value-ISRC",
		"value-ISWC",
		"value-releaseDate",
		"value-useType",
		100,
		"value-masterArtist1",
		10,
		"value-publishingArtist1",
		20,
		"value-masterArtist2",
		30,
		"value-publishingArtist2",
		40,
		true,
		"value-limitDistributionEarliest",
		"value-limitDistributionLatest",
	}

	licenseBytes2, _ := json.Marshal(license2)
	checkState(t, stub, licenseKey2, string(licenseBytes2))



	// search output
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"")
	buffer.WriteString(licenseKey1)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Record\":")
	buffer.WriteString(string(licenseBytes1))
	buffer.WriteString("}")
	buffer.WriteString(",")
	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"")
	buffer.WriteString(licenseKey2)
	buffer.WriteString("\"")
	buffer.WriteString(", \"Record\":")
	buffer.WriteString(string(licenseBytes2))
	buffer.WriteString("}")
	buffer.WriteString("]")

	expectedResponse := buffer.String()

	// check search
	checkQuery(t, stub, "searchForLicense", "value-useType", expectedResponse)

	fmt.Println("TestDRM_SearchForLicense passed")
}
*/



func TestDRM_TransferRights(t *testing.T) {
	fmt.Println("TestDRM_TransferRights started")

	scc := new(DRMChaincode)
	stub := shim.NewMockStub("DRM Workflow", scc)

	checkInit(t, stub, nil)

	checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte(LICENSEE1), []byte(strconv.Itoa(DEFAULT_WALLET_BALANCE))})
	checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("value-masterArtist1"), []byte(strconv.Itoa(DEFAULT_WALLET_BALANCE))})
	checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("value-publishingArtist1"), []byte(strconv.Itoa(DEFAULT_WALLET_BALANCE))})
	checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("value-masterArtist2"), []byte(strconv.Itoa(DEFAULT_WALLET_BALANCE))})
	checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("value-publishingArtist2"), []byte(strconv.Itoa(DEFAULT_WALLET_BALANCE))})

	checkInvoke(t, stub, [][]byte{[]byte("createLicense"),
		[]byte("value-title"),
		[]byte("value-artist"),
		[]byte("value-ISRC"),
		[]byte("value-ISWC"),
		[]byte("value-releaseDate"),
		[]byte("value-useType"),
		[]byte("100"),
		[]byte("value-masterArtist1"),
		[]byte("10"),
		[]byte("value-publishingArtist1"),
		[]byte("20"),
		[]byte("value-masterArtist2"),
		[]byte("30"),
		[]byte("value-publishingArtist2"),
		[]byte("40"),
		[]byte("true"),
		[]byte("value-limitDistributionEarliest"),
		[]byte("value-limitDistributionLatest")},
	)

	checkInvoke(t, stub, [][]byte{[]byte("transferRights"),	[]byte("value-title"), []byte(LICENSEE1)})

	transferOfRightsKey, _ := stub.CreateCompositeKey(TRANSFER_OF_RIGHTS_KEY_PREFIX, []string{"value-title", LICENSEE1})

	transferOfRights := &TransferOfRights{
		"value-title",
		LICENSEE1,
	}
	transferOfRightsBytes, _ := json.Marshal(transferOfRights)
	checkState(t, stub, transferOfRightsKey, string(transferOfRightsBytes))


	// check balances
	walletKey, _ := stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{LICENSEE1})
	checkState(t, stub, walletKey, strconv.Itoa(DEFAULT_WALLET_BALANCE - 110))
	walletKey, _ = stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{NETWORK_USER_ID})
	checkState(t, stub, walletKey, strconv.Itoa(10))
	walletKey, _ = stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{"value-masterArtist1"})
	checkState(t, stub, walletKey, strconv.Itoa(DEFAULT_WALLET_BALANCE + 10))
	walletKey, _ = stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{"value-publishingArtist1"})
	checkState(t, stub, walletKey, strconv.Itoa(DEFAULT_WALLET_BALANCE + 20))
	walletKey, _ = stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{"value-masterArtist2"})
	checkState(t, stub, walletKey, strconv.Itoa(DEFAULT_WALLET_BALANCE + 30))
	walletKey, _ = stub.CreateCompositeKey(WALLET_KEY_PREFIX, []string{"value-publishingArtist2"})
	checkState(t, stub, walletKey, strconv.Itoa(DEFAULT_WALLET_BALANCE + 40))

	fmt.Println("TestDRM_TransferRights passed")
}
