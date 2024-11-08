package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wanliqun/go-wallet-app/models"
	"github.com/wanliqun/go-wallet-app/services"
	"github.com/wanliqun/go-wallet-app/utils"
)

type WalletController struct {
	WalletService services.IWalletService
	UserService   services.IUserService
}

func NewWalletController(wallet services.IWalletService, user services.IUserService) *WalletController {
	return &WalletController{WalletService: wallet, UserService: user}
}

// POST /deposit
func (ctrl *WalletController) Deposit(c *gin.Context) {
	var cRequest DepositRequest
	if err := c.ShouldBindJSON(&cRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user := c.MustGet("user").(*models.User)
	err := ctrl.WalletService.Deposit(user.ID, cRequest.Currency, cRequest.Amount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// POST /withdraw
func (ctrl *WalletController) Withdraw(c *gin.Context) {
	var cRequest WithdrawRequest
	if err := c.ShouldBindJSON(&cRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user := c.MustGet("user").(*models.User)
	err := ctrl.WalletService.Withdraw(user.ID, cRequest.Currency, cRequest.Amount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// POST /transfer
func (ctrl *WalletController) Transfer(c *gin.Context) {
	var cRequest TransferRequest
	if err := c.ShouldBindJSON(&cRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user := c.MustGet("user").(*models.User)

	recipient, ok, err := ctrl.UserService.GetUserByName(cRequest.Recipient)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, services.ErrUserNotFound)
		return
	}

	if err := ctrl.WalletService.Transfer(user.ID, recipient.ID, cRequest.Currency, cRequest.Amount, cRequest.Memo); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// GET /balances
func (ctrl *WalletController) GetBalances(c *gin.Context) {
	var cRequest GetBalancesQuery
	if err := c.ShouldBindQuery(&cRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user := c.MustGet("user").(*models.User)

	vaults, err := ctrl.WalletService.GetBalances(user.ID, cRequest.Currencies)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, vaults)
}

// GET /transactions
func (ctrl *WalletController) GetTransactionHistory(c *gin.Context) {
	var cRequest GetTransactionHistoryQuery
	if err := c.ShouldBindQuery(&cRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user := c.MustGet("user").(*models.User)

	txnType := models.TransactionType(cRequest.Type)
	sortOrder := services.SortOrderDesc
	if cRequest.Order == "asc" {
		sortOrder = services.SortOrderAsc
	}

	transactions, nextCursor, err := ctrl.WalletService.GetTransactionHistory(user.ID, txnType, cRequest.Cursor, sortOrder, cRequest.Limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// Return transactions with the next cursor for pagination
	utils.SuccessResponse(c, GetTransactionHistoryResponse{
		Transactions: transactions,
		NextCursor:   nextCursor,
	})
}
