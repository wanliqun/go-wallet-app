package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Vault struct {
	gorm.Model
	UserID   uint            `gorm:"index;not null" json:"user_id"`
	Currency string          `gorm:"size:32;not null" json:"currency"`
	Amount   decimal.Decimal `gorm:"type:numeric(64,0);default:0" json:"amount"`
	User     User            `gorm:"foreignKey:UserID"`
}
