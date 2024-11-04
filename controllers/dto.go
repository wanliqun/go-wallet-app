package controllers

// DepositRequest represents the incoming request body for deposit operations
type DepositRequest struct {
	Currency string `json:"currency" binding:"required"`
	Amount   string `json:"amount" binding:"required"`
}

// WithdrawRequest represents the incoming request body for withdrawal operations
type WithdrawRequest DepositRequest

// TransferRequest represents the incoming request body for transfer operations
type TransferRequest struct {
	Recipient string `json:"recipient" binding:"required"`
	Currency  string `json:"currency" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	Memo      string `json:"memo,omitempty"`
}

// GetBalancesRequest represents the incoming request body for balance retrieval
type GetBalancesRequest struct {
	Offset uint `json:"offset,omitempty"`
	Limit  uint `json:"limit,omitempty"`
}

// GetTransactionHistoryRequest represents the incoming request body for transaction history retrieval
type GetTransactionHistoryRequest struct {
	Cursor uint   `json:"cursor,omitempty"`
	Limit  uint   `json:"limit,omitempty"`
	Type   string `json:"type,omitempty"`
	Order  string `json:"order,omitempty"`
}
