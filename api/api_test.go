package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"banking-ledger-service/db" // Replace with your actual import path
	"banking-ledger-service/models"
	"banking-ledger-service/queue" // Replace with your actual import path

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

var router *mux.Router

func setup() {
	db.InitPostgres()
	db.InitMongo()
	queue.InitRabbitMQ()

	// Initialize router
	router = mux.NewRouter()
	InitializeRoutes(router)
}

func teardown() {
	queue.CloseRabbitMQ()
	// Clean up database tables (implementation depends on your setup)
	// For example, you might truncate tables in Postgres and delete collections in MongoDB
	postgresDB := db.GetPostgresDB()
	if postgresDB != nil {
		_, err := postgresDB.Exec("TRUNCATE TABLE accounts;") // Reset account table
		if err != nil {
			panic(err) // Or handle more gracefully in production
		}
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestCreateAccountHandler(t *testing.T) {
	setup()
	defer teardown()

	var jsonStr = []byte(`{"initial_balance":100.00}`)
	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := executeRequest(req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["account_id"])
}

func TestProcessTransactionHandler(t *testing.T) {
	setup()
	defer teardown()

	var jsonStr = []byte(`{"initial_balance":100.00}`)
	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := executeRequest(req)

	var createResponse map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &createResponse)
	assert.NoError(t, err)
	accountID := createResponse["account_id"]

	transactionData := map[string]interface{}{
		"amount": 50.00,
		"type":   "deposit",
	}
	requestBody, _ := json.Marshal(transactionData)

	url := "/accounts/" + strconv.Itoa(accountID) + "/transactions"
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr = executeRequest(req)

	assert.Equal(t, http.StatusAccepted, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Transaction submitted for processing", response["status"])
}

func TestGetTransactionHistoryHandler(t *testing.T) {
	setup()
	defer teardown()

	var jsonStr = []byte(`{"initial_balance":100.00}`)
	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := executeRequest(req)

	var createResponse map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &createResponse)
	assert.NoError(t, err)
	accountID := createResponse["account_id"]

	transactionData := map[string]interface{}{
		"amount": 50.00,
		"type":   "deposit",
	}
	requestBody, _ := json.Marshal(transactionData)

	url := "/accounts/" + strconv.Itoa(accountID) + "/transactions"
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr = executeRequest(req)

	url = "/accounts/" + strconv.Itoa(accountID) + "/ledger"
	req, err = http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	rr = executeRequest(req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var ledger []models.Transaction
	err = json.Unmarshal(rr.Body.Bytes(), &ledger)

	assert.NoError(t, err)
	if len(ledger) > 0 {
		assert.Equal(t, accountID, ledger[0].AccountID)
		assert.Equal(t, 50.00, ledger[0].Amount)
		assert.Equal(t, "deposit", ledger[0].Type)
	}
}
