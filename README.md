# Wallet App

## Overview

This is a simple Wallet App demo implemented in Go. It provides a RESTful API that allows users to:

- Deposit money into their wallet
- Withdraw money from their wallet
- Send money to another user
- Check their wallet balance
- View their transaction history

## Folder Structure

```md
.
├── .gitignore                  # Git ignore file, excludes unnecessary files and folders from the repository
├── LICENSE                     # Project license
├── README.md                   # Main project documentation

├── config                      # Configuration settings for the project
│   ├── config.go               # Application configuration management code
│   ├── database.go             # Database connection setup
│   └── config.yml              # Configuration file (e.g., environment variables)

├── controllers                 # API Controllers for handling HTTP requests
│   ├── wallet.go               # Controller for wallet-related endpoints
│   ├── wallet_test.go          # Unit tests for wallet controller
│   └── dto.go                  # Data transfer objects (DTOs) for API request/response validation

├── docs                        # Documentation files for project design and usage
│   ├── system_design.md        # System design documentation

├── docker-compose.yml          # Docker Compose configuration for local environment setup
├── Dockerfile                  # Dockerfile for building the project container

├── main.go                     # Main application entry point

├── middlewares                 # Middleware functions for request handling
│   ├── auth.go                 # Authentication middleware
│   └── cors.go                 # CORS (Cross-Origin Resource Sharing) middleware

├── mocks                       # Mock services for testing
│   ├── mock_user_service.go    # Mock UserService for unit tests
│   └── mock_wallet_service.go  # Mock WalletService for unit tests

├── models                      # Database models representing core entities
│   ├── user.go                 # User model
│   ├── transaction.go          # Transaction model
│   └── vault.go                # Vault model

├── routes                      # API route definitions and setup
│   └── routes.go               # Router and API endpoint setup

├── scripts                     # Utility scripts for project setup and data management
│   └── setup-fixtures.sh       # Script to set up initial data or fixtures in the database

├── services                    # Business logic and service layer
│   ├── user.go                 # UserService containing user-related business logic
│   ├── user_test.go            # Unit tests for UserService
│   ├── wallet.go               # WalletService containing wallet-related business logic
│   └── wallet_test.go          # Unit tests for WalletService

├── utils                       # Utility functions and helper methods
│   ├── auth.go                 # Authorization helper functions
│   ├── auth_test.go            # Unit tests for authorization helpers
│   ├── pagination.go           # Pagination helper functions
│   ├── pagination_test.go      # Unit tests for pagination helpers
│   └── response.go             # Unified API response formatting functions

├── go.mod                      # Go module dependencies and versions
└── go.sum                      # Go module dependency checksums
```

## Setup and Run

### Prerequisites

- Go 1.22 or higher
- PostgreSQL database
- Docker (for containerized setup and unit testing)

### Steps

#### 1. Clone the Repository

```bash
git clone https://github.com/wanliqun/go-wallet-app.git
cd go-wallet-app
```

#### 2. Configure Environment Variables

Create a `.env` file in the root directory with the following content, adjusting the values as needed:

```env
APP_DATABASE_HOST={your_postgres_host}
APP_DATABASE_PORT={your_postgres_port}
APP_DATABASE_USER={your_database_user}
APP_DATABASE_PASSWORD={your_database_password}
```

#### 3. Install Dependencies

```bash
go mod tidy
```

#### 4. Run Database Migrations

No manual action required. The database and tables are automatically migrated at application startup.

#### 5. Run the Application

```bash
go run main.go
```

Or using Docker:

```bash
docker-compose up --build
```
#### 6. Setup test fixtures (Optional)

If you are using docker and want to load sample data into the database for testing, you can run the following command:

```bash
go run scripts/setup-fixtures.sh
```

## Project Retrospective

### Features Not Implemented

- User Authentication: JWT or session-based authentication was not implemented due to time constraints.
- Advanced Error Codes: Error responses are currently simplified.
- Integration and End-to-End Testing: Comprehensive tests for entire workflows are not yet in place.
- Additional Test Cases: More thorough test cases should be considered for greater coverage.

### Areas for Improvement

- Authentication: Enhance the current basic user authentication and authorization mechanisms.
- Input Validation: Enhance validation for request payloads.
- Logging: Introduce structured logging for better monitoring and debugging.
- Metrics and Alerts: Implement metrics and alerting to track system health and performance.
- Redis Caching: Consider caching frequent reads, such as balance checks, to improve response times.
- Scalability: Add rate limiting and plan for horizontal scaling to support future growth.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.