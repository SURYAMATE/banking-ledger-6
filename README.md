# Banking Ledger Service

## Overview

The Banking Ledger Service is a RESTful API built with Go (Golang) that allows users to manage bank accounts, perform transactions (deposits and withdrawals), and retrieve transaction history. The service is designed to handle high loads and ensures ACID-like consistency for core operations to prevent double spending or inconsistent balances. It uses a microservices architecture, with separate components for the API gateway, transaction processing, and data storage, communicating via RabbitMQ.

### Features

-   Create bank accounts with specified initial balances.
-   Perform deposits and withdrawals.
-   Maintain a detailed transaction log (ledger) for each account.
-   Ensure data consistency and prevent double spending.
-   Retrieve transaction history for each account.
-   Asynchronous transaction processing via RabbitMQ.
-   Built with Docker for easy deployment and scalability.

## Technologies Used

-   **Go (Golang)**: The programming language used to build the API and transaction processor.
-   **PostgreSQL**: Relational database for storing account balances.
-   **MongoDB**: NoSQL database for storing transaction logs, enabling efficient querying and scalability.
-   **RabbitMQ**: Message broker for managing transaction requests asynchronously, ensuring reliable message delivery.
-   **Docker**: Containerization technology for deploying the application and its dependencies.
-   **Gorilla Mux**: HTTP request router and dispatcher.
## project structure
banking-ledger-6/
├── api/
│ ├── Dockerfile # Dockerfile for the API Gateway
│ ├── main.go # Entry point for the API Gateway, # API request handlers,# Defines API routes
├── db/
│ ├── postgres.go # PostgreSQL database interactions
│ └── mongo.go # MongoDB database interactions
├── models/
│ ├── account.go # Defines the Account model
│ └── transaction.go # Defines the Transaction and TransactionRequest models
├── queue/
│ ├── rabbitmq.go # RabbitMQ integration
├── transactions/
│ ├── Dockerfile # Dockerfile for the Transaction Processor
│ ├── processor.go # Transaction processing logic
│ └── worker.go # Consumes transactions from RabbitMQ and processes them
├── tests/
│ ├── api_test.go # Integration tests for the API Gateway
│ └── processor_test.go # Unit tests for the Transaction Processor
└── docker-compose.yml # Defines services, networks, and volumes for multi-container Docker applications

## Getting Started

### Prerequisites

Make sure you have the following installed:

-   [Docker](https://www.docker.com/get-started)
-   [Docker Compose](https://docs.docker.com/compose/install/)
-   [Go](https://golang.org/doc/install) (Optional - only required for running tests or local development outside of Docker)

### Installation

1.  Clone the repository:

    ```
    git clone https://github.com/your-username/banking-ledger-service.git
    cd banking-ledger-service
    ```

2.  Build and run the application using Docker Compose:

    ```
    docker-compose up --build
    ```

    This command will:

    -   Build the Docker images for the `api` (API Gateway) and `transaction-processor`.
    -   Start the required services: `api`, `postgres`, `mongo`, `rabbitmq`, and `transaction-processor`.
    -   Create networks and volumes for the services to communicate and persist data.

3.  The API will be accessible at `http://localhost:8080`.

## API Endpoints

The Banking Ledger Service provides the following API endpoints:

### 1. Create Account

-   **Endpoint:** `POST /accounts`
-   **Description:** Creates a new bank account with an initial balance.
-   **Request Body:**

    ```
    {
        "initial_balance": 100.00
    }
    ```

-   **Response:**
    -   **Status Code:** `201 Created`
    -   **Body:**

        ```
        {
            "account_id": 1
        }
        ```

### 2. Process Transaction

#### Deposit Funds

-   **Endpoint:** `POST /accounts/{id}/transactions`
-   **Description:** Deposits funds into an account.
-   **Request Body:**

    ```
    {
        "amount": 50.00,
        "type": "deposit"
    }
    ```

-   **Response:**
    -   **Status Code:** `202 Accepted`
    -   **Body:**

        ```
        {
            "status": "Transaction submitted for processing"
        }
        ```

#### Withdraw Funds

-   **Endpoint:** `POST /accounts/{id}/transactions`
-   **Description:** Withdraws funds from an account.
-   **Request Body:**

    ```
    {
        "amount": 30.00,
        "type": "withdrawal"
    }
    ```

-   **Response:**
    -   **Status Code:** `202 Accepted`
    -   **Body:**

        ```
        {
            "status": "Transaction submitted for processing"
        }
        ```

### 3. Get Transaction History

-   **Endpoint:** `GET /accounts/{id}/ledger`
-   **Description:** Retrieves the transaction history for a specific account.
-   **Response:**
    -   **Status Code:** `200 OK`
    -   **Body:**

        ```
        [
            {
                "id": "605c72b8f1d2f8e3a8e7d5c5",
                "account_id": 1,
                "amount": 50,
                "type": "deposit"
            },
            {
                "id": "605c72b8f1d2f8e3a8e7d5c6",
                "account_id": 1,
                "amount": 30,
                "type": "withdrawal"
            }
        ]
        ```

## Testing the API with Postman

You can test the API endpoints using Postman by following these steps:

1.  Open Postman and create a new request.
2.  Set the request method (GET, POST) and enter the appropriate URL.
3.  For POST requests, set the request body in JSON format and add headers as needed (e.g., `Content-Type: application/json`).
4.  Click on “Send” to execute the request and view the response.

### Example Postman Requests

#### Create Account Request:

-   Method: `POST`
-   URL: `http://localhost:8080/accounts`
-   Headers: `Content-Type: application/json`
-   Body:
    ```
    {
        "initial_balance": 100.00
    }
    ```

#### Withdraw Funds Request:

-   Method: `POST`
-   URL: `http://localhost:8080/accounts/1/transactions` (Replace `1` with the account ID)
-   Headers: `Content-Type: application/json`
-   Body:

    ```
    {
        "amount": 30.00,
        "type": "withdrawal"
    }
    ```

#### Get Transaction History Request:

-   Method: `GET`
-   URL: `http://localhost:8080/accounts/1/ledger` (Replace `1` with the account ID)

## Additional Notes

-   **Scalability**: The use of RabbitMQ for asynchronous transaction processing allows the system to handle a high volume of transactions. The API gateway can be scaled independently from the transaction processors.
-   **Data Consistency**: The system uses PostgreSQL for account balances, which supports ACID properties to ensure data consistency.
-   **Error Handling**:  The code includes basic error checks, but a production-ready application needs comprehensive error handling, logging, and potentially retry mechanisms.
-   **Security**: This example does not include authentication or authorization. You would need to add security measures to protect the API endpoints and data.
-   **Configuration**:  The database and RabbitMQ connection strings are hardcoded.  Use environment variables or a configuration file for a more flexible setup.
-   **Testing**: The example includes basic API tests.  You should also write unit tests for the transaction processor and other components.  Consider using mocking to isolate components during testing.
-   **Dependencies**: Remember to run `go mod tidy` to download and manage dependencies.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs or feature requests.

## License

This project is licensed under the [MIT License](LICENSE).
