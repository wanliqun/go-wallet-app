package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionType string

const (
	Deposit      TransactionType = "deposit"
	Withdraw     TransactionType = "withdraw"
	TransferFrom TransactionType = "transfer_from"
	TransferTo   TransactionType = "transfer_to"
)

type Transaction struct {
	gorm.Model
	UserID         uint            `gorm:"index;not null" json:"user_id"`
	CounterpartyID *uint           `json:"counterparty_id"` // Pointer allows nulls
	Type           TransactionType `gorm:"type:enum('deposit', 'withdraw', 'transfer_from', 'transfer_to');not null" json:"type"`
	Amount         decimal.Decimal `gorm:"type:numeric(64,0);not null" json:"amount"`
	Currency       string          `gorm:"size:32;not null" json:"currency"`
	Memo           string          `gorm:"size:256" json:"memo,omitempty"`
	Timestamp      time.Time       `gorm:"default:current_timestamp" json:"timestamp"`
}
