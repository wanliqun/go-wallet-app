#!/bin/bash

# Configure database connection information, supporting environment variables
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-wallet_db}"

# Stop and remove any existing containers
echo "Stopping and removing any existing containers..."
docker-compose down

# Start up services defined in docker-compose
echo "Starting up docker containers..."
docker-compose up -d

# Wait for the database to be ready
echo "Waiting for Postgres to be ready..."
until docker exec postgres13 pg_isready -U "$DB_USER" -h "$DB_HOST" > /dev/null 2>&1; do
  echo -n "."
  sleep 1
done
echo "Postgres is ready!"

# Populate the database with test data
echo "Populating database with test data..."

# Build the SQL commands
SQL_COMMANDS=""

# Insert user data (10 users)
SQL_COMMANDS+="
-- Insert user data
INSERT INTO users (name, email, created_at, updated_at) VALUES
"

for i in {1..10}; do
  if [ $i -lt 10 ]; then
    SQL_COMMANDS+="('testuser_$i', 'testuser_$i@example.com', NOW(), NOW()),"
  else
    SQL_COMMANDS+="('testuser_$i', 'testuser_$i@example.com', NOW(), NOW());"
  fi
done

# Insert vault data (BTC and ETH for each user)
SQL_COMMANDS+="
-- Insert vault data
INSERT INTO vaults (user_id, currency, amount, created_at, updated_at) VALUES
"

for i in {1..10}; do
  BTC_AMOUNT=$((1000 + $i * 100))
  ETH_AMOUNT=$((500 + $i * 100))
  if [ $i -lt 10 ]; then
    SQL_COMMANDS+="
    ($i, 'BTC', $BTC_AMOUNT, NOW(), NOW()),
    ($i, 'ETH', $ETH_AMOUNT, NOW(), NOW()),
    "
  else
    SQL_COMMANDS+="
    ($i, 'BTC', $BTC_AMOUNT, NOW(), NOW()),
    ($i, 'ETH', $ETH_AMOUNT, NOW(), NOW());
    "
  fi
done

# Prepare batch insert for transactions
echo "Generating transactions for all users..."

SQL_COMMANDS+="
-- Insert transaction data
INSERT INTO transactions (user_id, counterparty_id, type, amount, currency, memo, timestamp, id) VALUES
"

TRANSACTION_VALUES=""

for i in {1..10}; do
  for j in {1..50}; do
    TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
    if (( j % 3 == 0 )); then
      TYPE="deposit"
      AMOUNT=$((RANDOM % 100 + 10))
      MEMO="Deposit transaction $j for user $i"
      COUNTERPARTY="NULL"
      CURRENCY="BTC"
    elif (( j % 3 == 1 )); then
      TYPE="withdrawal"
      AMOUNT=$((RANDOM % 50 + 5))
      MEMO="Withdrawal transaction $j for user $i"
      COUNTERPARTY="NULL"
      CURRENCY="ETH"
    else
      TYPE="transfer_out"
      COUNTERPARTY=$(((i % 10) + 1))
      AMOUNT=$((RANDOM % 30 + 1))
      MEMO="Transfer transaction $j from user $i to user $COUNTERPARTY"
      CURRENCY="BTC"
    fi

    # Add transaction values
    if [ "$TYPE" = "transfer_out" ]; then
      # Transfer out transaction for sender
      TRANSACTION_VALUES+="($i, $COUNTERPARTY, '$TYPE', $AMOUNT, '$CURRENCY', '$MEMO', '$TIMESTAMP', DEFAULT),"
      # Transfer in transaction for recipient
      TRANSACTION_VALUES+="($COUNTERPARTY, $i, 'transfer_in', $AMOUNT, '$CURRENCY', '$MEMO', '$TIMESTAMP', DEFAULT),"
    else
      # Deposit or withdrawal
      TRANSACTION_VALUES+="($i, $COUNTERPARTY, '$TYPE', $AMOUNT, '$CURRENCY', '$MEMO', '$TIMESTAMP', DEFAULT),"
    fi
  done
done

# Remove the last comma and add a semicolon
TRANSACTION_VALUES=${TRANSACTION_VALUES%,}
TRANSACTION_VALUES+=";"

SQL_COMMANDS+="$TRANSACTION_VALUES"

# Execute the SQL commands
echo "Executing batch insert..."
echo "$SQL_COMMANDS" | docker exec -i postgres13 psql -U "$DB_USER" -d "$DB_NAME"

echo "Database has been populated with test data."