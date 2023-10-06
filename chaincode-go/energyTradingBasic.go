/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the

wit"License"); you may not use this file except in compliance the License.  You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
	contractapi.Contract
}

// ============================================================================================================================
// User Definitions - The ledger with user
// ============================================================================================================================

type User struct {
	ID        int64        `json:"id"`
	Category  UserCategory `json:"category"`
	CreatedOn int64        `json:"createdOn"`
	UpdatedOn int64        `json:"updatedOn"`
	Location  string       `json:"location"` // For simplicity, using a string; consider more complex representations if needed
	MeterId   string       `json:"meterId"`
	Source    EnergySource `json:"source"`
}

type PlatformContract struct {
	UserID    int64 `json:"userId"`
	CreatedOn int64 `json:"createdOn"`
	UpdatedOn int64 `json:"updatedOn"`
}

// ============================================================================================================================
// Trading Definitions - The ledger with user
// ============================================================================================================================

// Order captures the details of an energy buy or sell bid.
// It includes attributes like total quantity, unit cost, and the total order cost.
// Struct fields are arranged alphabetically to ensure determinism across languages.
// Note: While Golang maintains field order when marshaling to JSON, it doesn't auto-sort them.
type Order struct {
	BidMatchID    int64           `json:"bidMatchId"`
	BidStatus     EnergyBidStatus `json:"bidStatus"`
	CreatedOn     int64           `json:"createdOn"`
	ID            int64           `json:"id"`
	OnMarketPrice string          `json:"onMarketPrice"`
	OrderCost     float64         `json:"status"`
	PaymentID     int64           `json:"paymentId"`
	SlotID        string          `json:"slotId"`
	SlotExecDate  int64           `json:"slotExecDate"`
	TotalQuantity int64           `json:"totalQuantity"`
	UnitCost      float64         `json:"unitCost"`
	UpdatedOn     int64           `json:"updatedOn"`
	UserAction    Action          `json:"action"`
	UserID        int64           `json:"userId"`
}

// BidMatch records the details of a matched bid in the energy market.
// Struct fields are alphabetically ordered for cross-language determinism.
type BidMatch struct {
	BidMatchTms       int64           `json:"bidMatchTms"`
	BidSlot           string          `json:"bidSlot"`
	BidStatus         EnergyBidStatus `json:"bidStatus"`
	BidUnitPrice      int64           `json:"bidUnitPrice"`
	BuyerUserId       int64           `json:"buyerUserId"`
	DeliveredBidUnits float64         `json:"deliveredBidUnits"`
	ID                int64           `json:"id"`
	OriginalBidUnits  float64         `json:"originalBidUnits"`
	SellerUserId      int64           `json:"sellerUserId"`
	TransactionBuyID  int64           `json:"transactionBuyId"`
	TransactionSellID int64           `json:"transactionSellId"`
}

// Payment logs transaction details for energy market payments.
// Struct fields are alphabetically ordered for cross-language determinism.
type Payment struct {
	CreatedOn       int64       `json:"createdOn"`
	ID              string      `json:"id"`
	PaymentDetailId int64       `json:"paymentDetail"`
	PaymentType     PaymentType `json:"paymentType"`
	TotalAmount     float64     `json:"totalAmount"`
	UserID          int64       `json:"userId"`
}

// PaymentDetail captures more granular transaction information.
// It includes attributes like the amount refunded, fees applied, and transaction parties.
// Struct fields are arranged alphabetically for consistent representation.
type PaymentDetail struct {
	ID                      int64   `json:"id"`
	DebitedFrom             string  `json:"debitedFrom"`
	CreditedTo              string  `json:"creditedTo"`
	TotalUnitCost           float64 `json:"totalUnitCost"`
	PlatformFee             float64 `json:"platformFee"`
	TokenAmount             float64 `json:"tokenAmount"`
	BidRefundAmount         float64 `json:"bidRefundAmount"`
	PlatformFeeRefundAmount float64 `json:"platformFeeRefundAmount"`
	TokenAmountRefund       float64 `json:"tokenAmountRefund"`
	PenaltyFromSeller       float64 `json:"penaltyFromSeller"`
}

// ============================================================================================================================
// Prefix Definitions - For creating composite keys and avoid id overlap (for future use)
// ============================================================================================================================

const BuyBidPrefix = "BuyBid"
const SellBidPrefix = "SellBid"

// ============================================================================================================================
// Enum Definitions - Absolute states of allowed status for different assets (WIP)
// ============================================================================================================================

type EnergyBidStatus int64
type EnergySource int64
type Action int64
type UserCategory int64
type PaymentType int64

const (
	BidCreated    EnergyBidStatus = iota // = 0
	BidAccepted                          // = 1
	BidRejected                          // = 2
	BidExecuted                          // = 3
	BidTerminated                        // = 4
)

var (
	energyBidMap = map[string]EnergyBidStatus{
		"BidCreated":    BidCreated,
		"BidAccepted":   BidAccepted,
		"BidRejected":   BidRejected,
		"BidExecuted":   BidExecuted,
		"BidTerminated": BidTerminated,
	}
)

func EnergyBidStatusString(status EnergyBidStatus) string {
	return []string{"BidCreated", "BidAccepted", "BidRejected", "BidExecuted", "BidTerminated"}[status]
}

const (
	Solar   EnergySource = iota // = 0
	Wind                        // = 1
	DGSet                       // = 2
	Battery                     // = 3
)

var (
	energySourceMap = map[string]EnergySource{
		"Solar":   Solar,
		"Wind":    Wind,
		"DG Set":  DGSet,
		"Battery": Battery,
	}
)

func EnergySourceString(status EnergySource) string {
	return []string{"Solar", "Wind", "DG Set", "Battery"}[status]
}

const (
	Buy  Action = iota // = 0
	Sell               // = 1
)

var (
	actionMap = map[string]Action{
		"Buy":  Buy,
		"Sell": Sell,
	}
)

func ActionString(status Action) string {
	return []string{"Buy", "Sell"}[status]
}

const (
	Prosumer UserCategory = iota // = 0
	Consumer                     // = 1
)

var (
	userCategoryMap = map[string]UserCategory{
		"Prosumer": Prosumer,
		"Consumer": Consumer,
	}
)

func UserCategoryString(status UserCategory) string {
	return []string{"Prosumer", "Consumer"}[status]
}

const (
	WalletRecharge              PaymentType = iota // = 0
	SellerTokenAmount                              // = 1
	BuyerEnergyPurchased                           // = 2
	BuyerSellerIncentive                           // = 3
	SellerEnergySoldTokenRefund                    // = 4
)

var (
	paymentTypeMap = map[string]PaymentType{
		"WalletRecharge":                         WalletRecharge,
		"Seller - Token Amount":                  SellerTokenAmount,
		"Buyer - Energy Purchased":               BuyerEnergyPurchased,
		"Buyer/Seller - Incentive":               BuyerSellerIncentive,
		"Seller - Energy Sold plus Token Refund": SellerEnergySoldTokenRefund,
	}
)

func PaymentTypeString(status PaymentType) string {
	return []string{"WalletRecharge", "Seller - Token Amount", "Buyer - Energy Purchased", "Buyer/Seller - Incentive", "Seller - Energy Sold plus Token Refund"}[status]
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}

// ============================================================================================================================
// Init - initialize the chaincode (needed for interface) - returning empty success response
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "Write" {
		return Write(stub, args)
	} else if function == "UpdateUserProfile" {
		return UpdateUserProfile(stub, args)
	} else if function == "SignPlatformContract" {
		return SignPlatformContract(stub, args)
	} else if function == "RecordPayment" {
		return RecordPayment(stub, args)
	} else if function == "RegisterOrder" {
		return RegisterOrder(stub, args)
	} else if function == "ProcessBidMatch" {
		return ProcessBidMatch(stub, args)
	} else if function == "ReadUserProfile" {
		return ReadUserProfile(stub, args)
	} else if function == "ReadPlatformContract" {
		return ReadPlatformContract(stub, args)
	} else if function == "ReadPayment" {
		return ReadPayment(stub, args)
	} else if function == "ReadPaymentDetail" {
		return ReadPaymentDetail(stub, args)
	} else if function == "ReadOrder" {
		return ReadOrder(stub, args)
	} else if function == "ReadBidMatch" {
		return ReadBidMatch(stub, args)
	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}

// ============================================================================================================================
// Query - legacy function (needed for interface)
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}
