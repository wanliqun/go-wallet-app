package services

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/wanliqun/go-wallet-app/models"
	"github.com/wanliqun/go-wallet-app/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SortOrder int

const (
	SortOrderAsc SortOrder = iota
	SortOrderDesc
)

var (
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientBalance = errors.New("insufficient balance")

	_ IWalletService = &WalletService{}
)

type IWalletService interface {
	Deposit(userID uint, currency string, amount decimal.Decimal) error
	Withdraw(userID uint, currency string, amount decimal.Decimal) error
	Transfer(senderID, recipientID uint, currency string, amount decimal.Decimal, memo string) error
	GetBalances(userID uint, currencies []string) ([]models.Vault, error)
	GetTransactionHistory(userID uint, txnType models.TransactionType, cursor string, order SortOrder, limit int) ([]models.Transaction, string, error)
}

// WalletService represents the service for wallet-related operations
type WalletService struct {
	DB *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{DB: db}
}

func (s *WalletService) Deposit(userID uint, currency string, amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Upsert the Vault record using ON CONFLICT clause
		vault := models.Vault{
			UserID:   userID,
			Currency: currency,
			Amount:   amount,
		}
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "currency"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"amount": gorm.Expr("EXCLUDED.amount + ?", amount),
			}),
		}).Create(&vault).Error
		if err != nil {
			return err
		}

		// Record the deposit in transaction history
		transaction := models.Transaction{
			UserID:   userID,
			Type:     models.Deposit,
			Amount:   amount,
			Currency: currency,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *WalletService) Withdraw(userID uint, currency string, amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Attempt to decrement the amount atomically, ensuring the balance doesn't go negative
		result := tx.Model(&models.Vault{}).
			Where("user_id = ? AND currency = ? AND amount >= ?", userID, currency, amount).
			Update("amount", gorm.Expr("amount - ?", amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrInsufficientBalance
		}

		// Record the withdrawal in transaction history
		transaction := models.Transaction{
			UserID:   userID,
			Type:     models.Withdrawal,
			Amount:   amount,
			Currency: currency,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *WalletService) Transfer(senderID, recipientID uint, currency string, amount decimal.Decimal, memo string) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	// Validate recipient (cannot be the sender)
	if recipientID == senderID {
		return errors.New("cannot transfer to self")
	}

	// Start a database transaction
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Deduct from sender's vault atomically
		result := tx.Model(&models.Vault{}).
			Where("user_id = ? AND currency = ? AND amount >= ?", senderID, currency, amount).
			Updates(map[string]interface{}{
				"amount": gorm.Expr("amount - ?", amount),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrInsufficientBalance
		}

		// Find or create recipient's vault
		var recipientVault models.Vault
		if err := tx.FirstOrCreate(&recipientVault, models.Vault{
			UserID:   recipientID,
			Currency: currency,
		}).Error; err != nil {
			// If there's any error other than duplicate key, return it
			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				return err
			}
		}

		// Add to recipient's vault
		result = tx.Model(&models.Vault{}).
			Where("user_id = ? AND currency = ?", recipientID, currency).
			Update("amount", gorm.Expr("amount + ?", amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("failed to update recipient's vault")
		}

		// Create transaction records for sender and recipient as a batch
		batchTxns := []*models.Transaction{
			{ // transfer out
				UserID:         senderID,
				Type:           models.TransferOut,
				Amount:         amount,
				Currency:       currency,
				Memo:           memo,
				CounterpartyID: &recipientID,
			},
			{ // transfer in
				UserID:         recipientID,
				Type:           models.TransferIn,
				Amount:         amount,
				Currency:       currency,
				Memo:           memo,
				CounterpartyID: &senderID,
			},
		}
		// Batch insert the transactions
		return tx.Create(batchTxns).Error
	})
}

func (s *WalletService) GetBalances(userID uint, currencies []string) ([]models.Vault, error) {
	var vaults []models.Vault
	err := s.DB.Model(&models.Vault{}).
		Where("user_id = ? AND currency IN ?", userID, currencies).
		Find(&vaults).Error
	if err != nil {
		return nil, err
	}

	return vaults, nil
}

// GetTransactionHistory retrieves paginated transaction history using a unique cursor with filters
func (s *WalletService) GetTransactionHistory(
	userID uint, txnType models.TransactionType, cursor string, order SortOrder, limit int) ([]models.Transaction, string, error) {
	var transactions []models.Transaction

	query := s.DB.Where("user_id = ?", userID)

	// Apply transaction type filter if provided
	if txnType != "" {
		query = query.Where("type = ?", txnType)
	}

	// Decode the cursor if provided for pagination
	if cursor != "" {
		timestamp, txnID, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}

		// Apply keyset pagination using cursor values
		if order == SortOrderAsc {
			query = query.Where("(timestamp > ?) OR (timestamp = ? AND id > ?)", timestamp, timestamp, txnID)
		} else {
			query = query.Where("(timestamp < ?) OR (timestamp = ? AND id < ?)", timestamp, timestamp, txnID)
		}
	}

	// Apply sort order
	orderStr := "desc"
	if order == SortOrderAsc {
		orderStr = "asc"
	}
	query = query.Order("timestamp " + orderStr + ", id " + orderStr)

	// Limit the number of records retrieved
	if limit == 0 {
		limit = 10 // Default limit
	}
	query = query.Limit(limit)

	// Execute query
	if err := query.Find(&transactions).Error; err != nil {
		return nil, "", err
	}

	// Generate next cursor if there are more results
	var nextCursor string
	if len(transactions) > 0 {
		lastTransaction := transactions[len(transactions)-1]
		nextCursor = utils.EncodeCursor(lastTransaction.Timestamp, lastTransaction.ID)
	}

	return transactions, nextCursor, nil
}
