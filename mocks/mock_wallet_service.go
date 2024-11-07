package mocks

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/wanliqun/go-wallet-app/models"
	"github.com/wanliqun/go-wallet-app/services"
)

var (
	_ services.IWalletService = &MockWalletService{}
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) Deposit(userID uint, currency string, amount decimal.Decimal) error {
	args := m.Called(userID, currency, amount)
	return args.Error(0)
}

func (m *MockWalletService) Withdraw(userID uint, currency string, amount decimal.Decimal) error {
	args := m.Called(userID, currency, amount)
	return args.Error(0)
}

func (m *MockWalletService) Transfer(senderID, recipientID uint, currency string, amount decimal.Decimal, memo string) error {
	args := m.Called(senderID, recipientID, currency, amount, memo)
	return args.Error(0)
}

func (m *MockWalletService) GetBalances(userID uint, currencies []string) ([]models.Vault, error) {
	args := m.Called(userID, currencies)
	return args.Get(0).([]models.Vault), args.Error(1)
}

func (m *MockWalletService) GetTransactionHistory(userID uint, txnType models.TransactionType, cursor string, order services.SortOrder, limit int) ([]models.Transaction, string, error) {
	args := m.Called(userID, txnType, cursor, order, limit)
	return args.Get(0).([]models.Transaction), args.String(1), args.Error(2)
}
