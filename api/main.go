package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"banking-ledger-service/db" // Replace with your actual import path
	"banking-ledger-service/models"
	"banking-ledger-service/queue" // Replace with your actual import path

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connections
	db.InitPostgres()
	db.InitMongo()
	queue.InitRabbitMQ()

	// Setup router
	router := mux.NewRouter()
	InitializeRoutes(router) // Use the handlers package

	// Start server
	port := 8080
	log.Printf("Server started on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

// InitializeRoutes sets up the API routes.  This is now in the `handlers` package.
func InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/accounts", createAccountHandler).Methods("POST")
	router.HandleFunc("/accounts/{id}/transactions", processTransactionHandler).Methods("POST")
	router.HandleFunc("/accounts/{id}/ledger", getTransactionHistoryHandler).Methods("GET")
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		InitialBalance float64 `json:"initial_balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accountID, err := db.CreateAccount(reqBody.InitialBalance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"account_id": accountID})
}

func processTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Account ID is required", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Amount float64 `json:"amount"`
		Type   string  `json:"type"` // "deposit" or "withdrawal"
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Basic Validation
	if reqBody.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	if reqBody.Type != "deposit" && reqBody.Type != "withdrawal" {
		http.Error(w, "Invalid transaction type", http.StatusBadRequest)
		return
	}

	// Send transaction to queue
	transactionData := models.TransactionRequest{
		AccountID: accountID,
		Amount:    reqBody.Amount,
		Type:      reqBody.Type,
	}

	err = queue.PublishTransaction(transactionData)
	if err != nil {
		http.Error(w, "Failed to publish transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted) // Indicate that the request is being processed.
	json.NewEncoder(w).Encode(map[string]string{"status": "Transaction submitted for processing"})
}

func getTransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Account ID is required", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	ledger, err := db.GetTransactionHistory(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ledger)
}
