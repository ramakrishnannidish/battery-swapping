package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	// Create a mock stub
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Positive Test Case: Writing a key-value pair
	t.Run("Successfully Write key-value pair", func(t *testing.T) {
		// Invoke the 'Write' function with a key and a value
		response := stub.MockInvoke("1", [][]byte{[]byte("Write"), []byte("TestKey"), []byte("TestValue")})

		// Assert that the function completed successfully
		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		// Fetch the value for the key from the ledger and check it
		value, err := stub.GetState("TestKey")
		assert.NoError(t, err, "Error getting value from ledger")
		assert.Equal(t, "TestValue", string(value), "Incorrect value retrieved from ledger")
	})

	// Negative Test Case: Wrong number of arguments
	t.Run("Invalid number of arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{[]byte("Write"), []byte("OnlyOneArg")})

		// Assert that the function did not complete successfully
		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})
}

func TestUpdateUserProfile(t *testing.T) {
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Test Case 1: Successfully Update User Profile
	t.Run("Successfully Update User Profile", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("UpdateUserProfile"),
			[]byte("1"),
			[]byte("Prosumer"),
			[]byte("Location 1"),
			[]byte("MeterId 1"),
			[]byte("Solar"),
		})

		// Assert the function completed successfully
		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		// Check if the user is stored in the ledger
		userAsBytes, err := stub.GetState("1")
		assert.NoError(t, err, "Error getting value from ledger")

		var user User
		err = json.Unmarshal(userAsBytes, &user)
		assert.NoError(t, err, "Error unmarshalling user")
		assert.Equal(t, "Location 1", user.Location, "Incorrect value retrieved from ledger")
	})

	// Test Case 2: Incorrect Number of Arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("UpdateUserProfile"),
			[]byte("1"),
			[]byte("Prosumer"),
		})

		// Assert the function did not complete successfully
		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})

	// Test Case 3: Invalid User ID
	t.Run("Invalid User ID", func(t *testing.T) {
		response := stub.MockInvoke("3", [][]byte{
			[]byte("UpdateUserProfile"),
			[]byte("InvalidID"),
			[]byte("Prosumer"),
			[]byte("Location 2"),
			[]byte("MeterId 2"),
			[]byte("Solar"),
		})

		// Assert the function did not complete successfully
		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})

	// Test Case 4: Invalid User Category
	t.Run("Invalid User Category", func(t *testing.T) {
		response := stub.MockInvoke("4", [][]byte{
			[]byte("UpdateUserProfile"),
			[]byte("2"),
			[]byte("InvalidCategory"),
			[]byte("Location 3"),
			[]byte("MeterId 3"),
			[]byte("Solar"),
		})

		// Assert the function did not complete successfully
		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})

	// Test Case 5: Invalid Energy Source
	t.Run("Invalid Energy Source", func(t *testing.T) {
		response := stub.MockInvoke("5", [][]byte{
			[]byte("UpdateUserProfile"),
			[]byte("3"),
			[]byte("Prosumer"),
			[]byte("Location 4"),
			[]byte("MeterId 4"),
			[]byte("InvalidSource"),
		})

		// Assert the function did not complete successfully
		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})
}

func TestSignPlatformContract(t *testing.T) {
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	key := "12345"
	id, err := strconv.ParseInt(key, 10, 64)

	// Creating a dummy user to be used in the tests.
	user := User{ID: int64(id)}
	userBytes, _ := json.Marshal(user)
	// Start a transaction
	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	err = stub.PutState(key, userBytes)
	if err != nil {
		t.Fatalf("Failed to put the user into the stub: %s", err.Error())
	}

	// Test Case 1: Successfully Sign Platform Contract
	t.Run("Successfully Sign Platform Contract", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{[]byte("SignPlatformContract"), []byte("12345")})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		contractAsBytes, err := stub.GetState("PlatformContract_12345")
		assert.NoError(t, err, "Error getting value from ledger")

		var contract PlatformContract
		err = json.Unmarshal(contractAsBytes, &contract)
		assert.NoError(t, err, "Error unmarshalling contract")
		assert.Equal(t, int64(12345), contract.UserID, "Incorrect UserID in contract")
	})

	// Test Case 2: Incorrect Number of Arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{[]byte("SignPlatformContract"), []byte("12345"), []byte("ExtraArg")})

		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})

	// Test Case 3: Non-existing User
	t.Run("Non-existing User", func(t *testing.T) {
		response := stub.MockInvoke("3", [][]byte{[]byte("SignPlatformContract"), []byte("54321")})

		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})

	// Test Case 4: Invalid User ID
	t.Run("Invalid User ID", func(t *testing.T) {
		response := stub.MockInvoke("4", [][]byte{[]byte("SignPlatformContract"), []byte("InvalidUserID")})

		assert.NotEqual(t, shim.OK, response.GetStatus(), "Function unexpectedly succeeded")
	})
}

func TestRegisterOrder(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Test Case 1: Successfully register a new order
	t.Run("Successfully Register a New Order", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("RegisterOrder"),
			[]byte("1"),        // bidMatchID
			[]byte("0"),        // bidStatus
			[]byte("100"),      // endTime
			[]byte("4"),        // orderID
			[]byte("2.5"),      // onMarketPrice
			[]byte("200"),      // orderCost
			[]byte("5"),        // paymentID
			[]byte("slot1234"), // slotID
			[]byte("300"),      // totalQuantity
			[]byte("3.5"),      // unitCost
			[]byte("6"),        // userID
			[]byte("50"),       // slotExecDate
			[]byte("Buy"),      // action
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		orderAsBytes, err := stub.GetState("Order_4")
		assert.NoError(t, err, "Error getting order from ledger")

		var order Order
		err = json.Unmarshal(orderAsBytes, &order)
		assert.NoError(t, err, "Error unmarshalling order")

		assert.Equal(t, int64(4), order.ID, "Order ID mismatch")
		assert.Equal(t, int64(1), order.BidMatchID, "BidMatchID mismatch")
	})

	// Test Case 2: Provide incorrect number of arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("RegisterOrder"),
			[]byte("1"),
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "Incorrect number of arguments")
	})

	// Test Case 3: Successfully update an existing order
	t.Run("Successfully Register a New Order", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("RegisterOrder"),
			[]byte("1"),        // bidMatchID
			[]byte("0"),        // bidStatus
			[]byte("100"),      // endTime
			[]byte("4"),        // orderID
			[]byte("2.5"),      // onMarketPrice
			[]byte("200"),      // orderCost
			[]byte("5"),        // paymentID
			[]byte("slot1234"), // slotID
			[]byte("300"),      // totalQuantity
			[]byte("3.5"),      // unitCost
			[]byte("6"),        // userID
			[]byte("50"),       // slotExecDate
			[]byte("Buy"),      // action
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		orderAsBytes, err := stub.GetState("Order_4")
		assert.NoError(t, err, "Error getting order from ledger")

		var order Order
		err = json.Unmarshal(orderAsBytes, &order)
		assert.NoError(t, err, "Error unmarshalling order")

		assert.Equal(t, int64(4), order.ID, "Order ID mismatch")
		assert.Equal(t, int64(1), order.BidMatchID, "BidMatchID mismatch")

		response = stub.MockInvoke("1", [][]byte{
			[]byte("RegisterOrder"),
			[]byte("1"),        // bidMatchID
			[]byte("1"),        // bidStatus
			[]byte("100"),      // endTime
			[]byte("4"),        // orderID
			[]byte("2.5"),      // onMarketPrice
			[]byte("200"),      // orderCost
			[]byte("5"),        // paymentID
			[]byte("slot1234"), // slotID
			[]byte("300"),      // totalQuantity
			[]byte("3.5"),      // unitCost
			[]byte("6"),        // userID
			[]byte("50"),       // slotExecDate
			[]byte("Buy"),      // action
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		orderAsBytes, err = stub.GetState("Order_4")
		assert.NoError(t, err, "Error getting order from ledger")

		err = json.Unmarshal(orderAsBytes, &order)
		assert.NoError(t, err, "Error unmarshalling order")

		assert.Equal(t, int64(4), order.ID, "Order ID mismatch")
		assert.Equal(t, int64(1), order.BidMatchID, "BidMatchID mismatch")
	})

}

func TestProcessBidMatch(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Test Case 1: Successfully process a new BidMatch
	t.Run("Successfully Process a New BidMatch", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("ProcessBidMatch"),
			[]byte("1"),
			[]byte("Slot1"),
			[]byte("1"),
			[]byte("100"),
			[]byte("4"),
			[]byte("2.5"),
			[]byte("1"),
			[]byte("3.5"),
			[]byte("5"),
			[]byte("6"),
			[]byte("7"),
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), "Unexpected error: "+response.GetMessage())

		bidMatchAsBytes, err := stub.GetState("BidMatch_1")
		assert.NoError(t, err, "Error getting BidMatch from ledger")

		var bidMatch BidMatch
		err = json.Unmarshal(bidMatchAsBytes, &bidMatch)
		assert.NoError(t, err, "Error unmarshalling BidMatch")

		assert.Equal(t, int64(1), bidMatch.ID, "BidMatch ID mismatch")
		assert.Equal(t, int64(4), bidMatch.BuyerUserId, "BuyerUserId mismatch")
	})

	// Test Case 2: Provide incorrect number of arguments
	t.Run("Incorrect Number of Arguments", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("ProcessBidMatch"),
			[]byte("1"),
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "Incorrect number of arguments")
	})
}

func TestReadOrder(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Registering a new order
	response := stub.MockInvoke("1", [][]byte{
		[]byte("RegisterOrder"),
		[]byte("1"),        // bidMatchID
		[]byte("0"),        // bidStatus
		[]byte("100"),      // endTime
		[]byte("4"),        // orderID
		[]byte("2.5"),      // onMarketPrice
		[]byte("200"),      // orderCost
		[]byte("5"),        // paymentID
		[]byte("slot1234"), // slotID
		[]byte("300"),      // totalQuantity
		[]byte("3.5"),      // unitCost
		[]byte("6"),        // userID
		[]byte("50"),       // slotExecDate
		[]byte("Buy"),      // action
	})
	assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

	// Test Case: Successfully read an order
	t.Run("Successfully Read an Order", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("ReadOrder"),
			[]byte("4"), // orderID
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		var order Order
		err := json.Unmarshal(response.GetPayload(), &order)
		assert.NoError(t, err, "Error unmarshalling order")

		assert.Equal(t, int64(4), order.ID, "Order ID mismatch")
		assert.Equal(t, int64(1), order.BidMatchID, "BidMatchID mismatch")
	})

	// Test Case: Try to read an order that doesn't exist
	t.Run("Try to Read Nonexistent Order", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("ReadOrder"),
			[]byte("99"), // orderID
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "not found")
	})
}

func TestReadBidMatch(t *testing.T) {
	// Mock stub creation
	stub := shimtest.NewMockStub("testingStub", new(SimpleChaincode))

	// Registering a new BidMatch
	response := stub.MockInvoke("1", [][]byte{
		[]byte("ProcessBidMatch"),
		[]byte("1"),     // ID
		[]byte("Slot1"), // other args...
		[]byte("1"),
		[]byte("100"),
		[]byte("4"),
		[]byte("2.5"),
		[]byte("1"),
		[]byte("3.5"),
		[]byte("5"),
		[]byte("6"),
		[]byte("7"),
	})
	assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

	// Test Case 1: Successfully read a BidMatch
	t.Run("Successfully Read a BidMatch", func(t *testing.T) {
		response := stub.MockInvoke("1", [][]byte{
			[]byte("ReadBidMatch"),
			[]byte("1"), // bidMatchID
		})

		assert.Equal(t, int32(shim.OK), response.GetStatus(), fmt.Sprintf("Unexpected error: %s", response.GetMessage()))

		var bidMatch BidMatch
		err := json.Unmarshal(response.GetPayload(), &bidMatch)
		assert.NoError(t, err, "Error unmarshalling BidMatch")

		assert.Equal(t, int64(1), bidMatch.ID, "BidMatch ID mismatch")
		assert.Equal(t, int64(4), bidMatch.BuyerUserId, "BuyerUserId mismatch")
	})

	// Test Case 2: Try to read a BidMatch that doesn't exist
	t.Run("Try to Read Nonexistent BidMatch", func(t *testing.T) {
		response := stub.MockInvoke("2", [][]byte{
			[]byte("ReadBidMatch"),
			[]byte("99"), // bidMatchID
		})

		assert.Equal(t, int32(shim.ERROR), response.GetStatus(), "Function unexpectedly succeeded")
		assert.Contains(t, response.GetMessage(), "not found")
	})
}
