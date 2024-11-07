package services_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/wanliqun/go-wallet-app/models"
	"github.com/wanliqun/go-wallet-app/services"
)

func TestDeposit(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	walletService := services.NewWalletService(tx)

	testuser := userGenerator.Generate()
	tx.Create(testuser)

	currency := "USDT"

	t.Run("should deposit successfully", func(t *testing.T) {
		amount := decimal.NewFromFloat(100.0)

		err := walletService.Deposit(testuser.ID, currency, amount)
		assert.NoError(t, err)

		var vault models.Vault
		tx.First(&vault, "user_id = ? AND currency = ?", testuser.ID, currency)
		assert.True(t, amount.Equal(vault.Amount))
	})

	t.Run("should return error for invalid amount", func(t *testing.T) {
		amount := decimal.NewFromFloat(-50.0)

		err := walletService.Deposit(testuser.ID, currency, amount)
		assert.Error(t, err)
		assert.Equal(t, services.ErrInvalidAmount, err)
	})
}

func TestWithdraw(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	walletService := services.NewWalletService(tx)

	testuser := userGenerator.Generate()
	tx.Create(testuser)

	currency := "USDT"
	initialAmount := decimal.NewFromFloat(100.0)
	walletService.Deposit(testuser.ID, currency, initialAmount)

	t.Run("should withdraw successfully", func(t *testing.T) {
		withdrawAmount := decimal.NewFromFloat(50.0)

		err := walletService.Withdraw(testuser.ID, currency, withdrawAmount)
		assert.NoError(t, err)

		var vault models.Vault
		tx.First(&vault, "user_id = ? AND currency = ?", testuser.ID, currency)
		assert.True(t, initialAmount.Sub(withdrawAmount).Equal(vault.Amount))
	})

	t.Run("should return error for insufficient balance", func(t *testing.T) {
		withdrawAmount := decimal.NewFromFloat(200.0)

		err := walletService.Withdraw(testuser.ID, currency, withdrawAmount)
		assert.Error(t, err)
		assert.Equal(t, services.ErrInsufficientBalance, err)
	})
}

func TestTransfer(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	senderUser := userGenerator.Generate()
	recipientUser := userGenerator.Generate()
	tx.CreateInBatches([]*models.User{senderUser, recipientUser}, 2)

	walletService := services.NewWalletService(tx)

	currency := "USDT"
	amount := decimal.NewFromFloat(50.0)

	walletService.Deposit(senderUser.ID, currency, decimal.NewFromFloat(100.0))

	t.Run("should transfer successfully", func(t *testing.T) {
		err := walletService.Transfer(senderUser.ID, recipientUser.ID, currency, amount, "test transfer")
		assert.NoError(t, err)

		var senderVault, recipientVault models.Vault
		tx.First(&senderVault, "user_id = ? AND currency = ?", senderUser.ID, currency)
		tx.First(&recipientVault, "user_id = ? AND currency = ?", recipientUser.ID, currency)

		assert.True(t, decimal.NewFromFloat(50.0).Equal(senderVault.Amount))
		assert.True(t, amount.Equal(recipientVault.Amount))
	})

	t.Run("should return error for insufficient balance", func(t *testing.T) {
		err := walletService.Transfer(senderUser.ID, recipientUser.ID, currency, decimal.NewFromFloat(200.0), "test insufficient balance")
		assert.Error(t, err)
		assert.Equal(t, services.ErrInsufficientBalance, err)
	})

	t.Run("should return error when transferring to self", func(t *testing.T) {
		err := walletService.Transfer(senderUser.ID, senderUser.ID, currency, amount, "self transfer")
		assert.Error(t, err)
		assert.Equal(t, "cannot transfer to self", err.Error())
	})
}

func TestGetBalances(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	testuser := userGenerator.Generate()
	tx.Create(testuser)

	walletService := services.NewWalletService(tx)

	currency1 := "BTC"
	currency2 := "USDT"

	walletService.Deposit(testuser.ID, currency1, decimal.NewFromFloat(100.0))
	walletService.Deposit(testuser.ID, currency2, decimal.NewFromFloat(50.0))

	t.Run("should return all balances for the user", func(t *testing.T) {
		balances, err := walletService.GetBalances(testuser.ID, []string{currency1, currency2})
		assert.NoError(t, err)
		assert.Len(t, balances, 2)

		for _, balance := range balances {
			switch balance.Currency {
			case currency1:
				assert.True(t, decimal.NewFromFloat(100.0).Equal(balance.Amount))
			case currency2:
				assert.True(t, decimal.NewFromFloat(50.0).Equal(balance.Amount))
			}
		}
	})
}

func TestGetTransactionHistory(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	testuser := userGenerator.Generate()
	tx.Create(testuser)

	walletService := services.NewWalletService(tx)

	currency := "USDT"

	walletService.Deposit(testuser.ID, currency, decimal.NewFromFloat(100.0))
	walletService.Withdraw(testuser.ID, currency, decimal.NewFromFloat(20.0))
	walletService.Deposit(testuser.ID, currency, decimal.NewFromFloat(50.0))

	t.Run("should return transaction history for the user", func(t *testing.T) {
		transactions, cursor, err := walletService.GetTransactionHistory(testuser.ID, "", "", services.SortOrderDesc, 10)
		assert.NoError(t, err)
		assert.Len(t, transactions, 3)
		assert.NotEmpty(t, cursor)

		assert.Equal(t, models.Deposit, transactions[0].Type)
		assert.True(t, decimal.NewFromFloat(50.0).Equal(transactions[0].Amount))

		assert.Equal(t, models.Withdrawal, transactions[1].Type)
		assert.True(t, decimal.NewFromFloat(20.0).Equal(transactions[1].Amount))

		assert.Equal(t, models.Deposit, transactions[2].Type)
		assert.True(t, decimal.NewFromFloat(100.0).Equal(transactions[2].Amount))
	})

	t.Run("should return paginated transaction history", func(t *testing.T) {
		transactions, cursor, err := walletService.GetTransactionHistory(testuser.ID, "", "", services.SortOrderDesc, 2)
		assert.NoError(t, err)
		assert.Len(t, transactions, 2)
		assert.NotEmpty(t, cursor)

		nextTransactions, _, err := walletService.GetTransactionHistory(testuser.ID, "", cursor, services.SortOrderDesc, 2)
		assert.NoError(t, err)
		assert.Len(t, nextTransactions, 1)
	})
}
