package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionType string

const (
	Deposit     TransactionType = "deposit"
	Withdrawal  TransactionType = "withdrawal"
	TransferOut TransactionType = "transfer_out"
	TransferIn  TransactionType = "transfer_in"
)

type Transaction struct {
	gorm.Model
	UserID         uint            `gorm:"not null;index:idx_user_type_timestamp_id,priority:1;index:idx_user_timestamp_id,priority:1" json:"user_id"`
	CounterpartyID *uint           `json:"counterparty_id"` // Pointer allows nulls
	Type           TransactionType `gorm:"size:16;index:idx_user_type_timestamp_id,priority:2" json:"type"`
	Amount         decimal.Decimal `gorm:"type:numeric(64,0);not null" json:"amount"`
	Currency       string          `gorm:"size:32;not null" json:"currency"`
	Memo           string          `gorm:"size:256" json:"memo,omitempty"`
	Timestamp      time.Time       `gorm:"autoCreateTime:milli;index:idx_user_type_timestamp_id,priority:3;index:idx_user_timestamp_id,priority:2" json:"timestamp"`
	ID             uint            `gorm:"primaryKey;index:idx_user_type_timestamp_id,priority:4;index:idx_user_timestamp_id,priority:3"`
}
