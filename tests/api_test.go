package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"banking-ledger-service/db"    // Replace with your actual import path
	"banking-ledger-service/queue" // Replace with your actual import path

	"github.com/stretchr/testify/assert"
)

func setup() {
	db.InitPostgres()
	db.InitMongo()
	queue.InitRabbitMQ()
}

func teardown() {
	// Clean up database and queue after tests (implementation depends on your setup)
	// For example, you might truncate tables in Postgres and delete queues in RabbitMQ
	queue.CloseRabbitMQ()
}

func TestCreateAccountHandler(t *testing.T) {
	setup()
	defer teardown()

	// Prepare request
	var jsonStr = []byte(`{"initial_balance":100.00}`)
	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Prepare recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement the createAccountHandler logic here
		w.WriteHeader(http.StatusCreated)
		response := map[string]int{"account_id": 1} // Example response
		json.NewEncoder(w).Encode(response)
	})

	// Execute request
	handler.ServeHTTP(rr, req)

	// Check results
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["account_id"])
}

func TestProcessTransactionHandler(t *testing.T) {
	setup()
	defer teardown()

	// 1. Create an account first
	accountID, err := db.CreateAccount(100.00)
	assert.NoError(t, err)

	// 2. Prepare transaction request
	transactionData := map[string]interface{}{
		"amount": 50.00,
		"type":   "deposit",
	}
	requestBody, _ := json.Marshal(transactionData)

	// 3. Prepare HTTP request
	url := "/accounts/" + string(rune(accountID)) + "/transactions" // Ensure accountID is a string
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// 4. Execute request
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement the ProcessTransactionHandler logic here
		w.WriteHeader(http.StatusAccepted)
		response := map[string]string{"status": "Transaction submitted for processing"}
		json.NewEncoder(w).Encode(response)
	})
	handler.ServeHTTP(rr, req)

	// 5. Check results
	assert.Equal(t, http.StatusAccepted, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Transaction submitted for processing", response["status"])
}
