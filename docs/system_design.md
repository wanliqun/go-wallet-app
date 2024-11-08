# Wallet App Design

## Requirement Analysis

### Functional Requirements

#### 1. Main Business Objects

- User: Represents each account holder on the platform:
    - name: User name for display or login
    - email: Email address for verification or notification
- Vault: Stores funds in a specific currency for each user (replacing "wallet" for clarity):
    - currency: Type of currency (e.g., USDT, BTC)
    - amount: Current balance in that currency
- Transaction: Records every movement of funds:
    - type: Type of transaction (e.g., deposit, withdrawal, transfer)
    - user: User initiating or receiving the transaction
    - currency: Currency type involved in the transaction
    - amount: Transaction amount
    - memo: Optional note or reference for the transaction
    - timestamp: Date and time of the transaction

#### 2. Relationships Between Main Objects

|               | **User**                    | **Vault**                           | **Transaction**                        |
|---------------|-----------------------------|-------------------------------------|----------------------------------------|
| **User**      | Can transfer                | Can deposit, withdraw, check        | Can view transaction history           |
| **Vault**     | Belongs to a user           | N/A                                 | Involved in transactions               |
| **Transaction** | Triggered by a user       | Affects vault balance               | N/A       |

#### 3. Access Patterns

- User-Centric Queries:
    - Get all vault balances for a user: (User → Vaults)
    - Get all transaction history for a user: (User → Transactions)
- Vault Operations:
    - Deposit funds into a vault
    - Withdraw funds from a vault
- User-to-User Operations:
    -Transfer funds from one user to another

### Non-Functional Requirements

#### 1. Performance

- Balance Retrieval: Checking vault balances is a frequent operation, especially if shown on the app’s landing page. Ideally, balance retrieval should occur within 200ms.
- Transaction Retrieval: Retrieving recent transactions should be quick, with an aim for under 500ms latency.

#### 2. Availability

- High Availability: Given the financial nature of the app, availability is crucial to maintaining user trust. Aim for 99.99% availability.
- Transactional Integrity: Operations like transfers must be atomic, ensuring consistency. Some brief downtime may be tolerable but should be minimized.

#### 3. Security

- Authentication: Two-Factor Authentication (2FA) should be implemented for actions like withdrawals and transfers to prevent unauthorized access.
- Data Encryption: Ensure encryption for data-at-rest and data-in-transit, especially sensitive data.
- Access Control: Implement role-based access control (RBAC) to manage access levels for different user roles.

#### 4. Scalability

- User Base Growth: The system should scale to accommodate an increasing user base with thousands of concurrent users.
- Traffic Peaks: Design for scalable handling of traffic spikes, particularly during peak times or events that increase transaction volume.

## Database Design

### Table Schema

#### User Table

| Column   | Data Type           | Constraints                        | Description                         |
|----------|----------------------|------------------------------------|-------------------------------------|
| id       | `UNSIGNED INT(4)`   | `PRIMARY KEY`, `AUTO_INCREMENT`    | Unique identifier for each user     |
| name     | `VARCHAR(16)`       | `NOT NULL`, `UNIQUE`               | User’s unique name                  |
| email    | `VARCHAR(32)`       | `NOT NULL`, `UNIQUE`               | User’s unique email address         |

#### Vault Table

| Column   | Data Type           | Constraints                        | Description                                      |
|----------|----------------------|------------------------------------|--------------------------------------------------|
| id       | `UNSIGNED INT(4)`   | `PRIMARY KEY`, `AUTO_INCREMENT`    | Unique identifier for each vault                 |
| user_id  | `UNSIGNED INT(4)`   | `NOT NULL`, `INDEX (idx_user_id)`  | Foreign key referencing `User.id`                |
| currency | `VARCHAR(32)`       | `NOT NULL`                         | Type of currency (e.g., USDT, BTC)               |
| amount   | `NUMERIC(36, 18)`   | `DEFAULT 0`                        | Current balance in the specified currency        |

#### Transactions Table

| Column         | Data Type           | Constraints                                | Description                                                     |
|----------------|---------------------|--------------------------------------------|-----------------------------------------------------------------|
| id             | `UNSIGNED INT(4)`   | `PRIMARY KEY`, `AUTO_INCREMENT`            | Unique identifier for each transaction                          |
| user_id        | `UNSIGNED INT(4)`   | `NOT NULL`                                 | Foreign key referencing `User.id`; primary user in the txn      |
| counterpart_id | `UNSIGNED INT(4)`   | `DEFAULT NULL`                             | Foreign key referencing `User.id`; other user in a transfer     |
| type           | `VARCHAR(16)`      | `NOT NULL` | Type of transaction (deposit, withdraw, transfer in/out)       |
| amount         | `NUMERIC(36, 18)`   | `NOT NULL`                                 | Amount of currency in the transaction                           |
| currency       | `VARCHAR(32)`       | `NOT NULL`                                 | Currency type (matches `Vault.currency`)                        |
| memo           | `VARCHAR(256)`      | `NULL`                                     | Optional note for transaction                                   |
| timestamp      | `DATETIME`          | `DEFAULT CURRENT_TIMESTAMP`                | Timestamp of transaction creation                               |

---

### Notes

- **Amount Precision**: The `NUMERIC(64, 0)` type supports extremely large values, suitable for cryptocurrency balances with different precisions.
- **Transfer Records**: Two entries are created per transfer transaction—`transfer_out` for the sender and `transfer_in` for the recipient—allowing simple queries for all user-related transactions.
- **Keyset Pagination** The `(user_id, timestamp, id)` composite index is specifically designed to support efficient transaction history queries involving specific users, especially for keyset pagination. The index is structured to efficiently support paginated queries by user:
  - **user_id** as the first column, allowing the index to quickly filter all transactions related to a specific user.
  - **timestamp** as the second column, ensuring efficient ordering by time, which is critical for retrieving the most recent transactions.
  - **id** is the third column, providing uniqueness and ensuring consistent sorting when multiple transactions occur at the same timestamp. This helps keyset pagination avoid inconsistent results due to ties in timestamp values.

# API Design

## Authentication and Security

- **Authentication**: All API methods require authorization. A Bearer token must be provided in the `Authorization` header of the request. For simplification, we'll temporarily use the username as the Bearer token to identify the acting user.

  **Note**: Using usernames as Bearer tokens poses security risks. We would implement a more secure authentication mechanism, such as OAuth2 or JWT (JSON Web Token) for real production environment.

- **Unified API Response Format**:

  ```json
  {
      "code": 0,
      "message": "ok",
      "result": { ... }
  }
  ```

  - `code`: `0` indicates success; non-zero values indicate a specific error.
  - `message`: "ok" for success; error message if an error occurs.
  - `result`: The result object returned upon success.

## API Endpoints

1. **Deposit**

   - **Method**: `POST /deposit`
   - **Description**: Deposit funds into the user's vault.
   - **Request Parameters**:

     | Parameter | Type                  | Required | Description                             |
     |-----------|-----------------------|----------|-----------------------------------------|
     | currency  | `string`              | Yes      | Currency type (e.g., USDT, BTC)         |
     | amount    | `string` or `decimal` | Yes      | Deposit amount                          |

   - **Response**:

     Success or error message.

2. **Withdraw**

   - **Method**: `POST /withdraw`
   - **Description**: Withdraw funds from the user's vault.
   - **Request Parameters**:

     | Parameter | Type                  | Required | Description                               |
     |-----------|-----------------------|----------|-------------------------------------------|
     | currency  | `string`              | Yes      | Currency type                             |
     | amount    | `string` or `decimal` | Yes      | Withdrawal amount                         |

   - **Response**:

     Success or error message.

3. **Transfer**

   - **Method**: `POST /transfer`
   - **Description**: Transfer funds to another user.
   - **Request Parameters**:

     | Parameter  | Type                  | Required | Description                                 |
     |------------|-----------------------|----------|---------------------------------------------|
     | recipient  | `string`              | Yes      | Recipient's username                        |
     | currency   | `string`              | Yes      | Currency type                               |
     | amount     | `string`              | Yes      | Amount to transfer                          |
     | memo       | `string`              | No       | Transfer notes or description               |

   - **Response**:

     Success or error message.

4. **Get Balance**

   - **Method**: `GET /balances`
   - **Description**: Retrieve balances for the specified currencies for the user. Limits the currencies list to a maximum of 30 items.
   - **Request Parameters**:

     | Parameter | Type  | Required | Description                                   |
     |-----------|-------|----------|-----------------------------------------------|
     | currency  | `[]string` | YES       | List of currency codes to retrieve balances for (max 30)       |

   - **Response**:

     ```json
     {
         "code": 0,
         "message": "ok",
         "result": {
             "balances": [
                 {
                     "currency": "USDT",
                     "amount": "1000"
                 }
                 // More balance entries
             ]
         }
     }
     ```

5. **Transaction History**

   - **Method**: `GET /transactions`
   - **Description**: Retrieve the user's transaction history. 
   - **Request Parameters**:

     | Parameter | Type     | Required | Description                                                                    |
     |-----------|----------|----------|--------------------------------------------------------------------------------|
     | cursor    | `int`    | No       | Encoded cursor (timestamp + id) from the last transaction of the previous page |
     | limit     | `int`    | No       | Number of records per page (default `10`, max `50`)                            |
     | type      | `string` | No       | Filter by transaction type (`deposit`, `withdraw`, `transfer_out`, `transfer_in`) |
     | order     | `string` | No       | Sort order: `asc` or `desc` (default `desc`)                                   |

     **Note**: The cursor parameter is used for **keyset pagination**, which improves performance over traditional offset pagination by efficiently querying based on the last transaction’s position.

   - **Response**:

     ```json
     {
         "code": 0,
         "message": "ok",
         "result": {
             "transactions": [
                 {
                     "type": "transfer_out",
                     "counterparty": "alice",
                     "currency": "BTC",
                     "amount": "1",
                     "memo": "Payment for services",
                     "timestamp": "2023-11-04T12:34:56Z"
                 }
                 // More transaction records
             ],
             "next_cursor": "MTczMTA2NjM2MzAwMF82MzQ="
         }
     }
     ```

# Technical Decisions

- Language: Chose Go for its performance and built-in concurrency support.
- Follows MVC architecture for best engineering practice.
- Web Framework: Used Gin for its minimalistic and high-performance HTTP request handling.
- ORM: Utilized GORM for database interactions to simplify CRUD operations.
- Database Transactions: Use retional database (Postgres) to implement transactions where atomicity is required (e.g., transfers).
- Error Handling: Comprehensive error handling and meaningful HTTP responses.
- Testing: Wrote unit tests for critical components to ensure reliability.
- Dockerization: Added Docker support for easy setup and deployment.