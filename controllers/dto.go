package controllers

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/wanliqun/go-wallet-app/config"
	"github.com/wanliqun/go-wallet-app/models"
)

func init() {
	// set up custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register currency validation
		v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
			currency := fl.Field().String()
			_, ok := config.AppConfig.Concurrencies[currency]
			return ok
		})

		// Register amount validation
		v.RegisterValidation("positive_decimal", func(fl validator.FieldLevel) bool {
			amount, ok := fl.Field().Interface().(decimal.Decimal)
			return ok && amount.GreaterThan(decimal.Zero)
		})

		// Register currency limit validation
		v.RegisterValidation("currencyLimit", func(fl validator.FieldLevel) bool {
			currencies, ok := fl.Field().Interface().([]string)
			if !ok {
				return false
			}
			// Check if the length of currencies is between 1 and 30
			return len(currencies) >= 1 && len(currencies) <= 30
		})
	}
}

// DepositRequest represents the incoming request body for deposit operations
type DepositRequest struct {
	Currency string          `json:"currency" binding:"required,currency"`
	Amount   decimal.Decimal `json:"amount" binding:"required,positive_decimal"`
}

// WithdrawRequest represents the incoming request body for withdrawal operations
type WithdrawRequest DepositRequest

// TransferRequest represents the incoming request body for transfer operations
type TransferRequest struct {
	Recipient string          `json:"recipient" binding:"required"`
	Currency  string          `json:"currency" binding:"required,currency"`
	Amount    decimal.Decimal `json:"amount" binding:"required,positive_decimal"`
	Memo      string          `json:"memo,omitempty"`
}

// GetBalancesQuery represents the incoming request body for balance retrieval
type GetBalancesQuery struct {
	Currencies []string `form:"currency" binding:"required,currencyLimit"` // List of currencies to filter by
}

// GetTransactionHistoryRequest represents the request for retrieving paginated transaction history with filters
type GetTransactionHistoryRequest struct {
	Type   string `json:"type,omitempty" binding:"required,oneof=deposit withdrawal transfer_out transfer_in"` // Filter by transaction type (e.g., "deposit", "withdrawal")
	Cursor string `json:"cursor,omitempty"`                                                                    // Encoded cursor for keyset pagination
	Limit  int    `json:"limit,omitempty"`                                                                     // Number of records to fetch
	Order  string `json:"order,omitempty" binding:"omitempty,oneof=asc desc"`                                  // Sort order (e.g., "asc", "desc")
}

// GetTransactionHistoryResponse represents the response for paginated transaction history
type GetTransactionHistoryResponse struct {
	Data       []models.Transaction `json:"data"`        // Array of transaction objects
	NextCursor string               `json:"next_cursor"` // Encoded cursor for next page
}
