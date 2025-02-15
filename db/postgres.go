package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	// Replace with your actual import path
)

var postgresDB *sql.DB

// InitPostgres initializes the PostgreSQL connection.
func InitPostgres() {
	var err error
	connStr := "user=example password=example dbname=banking_db sslmode=disable host=postgres port=5432" // Modified connection string
	postgresDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}

	err = postgresDB.Ping()
	if err != nil {
		log.Fatalf("Error pinging PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	// Create table if not exists
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			balance DECIMAL(15,2) NOT NULL
		);
	`
	_, err = postgresDB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

// CreateAccount creates a new account in the database.
func CreateAccount(initialBalance float64) (int, error) {
	var id int
	err := postgresDB.QueryRow("INSERT INTO accounts(balance) VALUES($1) RETURNING id", initialBalance).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetAccountBalance retrieves the balance of an account.
func GetAccountBalance(accountID int) (float64, error) {
	var balance float64
	err := postgresDB.QueryRow("SELECT balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// UpdateAccountBalance updates the balance of an account.
func UpdateAccountBalance(accountID int, newBalance float64) error {
	_, err := postgresDB.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", newBalance, accountID)
	return err
}

// Implement other necessary database functions like GetAccountBalance, UpdateAccountBalance
func GetPostgresDB() *sql.DB {
	return postgresDB
}
